package email_test

import (
	"testing"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"

	"github.com/stretchr/testify/require"
)

type testDB struct {
	emails []string
}

func (t *testDB) AddEmail(email string) error {
	t.emails = append(t.emails, email)
	return nil
}

func (t *testDB) Emails() ([]string, error) {
	return t.emails, nil
}

func TestRepository_AddEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		db    email.Database
		email string
		err   error
	}{
		{
			name:  "new email",
			db:    &testDB{},
			email: "test@ex.com",
			err:   nil,
		},
		{
			name: "existing email",
			db: &testDB{
				emails: []string{"test@ex.com"},
			},
			email: "test@ex.com",
			err:   email.ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := email.NewRepository(tt.db)
			err := r.AddEmail(tt.email)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestRepository_EmailList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		db    email.Database
		toAdd []string
		want  []string
	}{
		{
			name:  "empty list",
			db:    &testDB{},
			toAdd: nil,
			want:  []string{},
		},
		{
			name: "non-empty list",
			db: &testDB{
				emails: []string{"1", "2", "3"},
			},
			toAdd: nil,
			want:  []string{"1", "2", "3"},
		},
		{
			name: "non-empty list with new emails",
			db: &testDB{
				emails: []string{"1", "2", "3"},
			},
			toAdd: []string{"4", "5"},
			want:  []string{"1", "2", "3", "4", "5"},
		},
		{
			name: "non-empty list with existing emails",
			db: &testDB{
				emails: []string{"1", "2", "3"},
			},
			toAdd: []string{"2", "3"},
			want:  []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := email.NewRepository(tt.db)

			for _, e := range tt.toAdd {
				_ = r.AddEmail(e)
			}

			got := r.EmailList()
			require.ElementsMatch(t, tt.want, got)
		})
	}
}
