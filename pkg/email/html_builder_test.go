package email_test

import (
	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/email"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHTMLMessageBuilder_Construct(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		body    string
		want    []byte
	}{
		{
			name:    "without new lines",
			subject: "subject",
			body:    "body",
			want: []byte("Subject: subject\r\n" +
				"MIME-version: 1.0;\r\nContent-Type: text/html; " +
				"charset=\"UTF-8\";\r\n\r\n<html><body><h3>body</h3></body></html>"),
		},
		{
			name:    "with new lines",
			subject: "subject",
			body:    "body\nwith\nnew\nlines",
			want: []byte("Subject: subject\r\n" +
				"MIME-version: 1.0;\r\nContent-Type: text/html; " +
				"charset=\"UTF-8\";\r\n\r\n<html><body><h3>body<br>with<br>new<br>lines</h3></body></html>"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := email.HTMLMessageBuilder{}

			got := builder.Construct(tt.subject, tt.body)
			require.Equal(t, tt.want, got)
		})
	}
}
