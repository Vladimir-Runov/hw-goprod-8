package internal

import "fmt"

type Notifier interface {
	Send(customer string) error
}

type EmailSender struct{}

func (e *EmailSender) Send(customer string) error {
	fmt.Printf("E-mail: Уведомление отправлено клиенту %s\n", customer)
	return nil
}

type SMSSender struct{}

func (e *SMSSender) Send(customer string) error {
	fmt.Printf("SMS: Уведомление отправлено клиенту %s\n", customer)
	return nil
}
