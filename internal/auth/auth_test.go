package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
)

func TestKeyringFunctions(t *testing.T) {
	// Initialize mock keyring for testing
	keyring.MockInit()

	userID := "user-123"
	token := "secure-token-xyz"

	t.Run("SaveToken", func(t *testing.T) {
		err := SaveToken(userID, token)
		assert.NoError(t, err)
	})

	t.Run("GetToken", func(t *testing.T) {
		gotToken, gotUserID, err := GetToken()
		assert.NoError(t, err)
		assert.Equal(t, token, gotToken)
		assert.Equal(t, userID, gotUserID)
	})

	t.Run("DeleteToken", func(t *testing.T) {
		err := DeleteToken()
		assert.NoError(t, err)

		// Verify it's deleted
		_, _, err = GetToken()
		assert.Error(t, err)
	})
}
