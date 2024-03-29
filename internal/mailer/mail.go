package mailer

type MailMessage struct {
	To      []string
	Subject string
	Body    []byte
	IsHtml  bool
}

type MailService interface {
	Send(msg MailMessage) error
}
