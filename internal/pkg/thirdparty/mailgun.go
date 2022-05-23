package thirdparty

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	thirdparty "github.com/vbetsun/surgeon-intern-app/internal/pkg/thirdparty/assets"
)

type MailgunCredentials struct {
	ApiKey            string
	Domain            string
	InviteEmailSender string
}

type Mailgun struct {
	client      *mailgun.MailgunImpl
	senderEmail string
}

func NewMailgun(mailJetCredentials *MailgunCredentials) *Mailgun {
	mailgun.Debug = true
	client := mailgun.NewMailgun(mailJetCredentials.Domain, mailJetCredentials.ApiKey)
	client.SetAPIBase(mailgun.APIBaseEU)
	return &Mailgun{client: client, senderEmail: mailJetCredentials.InviteEmailSender}
}

func (m *Mailgun) SendInviteMail(email string, activationLink string) error {
	subject := "Inbjudan till Ogbook"
	body := strings.Replace(thirdparty.UserInviteHtml, "{{USER_EMAIL}}", email, 1)
	body = strings.Replace(body, "{{PASSWORD_RESET_LINK}}", activationLink, 1)
	// The message object allows you to add attachments and Bcc recipients
	message := m.client.NewMessage(m.senderEmail, subject, "", email)
	message.SetHtml(body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, _, err := m.client.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
