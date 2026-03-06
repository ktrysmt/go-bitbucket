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
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"type":            "user",
			"username":        "currentuser",
			"display_name":    "Current User",
			"account_id":      "acc-123",
			"account_status":  "active",
			"uuid":            "{user-uuid}",
			"nickname":        "current",
			"has_2fa_enabled": true,
			"created_on":      "2019-06-01T12:00:00.000000+00:00",
			"links": map[string]interface{}{
				"self": map[string]interface{}{
					"href": "https://api.bitbucket.org/2.0/users/currentuser",
				},
			},
		})
	})
	defer server.Close()

	user, err := client.User.Profile()

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/user", receivedPath)
	assert.Equal(t, "currentuser", user.Username)
	assert.Equal(t, "Current User", user.DisplayName)
	assert.Equal(t, "acc-123", user.AccountId)
	assert.Equal(t, "active", user.AccountStatus)
	assert.Equal(t, "{user-uuid}", user.Uuid)
	assert.Equal(t, "current", user.Nickname)
	assert.True(t, user.Has2faEnabled)
	assert.Equal(t, "2019-06-01T12:00:00.000000+00:00", user.CreatedOn)
	assert.NotNil(t, user.Links)
}

func TestUserProfile_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"error": map[string]interface{}{"message": "unauthorized"},
		})
	})
	defer server.Close()

	user, err := client.User.Profile()

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "401")
}

func TestUserProfile_ServerError(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{"message": "internal server error"},
		})
	})
	defer server.Close()

	user, err := client.User.Profile()

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "500")
}

func TestUserEmails_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"pagelen": 10,
			"size":    2,
			"values": []interface{}{
				map[string]interface{}{
					"email":        "primary@example.com",
					"is_primary":   true,
					"is_confirmed": true,
					"type":         "email",
					"links": map[string]interface{}{
						"self": map[string]interface{}{
							"href": "https://api.bitbucket.org/2.0/user/emails/primary@example.com",
						},
					},
				},
				map[string]interface{}{
					"email":        "secondary@example.com",
					"is_primary":   false,
					"is_confirmed": true,
					"type":         "email",
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
	assert.Len(t, values, 2)
	assert.Equal(t, float64(2), resultMap["size"])

	primaryEmail := values[0].(map[string]interface{})
	assert.Equal(t, "primary@example.com", primaryEmail["email"])
	assert.Equal(t, true, primaryEmail["is_primary"])
	assert.Equal(t, true, primaryEmail["is_confirmed"])
	assert.Equal(t, "email", primaryEmail["type"])

	secondaryEmail := values[1].(map[string]interface{})
	assert.Equal(t, "secondary@example.com", secondaryEmail["email"])
	assert.Equal(t, false, secondaryEmail["is_primary"])
}

func TestDecodeUser_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type":            "user",
		"username":        "testuser",
		"display_name":    "Test User",
		"account_id":      "557058:abcdef01-2345-6789-abcd-ef0123456789",
		"account_status":  "active",
		"uuid":            "{user-uuid-1234}",
		"nickname":        "tester",
		"website":         "https://example.com",
		"created_on":      "2020-01-15T10:30:00.000000+00:00",
		"has_2fa_enabled": true,
		"links": map[string]interface{}{
			"self": map[string]interface{}{
				"href": "https://api.bitbucket.org/2.0/users/testuser",
			},
			"avatar": map[string]interface{}{
				"href": "https://avatar.example.com/testuser.png",
			},
		},
	}

	user, err := decodeUser(response)

	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "Test User", user.DisplayName)
	assert.Equal(t, "557058:abcdef01-2345-6789-abcd-ef0123456789", user.AccountId)
	assert.Equal(t, "active", user.AccountStatus)
	assert.Equal(t, "{user-uuid-1234}", user.Uuid)
	assert.Equal(t, "tester", user.Nickname)
	assert.Equal(t, "https://example.com", user.Website)
	assert.Equal(t, "2020-01-15T10:30:00.000000+00:00", user.CreatedOn)
	assert.True(t, user.Has2faEnabled)
	assert.NotNil(t, user.Links)
	selfLink := user.Links["self"].(map[string]interface{})
	assert.Equal(t, "https://api.bitbucket.org/2.0/users/testuser", selfLink["href"])
	avatarLink := user.Links["avatar"].(map[string]interface{})
	assert.Equal(t, "https://avatar.example.com/testuser.png", avatarLink["href"])
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
