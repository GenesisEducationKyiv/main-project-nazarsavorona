package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/services"
)

type testEmailSender struct {
	count                int
	failedRequestAttempt int
}

func newTestEmailSender(failedRequestAttempt int) *testEmailSender {
	return &testEmailSender{failedRequestAttempt: failedRequestAttempt}
}

func (t *testEmailSender) SendEmail(_ string, _ []byte) error {
	t.count++
	if t.count == t.failedRequestAttempt {
		return fmt.Errorf("test error")
	}

	return nil
}

func TestEmailService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		tries      int
		mailSender services.EmailSender
		expectErr  require.ErrorAssertionFunc
	}{
		{
			name:       "all sent successfully",
			mailSender: newTestEmailSender(0),
			tries:      5,
			expectErr:  require.NoError,
		},
		{
			name:       "second fails",
			mailSender: newTestEmailSender(2),
			tries:      5,
			expectErr:  require.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			emails := make([]string, tt.tries)
			for i := range emails {
				emails[i] = "test"
			}

			s := services.NewEmailService(tt.mailSender, &email.HTMLMessageBuilder{})
			err := s.SendEmails(context.Background(), emails, &models.Message{})
			tt.expectErr(t, err)
		})
	}
}
