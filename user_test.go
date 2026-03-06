package bitbucket

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserProfile_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"type":         "user",
			"username":     "currentuser",
			"display_name": "Current User",
			"account_id":   "acc-123",
			"uuid":         "{user-uuid}",
		})
	})
	defer server.Close()

	user, err := client.User.Profile()

	require.NoError(t, err)
	assert.Equal(t, "/2.0/user", receivedPath)
	assert.Equal(t, "currentuser", user.Username)
	assert.Equal(t, "Current User", user.DisplayName)
}

func TestUserProfile_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"error": map[string]interface{}{"message": "unauthorized"},
		})
	})
	defer server.Close()

	_, err := client.User.Profile()

	assert.Error(t, err)
}

func TestUserEmails_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{
				map[string]interface{}{
					"email":      "user@example.com",
					"is_primary": true,
				},
			},
		})
	})
	defer server.Close()

	result, err := client.User.Emails()

	require.NoError(t, err)
	assert.Equal(t, "/2.0/user/emails", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
}

func TestDecodeUser_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type":           "user",
		"username":       "testuser",
		"display_name":   "Test User",
		"account_id":     "123",
		"account_status": "active",
		"uuid":           "{uuid}",
		"nickname":       "tester",
	}

	user, err := decodeUser(response)

	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "Test User", user.DisplayName)
	assert.Equal(t, "123", user.AccountId)
	assert.Equal(t, "active", user.AccountStatus)
}

func TestDecodeUser_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type": "error",
		"error": map[string]interface{}{
			"message": "user not found",
		},
	}

	_, err := decodeUser(response)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}
