/**
 * @file SDK authenticator for Cline provider
 * @description Implements the Authenticator interface for Cline authentication.
 * Supports two authentication modes:
 *   1. OAuth (ClineAuthenticator): Uses a refresh token from the VSCode extension.
 *   2. API Key (ClineAPIKeyAuthenticator): Uses a Cline API key directly.
 */

package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/auth/cline"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
	coreauth "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/auth"
)

// ---------------------------------------------------------------------------
// ClineAuthenticator — OAuth refresh token flow
// ---------------------------------------------------------------------------

// ClineAuthenticator implements the authentication flow for Cline accounts.
// It uses a refresh token obtained from the Cline VSCode extension to generate
// JWT access tokens for API authentication.
type ClineAuthenticator struct{}

// NewClineAuthenticator constructs a Cline authenticator.
func NewClineAuthenticator() *ClineAuthenticator {
	return &ClineAuthenticator{}
}

// Provider returns the provider identifier for this authenticator.
func (a *ClineAuthenticator) Provider() string {
	return "cline"
}

// RefreshLead returns the recommended time before token expiration to trigger a refresh.
// Cline tokens typically expire in 10 minutes, so we refresh 2 minutes before expiration.
func (a *ClineAuthenticator) RefreshLead() *time.Duration {
	d := 2 * time.Minute
	return &d
}

// Login performs the Cline authentication flow.
func (a *ClineAuthenticator) Login(ctx context.Context, cfg *config.Config, opts *LoginOptions) (*coreauth.Auth, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cliproxy auth: configuration is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if opts == nil {
		opts = &LoginOptions{}
	}

	// Get refresh token from metadata or prompt
	refreshToken := ""
	if opts.Metadata != nil {
		refreshToken = opts.Metadata["refresh_token"]
	}

	if refreshToken == "" && opts.Prompt != nil {
		fmt.Println("\nTo authenticate with Cline:")
		fmt.Println("1. Ensure you have the Cline extension installed and are logged in to VS Code.")
		fmt.Println("2. Run the included helper script to extract your refresh token:")
		fmt.Println("   $ ./get-cline-token.sh")
		fmt.Println("3. Copy the output and paste it below.")
		fmt.Println()

		var err error
		refreshToken, err = opts.Prompt("Please paste your Cline refresh token:")
		if err != nil {
			return nil, err
		}
	}

	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required for Cline authentication")
	}

	// Exchange refresh token for access token
	authSvc := cline.NewClineAuth(cfg)
	tokenData, err := authSvc.RefreshTokens(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("cline token exchange failed: %w", err)
	}

	// Create token storage
	tokenStorage := authSvc.CreateTokenStorage(tokenData)

	// Get email from metadata or prompt
	email := tokenData.Email
	if email == "" && opts.Metadata != nil {
		email = opts.Metadata["email"]
		if email == "" {
			email = opts.Metadata["alias"]
		}
	}

	if email == "" && opts.Prompt != nil {
		email, err = opts.Prompt("Please input your email address or alias for Cline:")
		if err != nil {
			return nil, err
		}
	}

	email = strings.TrimSpace(email)
	if email == "" {
		return nil, &EmailRequiredError{Prompt: "Please provide an email address or alias for Cline."}
	}

	tokenStorage.Email = email

	fileName := fmt.Sprintf("cline-%s.json", tokenStorage.Email)
	metadata := map[string]any{
		"email": tokenStorage.Email,
	}

	fmt.Println("Cline authentication successful")

	return &coreauth.Auth{
		ID:       fileName,
		Provider: a.Provider(),
		FileName: fileName,
		Storage:  tokenStorage,
		Metadata: metadata,
	}, nil
}

// ---------------------------------------------------------------------------
// ClineAPIKeyAuthenticator — direct API key flow
// ---------------------------------------------------------------------------

// ClineAPIKeyAuthenticator implements the authentication flow for Cline API keys.
// Unlike the OAuth flow, this simply stores the API key for direct use with the
// Cline API endpoint (without the workos: prefix).
type ClineAPIKeyAuthenticator struct{}

// NewClineAPIKeyAuthenticator constructs a Cline API key authenticator.
func NewClineAPIKeyAuthenticator() *ClineAPIKeyAuthenticator {
	return &ClineAPIKeyAuthenticator{}
}

// Provider returns the provider identifier.
func (a *ClineAPIKeyAuthenticator) Provider() string {
	return "cline-api"
}

// RefreshLead returns nil — API keys don't expire and don't need refresh.
func (a *ClineAPIKeyAuthenticator) RefreshLead() *time.Duration {
	return nil
}

// Login prompts the user for a Cline API key and stores it.
func (a *ClineAPIKeyAuthenticator) Login(ctx context.Context, cfg *config.Config, opts *LoginOptions) (*coreauth.Auth, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cliproxy auth: configuration is required")
	}
	if opts == nil {
		opts = &LoginOptions{}
	}

	apiKey := ""
	if opts.Metadata != nil {
		apiKey = opts.Metadata["api_key"]
	}

	if apiKey == "" && opts.Prompt != nil {
		fmt.Println("\nTo authenticate with a Cline API key:")
		fmt.Println("1. Get your API key from https://cline.bot")
		fmt.Println("2. Paste it below.")
		fmt.Println()

		var err error
		apiKey, err = opts.Prompt("Please paste your Cline API key:")
		if err != nil {
			return nil, err
		}
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for Cline API key authentication")
	}

	// Get a label/alias
	label := ""
	if opts.Metadata != nil {
		label = opts.Metadata["alias"]
	}
	if label == "" && opts.Prompt != nil {
		var err error
		label, err = opts.Prompt("Please input a label or alias for this API key:")
		if err != nil {
			return nil, err
		}
	}
	label = strings.TrimSpace(label)
	if label == "" {
		return nil, &EmailRequiredError{Prompt: "Please provide a label or alias for the Cline API key."}
	}

	fileName := fmt.Sprintf("cline-%s.json", label)

	// Store as a simple JSON with api_key and type fields.
	// The watcher/synthesizer will pick up api_key from Attributes.
	storage := &cline.ClineTokenStorage{
		AccessToken: apiKey,
		Email:       label,
		Type:        "cline",
	}

	metadata := map[string]any{
		"email":    label,
		"api_key":  apiKey,
		"auth_kind": "api_key",
	}

	fmt.Println("Cline API key authentication successful")

	return &coreauth.Auth{
		ID:       fileName,
		Provider: "cline", // same provider so the same executor handles it
		FileName: fileName,
		Storage:  storage,
		Metadata: metadata,
		Attributes: map[string]string{
			"api_key":   apiKey,
			"auth_kind": "api_key",
		},
	}, nil
}

