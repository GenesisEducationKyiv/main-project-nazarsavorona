package email

import (
	"log"
	"sync"
)

type Database interface {
	AddEmail(email string) error
	Emails() ([]string, error)
}

type Repository struct {
	emails map[string]bool
	db     Database
	mutex  sync.RWMutex
}

func NewRepository(db Database) *Repository {
	r := &Repository{
		emails: map[string]bool{},
		db:     db,
	}

	err := r.emailsFromDatabase()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return r
}

func (r *Repository) AddEmail(email string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	err := r.db.AddEmail(email)
	if err != nil {
		return err
	}

	r.emails[email] = true

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

func (r *Repository) emailsFromDatabase() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	emailList, err := r.db.Emails()
	if err != nil {
		return err
	}

	for _, currentEmail := range emailList {
		r.emails[currentEmail] = true
	}

	return nil
}
