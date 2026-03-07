/**
 * @file Cline login command implementation
 * @description Handles both Cline authentication flows:
 *   1. OAuth refresh token flow (--cline-login): uses refresh token from VSCode extension
 *   2. API key flow (--cline-api-login): uses a Cline API key directly
 */

package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
	sdkAuth "github.com/router-for-me/CLIProxyAPI/v6/sdk/auth"
	log "github.com/sirupsen/logrus"
)

// DoClineLogin handles the Cline authentication flow using the shared authentication manager.
// It prompts the user for a refresh token (exported from VSCode), exchanges it for access tokens,
// and saves the authentication credentials to the configured auth directory.
func DoClineLogin(cfg *config.Config, options *LoginOptions) {
	if options == nil {
		options = &LoginOptions{}
	}

	manager := newAuthManager()

	promptFn := options.Prompt
	if promptFn == nil {
		promptFn = func(prompt string) (string, error) {
			fmt.Println()
			fmt.Println(prompt)
			var value string
			_, err := fmt.Scanln(&value)
			return value, err
		}
	}

	authOpts := &sdkAuth.LoginOptions{
		NoBrowser: true, // Cline doesn't use browser-based OAuth
		Metadata:  map[string]string{},
		Prompt:    promptFn,
	}

	_, savedPath, err := manager.Login(context.Background(), "cline", cfg, authOpts)
	if err != nil {
		var emailErr *sdkAuth.EmailRequiredError
		if errors.As(err, &emailErr) {
			log.Error(emailErr.Error())
			return
		}
		fmt.Printf("Cline authentication failed: %v\n", err)
		return
	}

	if savedPath != "" {
		fmt.Printf("Authentication saved to %s\n", savedPath)
	}

	fmt.Println("Cline authentication successful!")
}

// DoClineAPILogin handles the Cline API key authentication flow.
// It prompts the user for a Cline API key and saves it for direct use with the Cline API.
func DoClineAPILogin(cfg *config.Config, options *LoginOptions) {
	if options == nil {
		options = &LoginOptions{}
	}

	manager := newAuthManager()

	promptFn := options.Prompt
	if promptFn == nil {
		promptFn = func(prompt string) (string, error) {
			fmt.Println()
			fmt.Println(prompt)
			var value string
			_, err := fmt.Scanln(&value)
			return value, err
		}
	}

	authOpts := &sdkAuth.LoginOptions{
		NoBrowser: true,
		Metadata:  map[string]string{"auth_mode": "api_key"},
		Prompt:    promptFn,
	}

	_, savedPath, err := manager.Login(context.Background(), "cline-api", cfg, authOpts)
	if err != nil {
		var emailErr *sdkAuth.EmailRequiredError
		if errors.As(err, &emailErr) {
			log.Error(emailErr.Error())
			return
		}
		fmt.Printf("Cline API key authentication failed: %v\n", err)
		return
	}

	if savedPath != "" {
		fmt.Printf("Authentication saved to %s\n", savedPath)
	}

	fmt.Println("Cline API key authentication successful!")
}

