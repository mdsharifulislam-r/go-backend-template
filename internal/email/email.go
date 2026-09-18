package email

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/mdsharifulislam-r/go-backend-template/config"
)

type Message struct {
	To      string
	Subject string
	HTML    string
}

type Service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Send(msg Message) error {
	if s.cfg.EmailUser == "" || s.cfg.EmailPass == "" {
		log.Printf("[email:dev] to=%s subject=%s\n%s", msg.To, msg.Subject, msg.HTML)
		return nil
	}

	from := s.cfg.EmailFrom
	if from == "" {
		from = s.cfg.EmailUser
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.EmailHost, s.cfg.EmailPort)
	auth := smtp.PlainAuth("", s.cfg.EmailUser, s.cfg.EmailPass, s.cfg.EmailHost)

	headers := []string{
		fmt.Sprintf("From: Go Backend <%s>", from),
		fmt.Sprintf("To: %s", msg.To),
		fmt.Sprintf("Subject: %s", msg.Subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=\"UTF-8\"",
		"",
		msg.HTML,
	}

	if err := smtp.SendMail(addr, auth, from, []string{msg.To}, []byte(strings.Join(headers, "\r\n"))); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Mail sent successfully to %s", msg.To)
	return nil
}
