package bitbucket

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeamsList_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	_, err := client.Teams.List("admin")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/teams/", receivedPath)
}

func TestTeamsProfile_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"username":     "myteam",
			"display_name": "My Team",
		})
	})
	defer server.Close()

	result, err := client.Teams.Profile("myteam")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/teams/myteam/", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, "myteam", resultMap["username"])
}

func TestTeamsMembers_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	_, err := client.Teams.Members("myteam")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/teams/myteam/members", receivedPath)
}

func TestTeamsFollowers_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	_, err := client.Teams.Followers("myteam")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/teams/myteam/followers", receivedPath)
}

func TestTeamsFollowing_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	_, err := client.Teams.Following("myteam")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/teams/myteam/following", receivedPath)
}

func TestTeamsRepositories_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	_, err := client.Teams.Repositories("myteam")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/teams/myteam/repositories", receivedPath)
}

func TestTeamsProjects_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	_, err := client.Teams.Projects("myteam")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/teams/myteam/projects/", receivedPath)
}

func TestTeams_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "team not found"},
		})
	})
	defer server.Close()

	_, err := client.Teams.Profile("nonexistent")

	assert.Error(t, err)
}
