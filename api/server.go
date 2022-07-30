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
	"strconv"
	"strings"
	"sync"
)

const dataFolder = "./resources"
const dataPath = "./resources/emails.dat"

type UnsentEmailsJSON struct {
	UnsentEmails []string
}

type Server struct {
	*mux.Router
	smtp.Auth
	email    string
	emails   []string
	template *template.Template
}

func NewServer(email, password string) *Server {
	functionMap := template.FuncMap{"add": func(x, y int) int { return x + y }}

	server := &Server{
		Router:   mux.NewRouter(),
		Auth:     internal.NewLoginAuth(email, password),
		email:    email,
		emails:   []string{},
		template: template.Must(template.New("").Funcs(functionMap).ParseGlob("./templates/*.gohtml")),
	}

	if _, err := os.Stat(dataFolder); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dataFolder, os.ModePerm)
		if err != nil {
			log.Println(err.Error())
		}
	}

	server.emails, _ = internal.ReadLines(dataPath)

	server.routes()

	return server
}

func (server *Server) routes() {
	server.HandleFunc("/", server.index()).Methods("GET")
	server.HandleFunc("/api/rate", server.rate()).Methods("GET")
	server.HandleFunc("/api/subscribe", server.subscribe()).Methods("POST")
	server.HandleFunc("/api/sendEmails", server.sendEmails()).Methods("POST")
	server.HandleFunc("/subscribe", server.webSubscribe()).Methods("POST")
	server.HandleFunc("/sendEmails", server.webSendEmails()).Methods("POST")

	http.Handle("/", server)
}

func (server *Server) rate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		btcRate, err := server.getBTCRate("UAH")

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
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

		if internal.BinarySearch(server.emails, email) {
			http.Error(writer, email+" is already subscribed!", http.StatusConflict)
			return
		} else {
			if server.handleNewSubscriber(email, err) == nil {
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

func (server *Server) handleNewSubscriber(email string, err error) error {
	err = server.addNewEmail(email, err)

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

func (server *Server) addNewEmail(email string, err error) error {
	var mutex sync.Mutex

	mutex.Lock()

	server.emails = internal.InsertSorted(server.emails, email)
	err = internal.WriteLines(server.emails, dataPath)

	mutex.Unlock()

	return err
}

func (server *Server) sendEmails() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err, unsentEmails := server.startSendingEmails()

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

func (server *Server) startSendingEmails() (error, []string) {
	subject := "BTC to UAH"
	rate, err := server.getBTCRate("UAH")
	body := fmt.Sprintf("Current exchange rate:\n 1 BTC = %s UAH", getFormattedCurrency(rate))

	unsentEmails := []string{}

	if err != nil {
		return err, unsentEmails
	}

	var mutex sync.Mutex

	for _, email := range server.emails {
		go func(email, subject, body string, mutex *sync.Mutex) {
			if sendErr := server.sendEmail(email, subject, body); sendErr != nil {
				mutex.Lock()
				unsentEmails = append(unsentEmails, email)
				mutex.Unlock()
			}
		}(email, subject, body, &mutex)
	}

	return err, unsentEmails
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
	return internal.SendEmail(server.Auth, server.email, email, subject, body)
}

func (server *Server) index() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		rate, _ := server.getBTCRate("UAH")

		err := server.template.ExecuteTemplate(writer, "index.gohtml", struct {
			Rate   string
			Emails []string
		}{getFormattedCurrency(rate), server.emails})

		if err != nil {
			http.Redirect(writer, request, "/", http.StatusInternalServerError)
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
			http.Redirect(writer, request, "/", http.StatusBadRequest)
			return
		}

		if internal.BinarySearch(server.emails, email) {
			http.Redirect(writer, request, "/", http.StatusSeeOther)
			return
		} else {
			err = server.handleNewSubscriber(email, err)
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
		err, _ := server.startSendingEmails()

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
