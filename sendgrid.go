package sendgrid

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type Config struct {
	APIKey       string
	SenderMail   string
	SenderName   string
	BccAddresses string
	CcAddresses  string
	Logger       Logger
}

type Email struct {
	Email string
	Name  string
}

type Attachment struct {
	Content     string `json:"content,omitempty"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Disposition string `json:"disposition,omitempty"`
	ContentID   string `json:"content_id,omitempty"`
}

type Mailer struct {
	config *Config
	logger Logger
}

func New(config *Config) *Mailer {
	if config.APIKey == "" {
		panic("Sendgrid API Key is required")
	}
	var mailer = &Mailer{
		config: config,
		logger: log.New(os.Stdout, "\r\n", 0),
	}

	if config.Logger != nil {
		mailer.logger = config.Logger
	}
	return mailer
}

type SendMailParams struct {
	Subject     string
	Name        string
	Email       string
	Attachments []*Attachment

	// Custom data
	HTMLContent string

	// Dynamic template
	TemplateID string
	Data       map[string]interface{}

	AsmGroupID         int
	AsmGroupsToDisplay []int
}

func (mailer *Mailer) SendMail(params SendMailParams) error {

	var apiKey = mailer.config.APIKey
	var addressAlias = mailer.config.SenderMail
	var nameAlias = mailer.config.SenderName
	var bccMails = mailer.config.BccAddresses
	var ccMails = mailer.config.CcAddresses
	var userName = params.Name
	var userEmail = params.Email

	var m = mail.NewV3Mail()
	var e = mail.NewEmail(nameAlias, addressAlias)
	m.SetFrom(e)

	if params.AsmGroupID > 0 && len(params.AsmGroupsToDisplay) > 0 {
		m.SetASM(&mail.Asm{
			GroupID:         params.AsmGroupID,
			GroupsToDisplay: params.AsmGroupsToDisplay,
		})
	}

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(userName, userEmail),
	}
	p.AddTos(tos...)

	var bcc = mailer.splitEmails(bccMails, userEmail)
	if len(bcc) > 0 {
		p.AddBCCs(bcc...)
	}

	var cc = mailer.splitEmails(ccMails, userEmail)
	if len(cc) > 0 {
		p.AddCCs(cc...)
	}

	m.AddPersonalizations(p)

	m.Subject = params.Subject

	if params.TemplateID != "" {
		m.SetTemplateID(params.TemplateID)
		for key, value := range params.Data {
			p.SetDynamicTemplateData(key, value)
		}
	}

	if params.HTMLContent != "" {
		m.Content = []*mail.Content{
			mail.NewContent("text/html", params.HTMLContent),
		}
	}

	// Add attachments
	for _, at := range params.Attachments {
		m.AddAttachment(&mail.Attachment{
			Content:     at.Content,
			ContentID:   at.ContentID,
			Disposition: at.Disposition,
			Filename:    at.Filename,
			Name:        at.Name,
			Type:        at.Type,
		})
	}

	var request = sendgrid.GetRequest(apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)

	if err != nil {
		mailer.logger.Printf("Send mail error: %v\n", err)
		return err
	}

	if response.StatusCode >= 400 {
		mailer.logger.Printf("Send mail status_code = %d response = %v\n", response.StatusCode, response.Body)
		var errs ErrorResponse
		var err = json.Unmarshal([]byte(response.Body), &errs)
		if err == nil && len(errs.Errors) > 0 {
			return errs.Errors
		}

		return nil
	}

	return nil
}

func (mailer *Mailer) splitEmails(mailList, mailTo string) []*mail.Email {
	var cc = []*mail.Email{}
	if mailList != "" {
		var s1 = strings.Split(mailList, "|")
		if len(s1) > 0 {
			for _, nameMail := range s1 {
				var s2 = strings.Split(nameMail, ",")
				if len(s2) == 2 {
					if !strings.EqualFold(s2[1], mailTo) {
						cc = append(cc, mail.NewEmail(s2[0], s2[1]))

					}
				}
			}
		} else {
			var s2 = strings.Split(mailList, ",")
			if len(s2) == 2 {
				if s2[1] != mailTo {
					cc = append(cc, mail.NewEmail(s2[0], s2[1]))
				}
			}
		}
	}
	return cc
}
