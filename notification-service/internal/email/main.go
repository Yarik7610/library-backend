package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/Yarik7610/library-backend/api-gateway/config"
	"go.uber.org/zap"
)

type Sender interface {
	Addr() string
	WithSubject(subject string)
	Send(body string, to []string)
}

type sender struct {
	host    string
	port    string
	from    string
	auth    smtp.Auth
	subject string
}

func NewSender() Sender {
	s := sender{}
	s.host = "smtp.gmail.com"
	s.port = "587"
	s.from = config.Data.Mail
	auth := smtp.PlainAuth("", s.from, config.Data.MailPassword, s.host)
	s.auth = auth
	return &s
}

func (s *sender) Addr() string {
	return s.host + ":" + s.port
}

func (s *sender) WithSubject(subject string) {
	s.subject = subject
}

func (s *sender) Send(body string, to []string) {
	message := s.buildMessage(body, to)
	err := smtp.SendMail(s.Addr(), s.auth, s.from, to, message)
	if err != nil {
		zap.S().Errorf("Mail send to %v error: %v", to, err)
	}
}

func (s *sender) buildMessage(body string, to []string) []byte {
	headers := map[string]string{
		"From":    s.from,
		"To":      strings.Join(to, ", "),
		"Subject": s.subject,
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return []byte(msg.String())
}
