package services_test

import (
	"errors"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"
	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"
)

type testEmailRepository struct {
	err    error
	emails []string
}

func newEmailRepository(err error, emails []string) *testEmailRepository {
	return &testEmailRepository{err: err, emails: emails}
}

func (e *testEmailRepository) AddEmail(_ string) error {
	return e.err
}

func (e *testEmailRepository) EmailList() []string {
	return e.emails
}

func TestSubscribeService_EmailList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		emails []string
	}{
		{
			name:   "check email list",
			emails: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newEmailRepository(nil, tt.emails)

			s := services.NewSubscribeService(r)
			require.Equal(t, tt.emails, s.EmailList())
		})
	}
}

func TestSubscribeService_Subscribe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		err       error
		expectErr require.ErrorAssertionFunc
	}{

		{
			name:      "without error",
			err:       nil,
			expectErr: require.NoError,
		},
		{
			name:      "with known error",
			err:       email.ErrAlreadyExists,
			expectErr: require.Error,
		},
		{
			name:      "with unknown error",
			err:       errors.New("unknown error"),
			expectErr: require.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newEmailRepository(tt.err, nil)

			s := services.NewSubscribeService(r)
			tt.expectErr(t, s.Subscribe("test"))
		})
	}
}
