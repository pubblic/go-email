package email

import (
	"crypto/tls"
	"net/smtp"
)

type Auth struct {
	Username string
	Password string
}

func connect(a *Auth, server, serverTls string) (*smtp.Client, error) {
	c, err := smtp.Dial(serverTls)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{ServerName: server}
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
		server,
	))
	if err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}
