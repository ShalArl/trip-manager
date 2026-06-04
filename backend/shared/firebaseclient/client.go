package firebaseclient

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type Client struct {
	auth *auth.Client
}

func New(ctx context.Context, projectID string) (*Client, error) {
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init firebase app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get firebase auth client: %w", err)
	}

	return &Client{auth: authClient}, nil
}

func (c *Client) SetCustomClaims(ctx context.Context, uid string, claims map[string]interface{}) error {
	return c.auth.SetCustomUserClaims(ctx, uid, claims)
}

func (c *Client) RevokeRefreshTokens(ctx context.Context, uid string) error {
	return c.auth.RevokeRefreshTokens(ctx, uid)
}
