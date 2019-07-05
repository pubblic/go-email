package email

import "net/smtp"

const (
	NAVER_SMTP     = "smtp.naver.com"
	NAVER_SMTP_TLS = "smtp.naver.com:587"
)

func Naver(a *Auth) (*smtp.Client, error) {
	return connect(a, NAVER_SMTP, NAVER_SMTP_TLS)
}
