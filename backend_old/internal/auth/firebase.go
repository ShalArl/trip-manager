package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type FirebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(ctx context.Context, projectID string) (*FirebaseAuth, error) {
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase middleware client: %w", err)
	}

	return &FirebaseAuth{client: client}, nil
}

// VerifyToken validate Firebase-ID-Token and returns verified claims.
func (f *FirebaseAuth) VerifyToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := f.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}
	return token, nil
}

// DeleteUser deletes a user from firebase.
func (f *FirebaseAuth) DeleteUser(ctx context.Context, firebaseUID string) error {
	return f.client.DeleteUser(ctx, firebaseUID)
}
