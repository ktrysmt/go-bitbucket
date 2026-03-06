package bitbucket

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"type":         "user",
			"username":     "testuser",
			"display_name": "Test User",
			"account_id":   "123456",
			"uuid":         "{abc-def}",
		})
	})
	defer server.Close()

	user, err := client.Users.Get("testuser")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/users/testuser/", receivedPath)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "Test User", user.DisplayName)
	assert.Equal(t, "123456", user.AccountId)
}

func TestUsersGet_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "user not found"},
		})
	})
	defer server.Close()

	_, err := client.Users.Get("nonexistent")

	assert.Error(t, err)
}

func TestUsersFollowers_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	_, err := client.Users.Followers("testuser")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/users/testuser/followers", receivedPath)
}

func TestUsersFollowing_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	_, err := client.Users.Following("testuser")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/users/testuser/following", receivedPath)
}

func TestUsersRepositories_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	_, err := client.Users.Repositories("testuser")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/users/testuser/repositories", receivedPath)
}
