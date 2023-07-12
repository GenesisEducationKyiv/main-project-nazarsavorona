package email

import (
	"fmt"
	"log"
	"sync"
)

type Database interface {
	AddEmail(email string) error
	Emails() ([]string, error)
}

type emptyStruct struct{}

type Repository struct {
	emails map[string]emptyStruct
	db     Database
	mutex  sync.RWMutex
}

func NewRepository(db Database) *Repository {
	r := &Repository{
		emails: map[string]emptyStruct{},
		db:     db,
	}

	err := r.fillEmailsFromDatabase()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return r
}

var ErrAlreadyExists = fmt.Errorf("email already exists")

func (r *Repository) AddEmail(email string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.emails[email]; ok {
		return ErrAlreadyExists
	}

	err := r.db.AddEmail(email)
	if err != nil {
		return err
	}

	r.emails[email] = emptyStruct{}

	return nil
}

func (r *Repository) EmailList() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	emailList := make([]string, len(r.emails))

	i := 0
	for k := range r.emails {
		emailList[i] = k
		i++
	}

	return emailList
}

func (r *Repository) fillEmailsFromDatabase() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	emailList, err := r.db.Emails()
	if err != nil {
		return err
	}

	for _, currentEmail := range emailList {
		r.emails[currentEmail] = emptyStruct{}
	}

	return nil
}
