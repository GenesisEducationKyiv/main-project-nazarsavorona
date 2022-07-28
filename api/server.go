package api

import (
	"BTCRateCheckService/internal"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/smtp"
	"strings"
	"sync"
)

const dataPath = "./resources/emails.dat"

type UnsentEmailsJSON struct {
	UnsentEmails []string
}

type Server struct {
	*mux.Router
	smtp.Auth
	email  string
	emails []string
}

func NewServer(email, password string) *Server {
	server := &Server{
		Router: mux.NewRouter(),
		Auth:   internal.NewLoginAuth(email, password),
		email:  email,
		emails: []string{},
	}

	server.emails, _ = internal.ReadLines(dataPath)

	server.routes()

	return server
}

func (server *Server) routes() {
	server.HandleFunc("/api/rate", server.rate()).Methods("GET")
	server.HandleFunc("/api/subscribe", server.subscribe()).Methods("POST")
	server.HandleFunc("/api/sendEmails", server.sendEmails()).Methods("POST")
}

func (server *Server) rate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		btcRate, err := server.getBTCRate("UAH")

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(writer).Encode(btcRate)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (server *Server) getBTCRate(currency string) (string, error) {
	response, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=BTC%s", currency))

	if err != nil {
		return "", err
	}

	var btcRate struct {
		Price string `json:"price"`
	}

	err = json.NewDecoder(response.Body).Decode(&btcRate)

	if err != nil {
		return "", err
	}

	return btcRate.Price, nil
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

		if internal.BinarySearch(server.emails, email) {
			http.Error(writer, email+" is already subscribed!", http.StatusConflict)
			return
		} else {
			server.emails = internal.InsertSorted(server.emails, email)

			err = internal.WriteLines(server.emails, dataPath)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			writer.WriteHeader(http.StatusOK)

			_, err = writer.Write([]byte(email + " has been added successfully!"))

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func (server *Server) sendEmails() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		subject := "BTC to UAH"
		rate, err := server.getBTCRate("UAH")

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		unsentEmails := []string{}
		var mutex sync.Mutex

		for _, email := range server.emails {
			go func(email string, mutex *sync.Mutex) {
				if err = internal.SendEmail(server.Auth, server.email, email, subject, rate); err != nil {
					mutex.Lock()
					unsentEmails = append(unsentEmails, email)
					mutex.Unlock()
				}
			}(email, &mutex)
		}

		writer.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(writer).Encode(&UnsentEmailsJSON{UnsentEmails: unsentEmails})

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
