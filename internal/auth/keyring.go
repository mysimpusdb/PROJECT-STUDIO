package auth

import (
	"log"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "PROJECT_STUDIO_APP"
	tokenKey    = "auth_token"
	userIDKey   = "auth_user_id"
)

// SaveToken securely saves the session token in the OS keyring
func SaveToken(userID, token string) error {
	// We save the token against the user ID or a generic key
	err := keyring.Set(serviceName, tokenKey, token)
	if err != nil {
		log.Printf("Failed to save token to keyring: %v\n", err)
		return err
	}
	
	err = keyring.Set(serviceName, userIDKey, userID)
	if err != nil {
		log.Printf("Failed to save user ID to keyring: %v\n", err)
		return err
	}
	return nil
}

// GetToken retrieves the stored token from the OS keyring
func GetToken() (string, string, error) {
	token, err := keyring.Get(serviceName, tokenKey)
	if err != nil {
		return "", "", err
	}
	
	userID, err := keyring.Get(serviceName, userIDKey)
	if err != nil {
		return token, "", err
	}
	
	return token, userID, nil
}

// DeleteToken removes the token from the OS keyring
func DeleteToken() error {
	err1 := keyring.Delete(serviceName, tokenKey)
	err2 := keyring.Delete(serviceName, userIDKey)
	
	if err1 != nil && err1 != keyring.ErrNotFound {
		return err1
	}
	if err2 != nil && err2 != keyring.ErrNotFound {
		return err2
	}
	return nil
}
