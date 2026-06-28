package tenant

import "github.com/ShalArl/trip-manager/backend/shared/email"

type EmailService = email.Service

func NewEmailService(apiKey string) *EmailService {
	return email.NewService(apiKey)
}
