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

	user, err := client.Users.Get("nonexistent")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "404")
}

func TestUsersFollowers_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"pagelen": 10,
			"size":    2,
			"values": []interface{}{
				map[string]interface{}{
					"username":     "follower1",
					"display_name": "Follower One",
					"uuid":         "{follower1-uuid}",
					"type":         "user",
					"account_id":   "acc-f1",
				},
				map[string]interface{}{
					"username":     "follower2",
					"display_name": "Follower Two",
					"uuid":         "{follower2-uuid}",
					"type":         "user",
					"account_id":   "acc-f2",
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Users.Followers("testuser")

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/users/testuser/followers", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 2)
	follower1 := values[0].(map[string]interface{})
	assert.Equal(t, "follower1", follower1["username"])
	assert.Equal(t, "Follower One", follower1["display_name"])
	assert.Equal(t, "user", follower1["type"])
	follower2 := values[1].(map[string]interface{})
	assert.Equal(t, "follower2", follower2["username"])
	assert.Equal(t, float64(2), resultMap["size"])
}

func TestUsersFollowing_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"pagelen": 10,
			"size":    1,
			"values": []interface{}{
				map[string]interface{}{
					"username":     "followed-user",
					"display_name": "Followed User",
					"uuid":         "{followed-uuid}",
					"type":         "user",
					"account_id":   "acc-followed",
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Users.Following("testuser")

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/users/testuser/following", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	followed := values[0].(map[string]interface{})
	assert.Equal(t, "followed-user", followed["username"])
	assert.Equal(t, "Followed User", followed["display_name"])
	assert.Equal(t, "acc-followed", followed["account_id"])
	assert.Equal(t, float64(1), resultMap["size"])
}

func TestUsersRepositories_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"pagelen": 10,
			"size":    1,
			"values": []interface{}{
				map[string]interface{}{
					"slug":       "my-repo",
					"full_name":  "testuser/my-repo",
					"name":       "My Repo",
					"uuid":       "{repo-uuid}",
					"type":       "repository",
					"scm":        "git",
					"is_private": false,
					"owner": map[string]interface{}{
						"username":     "testuser",
						"display_name": "Test User",
						"type":         "user",
					},
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Users.Repositories("testuser")

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/users/testuser/repositories", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	repo := values[0].(map[string]interface{})
	assert.Equal(t, "my-repo", repo["slug"])
	assert.Equal(t, "testuser/my-repo", repo["full_name"])
	assert.Equal(t, "repository", repo["type"])
	assert.Equal(t, "git", repo["scm"])
	assert.Equal(t, false, repo["is_private"])
	owner := repo["owner"].(map[string]interface{})
	assert.Equal(t, "testuser", owner["username"])
}

func TestUsersFollowers_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"error": map[string]interface{}{"message": "unauthorized"},
		})
	})
	defer server.Close()

	result, err := client.Users.Followers("testuser")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")
}
