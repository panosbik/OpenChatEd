package helpers

import (
	"bytes"
	"html/template"

	"gopkg.in/gomail.v2"
)

func SendEmail(dialer *gomail.Dialer, subject, templatePath string, data any, to []string) (err error) {
	// Get html
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)

	if err != nil {
		return
	}

	// Apply a parsed template to the specified data object
	t.Execute(&body, data)

	m := gomail.NewMessage()
	m.SetHeader("From", dialer.Username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	// Send the email
	if err = dialer.DialAndSend(m); err != nil {
		return
	}
	return nil
}
