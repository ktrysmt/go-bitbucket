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
	var receivedQuery string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedQuery = r.URL.RawQuery
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"pagelen": 10,
			"size":    2,
			"values": []interface{}{
				map[string]interface{}{
					"username":     "team-alpha",
					"display_name": "Team Alpha",
					"uuid":         "{team-alpha-uuid}",
					"type":         "team",
				},
				map[string]interface{}{
					"username":     "team-beta",
					"display_name": "Team Beta",
					"uuid":         "{team-beta-uuid}",
					"type":         "team",
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Teams.List("admin")

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/teams/", receivedPath)
	assert.Contains(t, receivedQuery, "role=admin")
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 2)
	firstTeam := values[0].(map[string]interface{})
	assert.Equal(t, "team-alpha", firstTeam["username"])
	assert.Equal(t, "Team Alpha", firstTeam["display_name"])
	assert.Equal(t, "{team-alpha-uuid}", firstTeam["uuid"])
	assert.Equal(t, "team", firstTeam["type"])
	secondTeam := values[1].(map[string]interface{})
	assert.Equal(t, "team-beta", secondTeam["username"])
	assert.Equal(t, float64(2), resultMap["size"])
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
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"pagelen": 10,
			"size":    2,
			"values": []interface{}{
				map[string]interface{}{
					"username":     "member1",
					"display_name": "Member One",
					"uuid":         "{member1-uuid}",
					"type":         "user",
					"account_id":   "acc-001",
				},
				map[string]interface{}{
					"username":     "member2",
					"display_name": "Member Two",
					"uuid":         "{member2-uuid}",
					"type":         "user",
					"account_id":   "acc-002",
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Teams.Members("myteam")

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/teams/myteam/members", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 2)
	member1 := values[0].(map[string]interface{})
	assert.Equal(t, "member1", member1["username"])
	assert.Equal(t, "Member One", member1["display_name"])
	assert.Equal(t, "user", member1["type"])
	member2 := values[1].(map[string]interface{})
	assert.Equal(t, "member2", member2["username"])
	assert.Equal(t, "acc-002", member2["account_id"])
}

func TestTeamsFollowers_Success(t *testing.T) {
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
					"username":     "follower1",
					"display_name": "Follower One",
					"uuid":         "{follower1-uuid}",
					"type":         "user",
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Teams.Followers("myteam")

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/teams/myteam/followers", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	follower := values[0].(map[string]interface{})
	assert.Equal(t, "follower1", follower["username"])
	assert.Equal(t, "Follower One", follower["display_name"])
	assert.Equal(t, float64(1), resultMap["size"])
}

func TestTeamsFollowing_Success(t *testing.T) {
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
					"username":     "followed-team",
					"display_name": "Followed Team",
					"uuid":         "{followed-uuid}",
					"type":         "team",
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Teams.Following("myteam")

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/teams/myteam/following", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	followed := values[0].(map[string]interface{})
	assert.Equal(t, "followed-team", followed["username"])
	assert.Equal(t, "Followed Team", followed["display_name"])
	assert.Equal(t, "team", followed["type"])
}

func TestTeamsRepositories_Success(t *testing.T) {
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
					"slug":      "my-repo",
					"full_name": "myteam/my-repo",
					"name":      "My Repo",
					"uuid":      "{repo-uuid}",
					"type":      "repository",
					"scm":       "git",
					"is_private": true,
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Teams.Repositories("myteam")

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/teams/myteam/repositories", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	repo := values[0].(map[string]interface{})
	assert.Equal(t, "my-repo", repo["slug"])
	assert.Equal(t, "myteam/my-repo", repo["full_name"])
	assert.Equal(t, "repository", repo["type"])
	assert.Equal(t, "git", repo["scm"])
	assert.Equal(t, true, repo["is_private"])
}

func TestTeamsProjects_Success(t *testing.T) {
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
					"key":         "PROJ",
					"name":        "My Project",
					"uuid":        "{project-uuid}",
					"type":        "project",
					"description": "A test project",
					"is_private":  false,
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Teams.Projects("myteam")

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/teams/myteam/projects/", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	project := values[0].(map[string]interface{})
	assert.Equal(t, "PROJ", project["key"])
	assert.Equal(t, "My Project", project["name"])
	assert.Equal(t, "project", project["type"])
	assert.Equal(t, "A test project", project["description"])
	assert.Equal(t, false, project["is_private"])
}

func TestTeams_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "team not found"},
		})
	})
	defer server.Close()

	result, err := client.Teams.Profile("nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "404")
}

func TestTeamsList_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	result, err := client.Teams.List("admin")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "403")
}

func TestTeamsMembers_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{"message": "internal server error"},
		})
	})
	defer server.Close()

	result, err := client.Teams.Members("myteam")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "500")
}
