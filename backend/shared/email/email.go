package email

import (
	"fmt"

	"github.com/resend/resend-go/v3"
)

type Service struct {
	client   *resend.Client
	fromAddr string
}

func NewService(apiKey string) *Service {
	if apiKey == "" {
		return nil
	}
	return &Service{
		client:   resend.NewClient(apiKey),
		fromAddr: "noreply@neatnode.xyz",
	}
}

func (s *Service) SendInvitation(toEmail, tenantName, inviteLink, role string) error {
	if s == nil {
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

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.fromAddr,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("Einladung zu %s auf Trip Manager", tenantName),
		Html:    html,
	})
	return err
}

func (s *Service) SendContactRequest(toEmail, advertiserName, advertiserEmail, tenantName, message string) error {
	if s == nil {
		return nil
	}
	msgHTML := ""
	if message != "" {
		msgHTML = fmt.Sprintf(`<p style="margin: 8px 0 0;"><strong>Nachricht:</strong> %s</p>`, message)
	}
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h2 style="color: #0284c7;">Neue Kontaktanfrage von einem Werbepartner 📢</h2>
    <p><strong>%s</strong> möchte dein Reisebüro <strong>%s</strong> bewerben.</p>
    <div style="background: #f4f4f5; border-radius: 8px; padding: 16px; margin: 16px 0;">
        <p style="margin: 0;"><strong>Name:</strong> %s</p>
        <p style="margin: 8px 0 0;"><strong>Email:</strong> <a href="mailto:%s">%s</a></p>
        %s
    </div>
</body>
</html>`, advertiserName, tenantName, advertiserName, advertiserEmail, advertiserEmail, msgHTML)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.fromAddr,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("Kontaktanfrage von Werbepartner %s", advertiserName),
		Html:    html,
	})
	return err
}

type InsightSummary struct {
	TenantID       string
	TenantName     string
	TopDestination string
	TotalLikes     int64
	PeakMonth      string
}

func (s *Service) SendInsightsReport(toEmail, advertiserName string, insights []InsightSummary) error {
	if s == nil {
		return nil
	}

	tenantsHTML := ""
	for _, ins := range insights {
		tenantsHTML += fmt.Sprintf(`
        <div style="border: 1px solid #e4e4e7; border-radius: 8px; padding: 16px; margin: 12px 0;">
            <h3 style="margin: 0 0 12px; color: #18181b;">%s</h3>
            <table style="width: 100%%; border-collapse: collapse;">
                <tr>
                    <td style="padding: 4px 0; color: #71717a; font-size: 14px;">Top Destination</td>
                    <td style="padding: 4px 0; text-align: right; font-size: 14px;">%s</td>
                </tr>
                <tr>
                    <td style="padding: 4px 0; color: #71717a; font-size: 14px;">Likes gesamt</td>
                    <td style="padding: 4px 0; text-align: right; font-size: 14px;">%d</td>
                </tr>
                <tr>
                    <td style="padding: 4px 0; color: #71717a; font-size: 14px;">Peak-Monat</td>
                    <td style="padding: 4px 0; text-align: right; font-size: 14px;">%s</td>
                </tr>
            </table>
        </div>`, ins.TenantName, ins.TopDestination, ins.TotalLikes, ins.PeakMonth)
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h2 style="color: #0284c7;">Deine wöchentlichen Travel Insights 🌍</h2>
    <p>Hallo <strong>%s</strong>, hier sind deine Travel Insights für diese Woche:</p>
    %s
    <p style="margin-top: 24px;">
        <a href="https://www.neatnode.xyz/advertiser" style="display: inline-block; padding: 12px 24px; background: #0284c7; color: white; text-decoration: none; border-radius: 8px;">
            Vollständige Insights ansehen →
        </a>
    </p>
    <p style="color: #666; font-size: 12px; margin-top: 24px;">Trip Manager · neatnode.xyz</p>
</body>
</html>`, advertiserName, tenantsHTML)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.fromAddr,
		To:      []string{toEmail},
		Subject: "Deine wöchentlichen Travel Insights – Trip Manager",
		Html:    html,
	})
	return err
}
