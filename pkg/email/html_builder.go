package email

import (
	"strings"

	"github.com/GenesisEducationKyiv/main-project-nazarsavorona/pkg/models"
)

type HTMLMessageBuilder struct{}

func (*HTMLMessageBuilder) Construct(message *models.Message) []byte {
	body := strings.ReplaceAll(message.Body, "\n", "<br>")

	htmlBody := "<html><body><h3>" + body + "</h3></body></html>"

	return []byte("Subject: " + message.Subject + "\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		htmlBody)
}
