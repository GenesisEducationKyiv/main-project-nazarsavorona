package integration_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/database"
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		emails       []string
		addEmail     string
		resultEmails []string
		expectErr    require.ErrorAssertionFunc
	}{
		{
			name:         "Add new email to empty repository",
			emails:       []string{},
			addEmail:     "1",
			resultEmails: []string{"1"},
			expectErr:    require.NoError,
		},
		{
			name:         "Add new email to non-empty repository",
			emails:       []string{"1", "2", "3"},
			addEmail:     "4",
			resultEmails: []string{"1", "2", "3", "4"},
			expectErr:    require.NoError,
		},
		{
			name:         "Add duplicate email to non-empty repository",
			emails:       []string{"1", "2", "3"},
			addEmail:     "2",
			resultEmails: []string{"1", "2", "3"},
			expectErr:    require.Error,
		},
	}

	for i, tt := range tests {
		tt := tt
		i := i

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.CreateTemp("", fmt.Sprintf("repository-test-%d", i))
			if err != nil {
				t.Fatal(err)
			}

			defer func(name string) {
				err = os.Remove(name)
				if err != nil {
					log.Println(err)
				}
			}(file.Name())

			db := database.NewFileDatabase(file)

			for _, e := range tt.emails {
				err = db.AddEmail(e)
				if err != nil {
					t.Fatal(err)
				}
			}

			repo := email.NewRepository(db)

			err = repo.AddEmail(tt.addEmail)
			tt.expectErr(t, err)

			require.ElementsMatch(t, tt.resultEmails, repo.EmailList())
		})
	}
}
