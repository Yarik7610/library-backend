package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Sender interface {
	Addr() string
	From() string
	Subject() string
	WithSubject(subject string)
	Send(body string, to []string) error
}

type sender struct {
	host    string
	port    string
	from    string
	auth    smtp.Auth
	subject string
}

func NewSender(from, fromPassword string) Sender {
	s := sender{
		host: "smtp.gmail.com",
		port: "587",
		from: from,
	}
	s.auth = smtp.PlainAuth("", s.from, fromPassword, s.host)
	return &s
}

func (s *sender) Addr() string {
	return s.host + ":" + s.port
}

func (s *sender) From() string {
	return s.from
}

func (s *sender) Subject() string {
	return s.subject
}

func (s *sender) WithSubject(subject string) {
	s.subject = subject
}

func (s *sender) Send(body string, to []string) error {
	message := s.buildMessage(body, to)
	err := smtp.SendMail(s.Addr(), s.auth, s.from, to, message)
	return err
}

func (s *sender) buildMessage(body string, to []string) []byte {
	headers := map[string]string{
		"From":         s.from,
		"To":           strings.Join(to, ", "),
		"Subject":      s.subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return []byte(msg.String())
}
