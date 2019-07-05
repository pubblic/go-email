package email

import (
	"bytes"
	"io"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

func SyntaxMsgId(idLeft, idRight string) string {
	return "<" + idLeft + "@" + idRight + ">"
}

type HeaderBuilder struct {
	From        *mail.Address   // required
	To          []*mail.Address // required
	MessageId   string          // required
	Subject     string
	Date        time.Time
	ContentType string
	Header      mail.Header
}

func (b HeaderBuilder) Build() mail.Header {
	const e = 5
	h := make(mail.Header, len(b.Header)+e)
	for key, val := range b.Header {
		h[key] = val
	}
	set := func(key, val string) {
		if val == "" {
			return
		}
		h[key] = []string{val}
	}
	set("From", b.From.String())
	set("To", stringMailboxList(b.To))
	if !b.Date.IsZero() {
		set("Date", b.Date.Format(time.RFC1123Z))
	}
	set("Subject", b.Subject)
	set("Content-Type", b.ContentType)
	set("Message-ID", b.MessageId)
	return h
}

func stringMailboxList(list []*mail.Address) string {
	var buf strings.Builder
	for i, addr := range list {
		if i+1 < len(list) {
			buf.WriteString(", ")
		}
		buf.WriteString(addr.String())
	}
	return buf.String()
}

// Mail starts mail transaction and sends mail. Mail transaction is cancelled
// if error occurs.
func Mail(c *smtp.Client, from string, to []string, msg *mail.Message) error {
	// TODO: validate the arguments.
	return doMail(c, from, to, msg)
}

func bytesHeader(h mail.Header) *bytes.Buffer {
	const CRLF = "\r\n"
	buf := new(bytes.Buffer)
	for key, vals := range h {
		if len(vals) == 0 {
			continue
		}
		buf.WriteString(key)
		buf.WriteString(": ")
		buf.WriteString(vals[0])
		buf.WriteString(CRLF)
	}
	buf.WriteString(CRLF)
	return buf
}

func xCopy(dst io.Writer, src io.Reader) (int64, error) {
	if src == nil {
		return 0, nil
	}
	return io.Copy(dst, src)
}

func doMail(c *smtp.Client, from string, to []string, msg *mail.Message) (err error) {
	var startMail bool
	err = c.Mail(from)
	if err != nil {
		return err
	}
	startMail = true
	defer func() {
		if err != nil && startMail {
			c.Reset()
		}
	}()
	for _, address := range to {
		err := c.Rcpt(address)
		if err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = bytesHeader(msg.Header).WriteTo(w)
	if err != nil {
		return err
	}
	_, err = xCopy(w, msg.Body)
	if err != nil {
		return err
	}
	return w.Close()
}
