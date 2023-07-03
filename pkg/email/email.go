package email

import (
	"strconv"

	//"gopkg.in/gomail.v2"
	"net/smtp"
	"strings"
)

type SMTPServer struct {
	Host     string
	Port     int
	User     string
	Password string
}

func NewSMTPServer(Host string, Port int, User string, Password string) *SMTPServer {
	return &SMTPServer{
		Host:     Host,
		Port:     Port,
		User:     User,
		Password: Password,
	}
}

type EmailManager struct {
	SMTP        *SMTPServer
	SourceEmail string
}

func NewEmailManager(SMTP *SMTPServer, source string) *EmailManager {

	return &EmailManager{
		SMTP:        SMTP,
		SourceEmail: source,
	}
}

type SMTPManager interface {
	SendMail(to string, subject string, content string) error
}

func (m *EmailManager) SendMail(to []string, subject string, content string) error {

	auth := smtp.PlainAuth("", m.SMTP.User, m.SMTP.Password, m.SMTP.Host)

	msg := []byte(
		"From: " + m.SMTP.User + "\r\n" +
			"To: " + strings.Join(to, ", ") + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			content + "\r\n")

	err := smtp.SendMail(m.SMTP.Host+":"+strconv.Itoa(m.SMTP.Port), auth, m.SourceEmail, to, msg)

	return err
}
