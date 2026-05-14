package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UsersClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUsersClient(baseURL string) *UsersClient {
	return &UsersClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *UsersClient) GetMe(ctx context.Context, token string) (*UserResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/users/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call users service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("users service returned %d", resp.StatusCode)
	}

	var user UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &user, nil
}
