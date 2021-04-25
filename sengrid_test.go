package sendgrid

import (
	"os"
	"testing"
)

func TestSendMail(t *testing.T) {

	var sg = New(&Config{
		APIKey:       os.Getenv("SG_KEY"),
		SenderMail:   os.Getenv("SG_SENDER_MAIL"),
		SenderName:   os.Getenv("SG_SENDER_NAME"),
		BccAddresses: os.Getenv("SG_BCC"),
	})

	sg.SendMail(SendMailParams{
		Email:       os.Getenv("SG_RECEIVER_MAIL"),
		Name:        os.Getenv("SG_RECEIVER_NAME"),
		Subject:     "Test eail",
		HTMLContent: "<div>Hello world</div>",
	})
}
