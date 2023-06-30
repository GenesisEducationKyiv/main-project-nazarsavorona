package database

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFileDatabase(t *testing.T) {
	t.Parallel()

	filename := "AddEmailTest.txt"

	file, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = file.Close()
		_ = os.Remove(filename)
	}()

	tests := []struct {
		name   string
		db     *FileDatabase
		emails []string
		err    require.ErrorAssertionFunc
	}{
		{
			name:   "AddEmailTest",
			db:     NewFileDatabase(file),
			emails: []string{"test@test.com", "test1@test.com", "test2@test.com"},
			err:    require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, email := range tt.emails {
				err := tt.db.AddEmail(email)
				tt.err(t, err)
			}

			emails, err := tt.db.Emails()
			tt.err(t, err)
			require.Equal(t, tt.emails, emails)
		})
	}
}
