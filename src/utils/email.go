package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/VolodymyrShabat/Test_ATN/src/storage/models"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
	"strconv"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

func SendEmail(user *models.User, data *EmailData, templateName string, config models.Config) error {
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort, err := strconv.Atoi(config.SMTPPort)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	template, err := ParseTemplateDir("./src/templates")
	if err != nil {
		return err
	}

	template = template.Lookup(templateName)
	template.Execute(&body, &data)
	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	fmt.Println("success", err)
	return nil
}
