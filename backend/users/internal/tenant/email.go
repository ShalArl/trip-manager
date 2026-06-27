package tenant

import (
	"fmt"
	"github.com/resend/resend-go/v3"
)

type EmailService struct {
	client   *resend.Client
	fromAddr string
}

func NewEmailService(apiKey string) *EmailService {
	if apiKey == "" {
		return nil
	}
	return &EmailService{
		client:   resend.NewClient(apiKey),
		fromAddr: "noreply@neatnode.xyz",
	}
}

func (e *EmailService) SendInvitation(toEmail, tenantName, inviteLink, role string) error {
	if e == nil {
		return nil
	}

	roleLabel := "Mitarbeiter"
	if role == "tenant_admin" {
		roleLabel = "Admin"
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h2 style="color: #0284c7;">Du wurdest eingeladen! 🌍</h2>
    <p>Du wurdest als <strong>%s</strong> zu <strong>%s</strong> auf Trip Manager eingeladen.</p>
    <a href="%s" style="display: inline-block; padding: 12px 24px; background: #0284c7; color: white; text-decoration: none; border-radius: 8px; margin: 16px 0;">
        Einladung annehmen
    </a>
    <p style="color: #666; font-size: 12px;">Dieser Link ist 7 Tage gültig.</p>
</body>
</html>`, roleLabel, tenantName, inviteLink)

	params := &resend.SendEmailRequest{
		From:    e.fromAddr,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("Einladung zu %s auf Trip Manager", tenantName),
		Html:    html,
	}

	_, err := e.client.Emails.Send(params)
	return err
}
