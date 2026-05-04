package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TokenValidationResponse is the response from auth service
type TokenValidationResponse struct {
	Valid  bool                   `json:"valid"`
	UserID string                 `json:"userId,omitempty"`
	Email  string                 `json:"email,omitempty"`
	Claims map[string]interface{} `json:"claims,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

// Client is the auth service client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new auth service client
func NewClient(authServiceURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(authServiceURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * 1000000000, // 10 seconds
		},
	}
}

// ValidateToken validates a token by sending it to auth service
func (c *Client) ValidateToken(ctx context.Context, token string) (*TokenValidationResponse, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	body := map[string]string{"token": token}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/validate-token", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(resp.Body)

	var result TokenValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ValidateTokenFromHeader validates a token from Authorization header
func (c *Client) ValidateTokenFromHeader(ctx context.Context, authHeader string) (*TokenValidationResponse, error) {
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is required")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/validate-token", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", authHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call auth service: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(resp.Body)

	var result TokenValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ValidateBearerToken extracts and validates a Bearer token
func (c *Client) ValidateBearerToken(ctx context.Context, authHeader string) (*TokenValidationResponse, error) {
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader && authHeader != "" {
		// No "Bearer " prefix, try as-is
		return c.ValidateToken(ctx, authHeader)
	}
	return c.ValidateToken(ctx, token)
}

// GetUserIDFromHeader validates header and returns UserID
func (c *Client) GetUserIDFromHeader(ctx context.Context, authHeader string) (string, error) {
	result, err := c.ValidateBearerToken(ctx, authHeader)
	if err != nil {
		return "", err
	}

	if !result.Valid {
		return "", fmt.Errorf("invalid token: %s", result.Error)
	}

	return result.UserID, nil
}

// HealthCheck checks if auth service is available
func (c *Client) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth service unreachable: %w", err)
	}
	defer func(r io.Reader) {
		_, err := io.ReadAll(r)
		if err != nil {
			fmt.Println("failed to read response body")
		}
	}(resp.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	return nil
}
