package email

import (
	"net/smtp"
)

const (
	GOOGLE_SMTP     = "smtp.gmail.com"
	GOOGLE_SMTP_TLS = "smtp.gmail.com:587"
)

func Google(a *Auth) (*smtp.Client, error) {
	return connect(a, GOOGLE_SMTP, GOOGLE_SMTP_TLS)
}
