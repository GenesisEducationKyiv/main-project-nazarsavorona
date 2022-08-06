package api

import (
	"BTCRateCheckService/internal"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const dataFolder = "./resources"
const dataPath = "./resources/emails.dat"

type UnsentEmailsJSON struct {
	UnsentEmails []string
}

type Server struct {
	Router          *mux.Router
	auth            smtp.Auth
	email           string
	emails          map[string]bool
	btcRate         float64
	lastTimeRequest time.Time
	template        *template.Template
}

func NewServer(email, password string) *Server {
	functionMap := template.FuncMap{"add": func(x, y int) int { return x + y }}

	server := &Server{
		Router:   mux.NewRouter(),
		auth:     internal.NewLoginAuth(email, password),
		email:    email,
		emails:   map[string]bool{},
		template: template.Must(template.New("").Funcs(functionMap).ParseGlob("./templates/*.gohtml")),
	}

	if _, err := os.Stat(dataFolder); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dataFolder, os.ModePerm)
		if err != nil {
			log.Println(err.Error())
		}
	}

	server.readEmailsFromFile()

	server.routes()

	var err error = nil

	server.btcRate, err = server.getBTCRate("UAH")

	if err != nil {
		panic(err.Error())
	}

	server.lastTimeRequest = time.Now()

	return server
}

func (server *Server) readEmailsFromFile() {
	emailList, _ := internal.ReadLines(dataPath)

	for _, currentEmail := range emailList {
		server.emails[currentEmail] = true
	}
}

func (server *Server) routes() {
	server.Router.HandleFunc("/", server.index()).Methods("GET")
	server.Router.HandleFunc("/conflict", server.conflict()).Methods("GET")
	server.Router.HandleFunc("/api/rate", server.rate()).Methods("GET")
	server.Router.HandleFunc("/api/subscribe", server.subscribe()).Methods("POST")
	server.Router.HandleFunc("/api/sendEmails", server.sendEmails()).Methods("POST")
	server.Router.HandleFunc("/subscribe", server.webSubscribe()).Methods("POST")
	server.Router.HandleFunc("/sendEmails", server.webSendEmails()).Methods("POST")

	http.Handle("/", server.Router)
}

func (server *Server) rate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		btcRate := 0.0
		var err error = nil

		if time.Now().Sub(server.lastTimeRequest) <= time.Second*15 {
			btcRate = server.btcRate
		} else {
			btcRate, err = server.getBTCRate("UAH")

			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}

			server.lastTimeRequest = time.Now()
		}

		writer.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(writer).Encode(json.Number(getFormattedCurrency(btcRate)))

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (server *Server) subscribe() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		email := request.Form.Get("email")
		email = strings.TrimSpace(email)

		if !internal.ValidateEmail(email) {
			http.Error(writer, "Invalid email!", http.StatusBadRequest)
			return
		}

		if _, isPresent := server.emails[email]; isPresent {
			http.Error(writer, email+" is already subscribed!", http.StatusConflict)
			return
		} else {
			if server.handleNewSubscriber(email) == nil {
				_, err = writer.Write([]byte(email + " has been added successfully!"))

				if err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
				}
			} else {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func (server *Server) handleNewSubscriber(email string) error {
	err := server.addNewEmail(email)

	if err != nil {
		return err
	}

	go func() {
		err = server.sendEmail(email, "Thank You for subscription!",
			"You will be receiving information about BTC to UAH exchange rates from now on.\n\nStay tuned!")

		if err != nil {
			log.Printf(err.Error())
		}
	}()

	return nil
}

func (server *Server) addNewEmail(email string) error {
	var mutex sync.Mutex

	mutex.Lock()

	server.emails[email] = true

	emailList := server.getEmailList()

	err := internal.WriteLines(dataPath, emailList)

	mutex.Unlock()

	return err
}

func (server *Server) getEmailList() []string {
	emailList := make([]string, len(server.emails))

	i := 0
	for k := range server.emails {
		emailList[i] = k
		i++
	}
	return emailList
}

func (server *Server) sendEmails() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		unsentEmails, err := server.startSendingEmails()

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(writer).Encode(&UnsentEmailsJSON{UnsentEmails: unsentEmails})

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (server *Server) startSendingEmails() ([]string, error) {
	subject := "BTC to UAH"
	rate, err := server.getBTCRate("UAH")
	body := fmt.Sprintf("Current exchange rate:\n 1 BTC = %s UAH", getFormattedCurrency(rate))

	unsentEmails := []string{}

	if err != nil {
		return unsentEmails, err
	}

	var mutex sync.Mutex
	var waitGroup sync.WaitGroup

	for email := range server.emails {
		waitGroup.Add(1)
		go func(email, subject, body string, mutex *sync.Mutex) {
			defer waitGroup.Done()

			if sendErr := server.sendEmail(email, subject, body); sendErr != nil {
				mutex.Lock()
				defer mutex.Unlock()

				unsentEmails = append(unsentEmails, email)
			}
		}(email, subject, body, &mutex)
	}

	waitGroup.Wait()

	return unsentEmails, err
}

func (server *Server) getBTCRate(currency string) (float64, error) {
	response, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=BTC%s", currency))

	if err != nil {
		return 0, err
	}

	var btcRate struct {
		Price string `json:"price"`
	}

	err = json.NewDecoder(response.Body).Decode(&btcRate)

	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(btcRate.Price, 64)
}

func (server *Server) sendEmail(email, subject, body string) error {
	return internal.SendEmail(server.auth, server.email, email, subject, body)
}

func (server *Server) index() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		emails := server.getEmailList()
		sort.Strings(emails)

		indexData := struct {
			Rate   string
			Emails []string
		}{getFormattedCurrency(server.btcRate), emails}

		err := server.template.ExecuteTemplate(writer, "index.gohtml", indexData)

		if err != nil {
			http.Redirect(writer, request, "/", http.StatusInternalServerError)
		}
	}
}

func (server *Server) conflict() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := server.template.ExecuteTemplate(writer, "conflict.gohtml", nil)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (server *Server) webSubscribe() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()

		if err != nil {
			return
		}

		email := request.Form.Get("email")
		email = strings.TrimSpace(email)

		if !internal.ValidateEmail(email) {
			http.Redirect(writer, request, "/", http.StatusSeeOther)
			return
		}

		if _, isPresent := server.emails[email]; isPresent {
			http.Redirect(writer, request, "/conflict", http.StatusSeeOther)
			return
		} else {
			err = server.handleNewSubscriber(email)
			if err != nil {
				http.Redirect(writer, request, "/", http.StatusBadRequest)
				return
			}
		}

		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}
}

func (server *Server) webSendEmails() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := server.startSendingEmails()

		if err != nil {
			http.Redirect(writer, request, "/", http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}
}

func getFormattedCurrency(btcRate float64) string {
	return fmt.Sprintf("%.2f", btcRate)
}
