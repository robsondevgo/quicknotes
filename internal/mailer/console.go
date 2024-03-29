package mailer

import (
	"fmt"
	"strings"
)

type consoleMailService struct {
	from string
}

func NewConsoleMailService(from string) MailService {
	return consoleMailService{from: from}
}

func (cs consoleMailService) Send(msg MailMessage) error {
	fmt.Printf("From: %s \nTo: %s\nSubject: %s\nBody: %s\n",
		cs.from, strings.Join(msg.To, ","), msg.Subject, msg.Body)
	return nil
}
