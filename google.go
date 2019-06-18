package email

import (
	"crypto/tls"
	"net/smtp"
)

const (
	GOOGLE_SMTP     = "smtp.gmail.com"
	GOOGLE_SMTP_TLS = "smtp.gmail.com:587"
)

func Google(a *Auth) (*smtp.Client, error) {
	c, err := smtp.Dial(GOOGLE_SMTP_TLS)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{ServerName: GOOGLE_SMTP}
	err = c.StartTLS(config)
	if err != nil {
		c.Close()
		return nil, err
	}
	const identity = ""
	err = c.Auth(smtp.PlainAuth(
		identity,
		a.Username,
		a.Password,
		GOOGLE_SMTP,
	))
	if err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}
