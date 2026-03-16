package email

import (
	"crypto/tls"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Sender struct {
	SMTPHost string
	SMTPPort string
	Username string
	Password string
	From     string
}

// SendOTP sends a simple plaintext OTP email. It logs errors but does not return them
// so that signup flow can continue even if email fails in development.
func (s *Sender) SendOTP(toEmail, subject, body string) {
	if s == nil || s.SMTPHost == "" || s.From == "" {
		log.Printf("email: SMTP config missing, skipping send to %s", toEmail)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	port, err := strconv.Atoi(s.SMTPPort)
	if err != nil || port == 0 {
		port = 587
	}

	d := gomail.NewDialer(s.SMTPHost, port, s.Username, s.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("email: failed to send to %s: %v", toEmail, err)
		return
	}

	log.Printf("email: sent OTP email to %s", toEmail)
}
