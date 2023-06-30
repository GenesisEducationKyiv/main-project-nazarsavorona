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
	mutex  sync.Mutex
}

func NewRepository(db Database) *Repository {
	r := &Repository{
		emails: map[string]bool{},
		db:     db,
	}

	err := r.getEmailsFromDatabase()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return r
}

func (r *Repository) AddNewEmail(email string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	err := r.db.AddEmail(email)
	if err != nil {
		return err
	}

	r.emails[email] = true

	return nil
}

func (r *Repository) GetEmailList() []string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	emailList := make([]string, len(r.emails))

	i := 0
	for k := range r.emails {
		emailList[i] = k
		i++
	}

	return emailList
}

func (r *Repository) getEmailsFromDatabase() error {
	emailList, err := r.db.Emails()
	if err != nil {
		return err
	}

	for _, currentEmail := range emailList {
		r.emails[currentEmail] = true
	}

	return nil
}
