package email_test

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"
)

func TestSender_SendEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		returnErr error
		expectErr require.ErrorAssertionFunc
	}{
		{
			name:      "sent successfully",
			returnErr: nil,
			expectErr: require.NoError,
		},
		{
			name:      "sent with error",
			returnErr: fmt.Errorf("something wrong happened"),
			expectErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := email.NewSender("", "", nil,
				func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
					return tt.returnErr
				})

			err := s.SendEmail("", "", "")
			tt.expectErr(t, err)
		})
	}
}
