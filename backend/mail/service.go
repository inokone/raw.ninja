package mail

import (
	"bytes"
	_ "embed"
	"html/template"

	"github.com/inokone/photostorage/common"
	"github.com/rs/zerolog/log"
	"gopkg.in/mail.v2"
)

//go:embed "confirmation.html"
var ct string

//go:embed "passwordreset.html"
var pt string

// Service is a struct for a service sending mails for our users.
type Service struct {
	config       common.MailConfig
	dialer       *mail.Dialer
	confirmation *template.Template
	pwdReset     *template.Template
}

// NewService create a new `Service` entity based on the configuration.
// If SMTP server is not configured the service will not return error, just logs it as a warning.
func NewService(config common.MailConfig) Service {
	if len(config.SMTPAddress) == 0 {
		log.Warn().Msg("SMTP is not set up, e-mail sending functionality will not work correctly!")
	}
	return Service{
		config:       config,
		dialer:       mail.NewDialer(config.SMTPAddress, config.SMTPPort, config.SMTPUser, config.SMTPPassword),
		confirmation: mustLoadTemplate(ct),
		pwdReset:     mustLoadTemplate(pt),
	}
}

func mustLoadTemplate(tpl string) *template.Template {
	tmpl, err := template.New("email").Parse(tpl)
	if err != nil {
		panic("email template can not be parsed")
	}
	return tmpl
}

type templateData struct {
	Link string
}

// Send is a method of `Service` sends an e-mail to the recipient email address with the subject and body provided as parameters
// If SMTP server is not configured the service will not return error, just logs it as a warning.
func (s *Service) send(recipient string, subject string, body string) error {
	if len(s.config.SMTPAddress) == 0 {
		log.Warn().Msg("SMTP is not set up, failed to send the e-mail!")
		return nil
	}
	// Set up email message
	m := mail.NewMessage()
	m.SetHeader("From", s.config.NoReplyAddress)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send the email
	return s.dialer.DialAndSend(m)
}

// EmailConfirmation is a method of `Service` sends an e-mail confirmation message to the recipient email address
func (s *Service) EmailConfirmation(recipient string, confirmationURL string) error {
	var c bytes.Buffer
	if err := s.confirmation.Execute(&c, templateData{Link: confirmationURL}); err != nil {
		return err
	}
	return s.send(recipient, "E-mail Confirmation", c.String())
}

// PasswordReset is a method of `Service` sends a password reset message to the recipient email address
func (s *Service) PasswordReset(recipient string, resetURL string) error {
	var c bytes.Buffer
	if err := s.pwdReset.Execute(&c, templateData{Link: resetURL}); err != nil {
		return err
	}
	return s.send(recipient, "Password Reset", c.String())
}
