package database_test

import (
	"os"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/database"

	"github.com/stretchr/testify/require"
)

func TestFileDatabase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		emails []string
		err    require.ErrorAssertionFunc
	}{
		{
			name:   "AddEmailTest",
			emails: []string{"test@test.com", "test1@test.com", "test2@test.com"},
			err:    require.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			filename := tt.name + ".txt"

			file, err := os.Create(filename)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = file.Close()
				_ = os.Remove(filename)
			}()

			db := database.NewFileDatabase(file)

			for _, email := range tt.emails {
				err = db.AddEmail(email)
				tt.err(t, err)
			}

			emails, err := db.Emails()
			tt.err(t, err)
			require.Equal(t, tt.emails, emails)
		})
	}
}
