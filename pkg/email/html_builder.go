package email

import "strings"

type HTMLMessageBuilder struct{}

func (*HTMLMessageBuilder) Construct(subject, body string) []byte {
	body = strings.ReplaceAll(body, "\n", "<br>")

	htmlBody := "<html><body><h3>" + body + "</h3></body></html>"
	message := []byte("Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		htmlBody)

	return message
}
