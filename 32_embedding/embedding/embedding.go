package main

import (
	"errors"
	"fmt"
)

type Email struct {
	DisplayName string
	UserName    string
	Domain      string
}

func (e Email) Address() string {
	return fmt.Sprintf("%s %s@%s", e.DisplayName, e.UserName, e.Domain)
}

type Invite struct {
	ID int
	Email
}

func main() {
	invite := Invite{
		ID: 1,
		Email: Email{
			DisplayName: "Bor Kipe",
			UserName:    "bkipe",
			Domain:      "go.dev.com",
		},
	}

	if err := SendEmail(invite); err != nil {
		fmt.Printf("error sending email: %s\n", err)
	}
}

func SendEmail(i Invite) error {
	address := i.Address()
	if address == "" {
		return errors.New("no email address provided")
	}
	fmt.Printf("sent email to: %s\n", address)
	return nil
}
