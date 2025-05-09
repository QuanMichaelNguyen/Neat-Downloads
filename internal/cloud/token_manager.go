package cloud

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type TokenManager struct {
	tokenPath string
}

func NewTokenManager() *TokenManager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	tokenPath := filepath.Join(homeDir, ".neat-downloads", "dropbox_tokens.json")
	return &TokenManager{
		tokenPath: tokenPath,
	}
}

func (tm *TokenManager) SaveTokens(accessToken, refreshToken string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(tm.tokenPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create token directory: %v", err)
	}

	tokenData := TokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(4 * time.Hour), // Dropbox tokens typically expire in 4 hours
	}

	data, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %v", err)
	}

	if err := os.WriteFile(tm.tokenPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %v", err)
	}

	return nil
}

func (tm *TokenManager) LoadTokens() (*TokenData, error) {
	data, err := os.ReadFile(tm.tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read token file: %v", err)
	}

	var tokenData TokenData
	if err := json.Unmarshal(data, &tokenData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %v", err)
	}

	return &tokenData, nil
}

func (tm *TokenManager) ClearTokens() error {
	if err := os.Remove(tm.tokenPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove token file: %v", err)
	}
	return nil
} 