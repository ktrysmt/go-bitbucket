package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCommits_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"hash":    "abc123def456",
				"message": "initial commit",
				"date":    "2025-01-15T10:00:00+00:00",
				"author": map[string]interface{}{
					"raw":  "Test User <test@example.com>",
					"type": "author",
					"user": map[string]interface{}{"display_name": "Test User", "uuid": "{user-1}"},
				},
				"parents": []interface{}{
					map[string]interface{}{"hash": "parent123"},
				},
			},
		}))
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:       "owner",
		RepoSlug:    "repo",
		Branchortag: "main",
	}
	result, err := client.Repositories.Commits.GetCommits(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/commits/main", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	require.Len(t, values, 1)

	commit := values[0].(map[string]interface{})
	assert.Equal(t, "abc123def456", commit["hash"])
	assert.Equal(t, "initial commit", commit["message"])
	assert.Equal(t, "2025-01-15T10:00:00+00:00", commit["date"])
	author := commit["author"].(map[string]interface{})
	assert.Equal(t, "Test User <test@example.com>", author["raw"])
	parents := commit["parents"].([]interface{})
	assert.Len(t, parents, 1)
}

func TestGetCommits_WithIncludeExclude(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:       "owner",
		RepoSlug:    "repo",
		Branchortag: "main",
		Include:     "feature-branch",
		Exclude:     "old-branch",
	}
	_, err := client.Repositories.Commits.GetCommits(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "include=feature-branch")
	assert.Contains(t, receivedQuery, "exclude=old-branch")
}

func TestGetCommit_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"hash":    "abc123",
			"message": "test commit",
			"date":    "2025-01-15T10:00:00+00:00",
			"author": map[string]interface{}{
				"raw":  "Test User <test@example.com>",
				"type": "author",
			},
			"parents": []interface{}{
				map[string]interface{}{"hash": "parent456"},
			},
		})
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "abc123",
	}
	result, err := client.Repositories.Commits.GetCommit(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, "abc123", resultMap["hash"])
	assert.Equal(t, "test commit", resultMap["message"])
	assert.Equal(t, "2025-01-15T10:00:00+00:00", resultMap["date"])
	author := resultMap["author"].(map[string]interface{})
	assert.Equal(t, "Test User <test@example.com>", author["raw"])
	assert.Equal(t, "author", author["type"])
	parents := resultMap["parents"].([]interface{})
	require.Len(t, parents, 1)
	parent := parents[0].(map[string]interface{})
	assert.Equal(t, "parent456", parent["hash"])
}

func TestGetCommit_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "Commit not found"},
		})
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "nonexistent",
	}
	_, err := client.Repositories.Commits.GetCommit(opts)

	assert.Error(t, err)
}

func TestGetCommitComments_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"id": 1, "content": map[string]interface{}{"raw": "nice commit"}},
		}))
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "abc123",
	}
	result, err := client.Repositories.Commits.GetCommitComments(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/comments", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	require.Len(t, values, 1)
	comment := values[0].(map[string]interface{})
	assert.Equal(t, float64(1), comment["id"])
	content := comment["content"].(map[string]interface{})
	assert.Equal(t, "nice commit", content["raw"])
}

func TestGetCommitComments_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "Not found"},
		})
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "nonexistent",
	}
	_, err := client.Repositories.Commits.GetCommitComments(opts)

	assert.Error(t, err)
}

func TestGetCommitComment_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": 42, "content": map[string]interface{}{"raw": "a comment"},
		})
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:     "owner",
		RepoSlug:  "repo",
		Revision:  "abc123",
		CommentID: "42",
	}
	result, err := client.Repositories.Commits.GetCommitComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/comments/42", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, float64(42), resultMap["id"])
	content := resultMap["content"].(map[string]interface{})
	assert.Equal(t, "a comment", content["raw"])
}

func TestGetCommitStatuses_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"key": "build", "state": "SUCCESSFUL"},
		}))
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "abc123",
	}
	result, err := client.Repositories.Commits.GetCommitStatuses(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/statuses", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	require.Len(t, values, 1)
	status := values[0].(map[string]interface{})
	assert.Equal(t, "build", status["key"])
	assert.Equal(t, "SUCCESSFUL", status["state"])
}

func TestGetCommitStatus_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"key": "build-key", "state": "SUCCESSFUL",
		})
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "abc123",
	}
	result, err := client.Repositories.Commits.GetCommitStatus(opts, "build-key")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/statuses/build/build-key", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, "build-key", resultMap["key"])
	assert.Equal(t, "SUCCESSFUL", resultMap["state"])
}

func TestGetCommitStatus_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "Status not found"},
		})
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "abc123",
	}
	_, err := client.Repositories.Commits.GetCommitStatus(opts, "nonexistent")

	assert.Error(t, err)
}

func TestGiveApprove_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{"approved": true})
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "abc123",
	}
	result, err := client.Repositories.Commits.GiveApprove(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/approve", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, true, resultMap["approved"])
}

func TestRemoveApprove_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "abc123",
	}
	result, err := client.Repositories.Commits.RemoveApprove(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Nil(t, result)
}

func TestCreateCommitStatus_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &receivedBody)
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"key": "build", "state": "INPROGRESS", "url": "https://ci.example.com/build/1", "name": "Build #1",
		})
	})
	defer server.Close()

	cmo := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "abc123",
	}
	cso := &CommitStatusOptions{
		Key:   "build",
		State: "INPROGRESS",
		Url:   "https://ci.example.com/build/1",
		Name:  "Build #1",
	}
	result, err := client.Repositories.Commits.CreateCommitStatus(cmo, cso)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/statuses/build", receivedPath)

	// Verify request body contains the expected fields
	assert.Equal(t, "build", receivedBody["key"])
	assert.Equal(t, "INPROGRESS", receivedBody["state"])
	assert.Equal(t, "https://ci.example.com/build/1", receivedBody["url"])
	assert.Equal(t, "Build #1", receivedBody["name"])

	// Verify response
	resultMap := result.(map[string]interface{})
	assert.Equal(t, "build", resultMap["key"])
	assert.Equal(t, "INPROGRESS", resultMap["state"])
	assert.Equal(t, "https://ci.example.com/build/1", resultMap["url"])
}

func TestCreateCommitStatus_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{"message": "Invalid state"},
		})
	})
	defer server.Close()

	cmo := &CommitsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Revision: "abc123",
	}
	cso := &CommitStatusOptions{
		Key:   "build",
		State: "INVALID",
	}
	_, err := client.Repositories.Commits.CreateCommitStatus(cmo, cso)

	assert.Error(t, err)
}

func TestBuildCommitsQuery(t *testing.T) {
	t.Parallel()
	commits := &Commits{}

	tests := []struct {
		name     string
		include  string
		exclude  string
		contains []string
		empty    bool
	}{
		{"both include and exclude", "feat", "old", []string{"include=feat", "exclude=old"}, false},
		{"include only", "feat", "", []string{"include=feat"}, false},
		{"exclude only", "", "old", []string{"exclude=old"}, false},
		{"neither", "", "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := commits.buildCommitsQuery(tt.include, tt.exclude)
			if tt.empty {
				assert.Empty(t, result)
			} else {
				for _, s := range tt.contains {
					assert.Contains(t, result, s)
				}
			}
		})
	}
}

func TestGetCommits_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"error": map[string]interface{}{"message": "unauthorized"},
		})
	})
	defer server.Close()

	opts := &CommitsOptions{
		Owner:       "owner",
		RepoSlug:    "repo",
		Branchortag: "main",
	}
	_, err := client.Repositories.Commits.GetCommits(opts)

	assert.Error(t, err)
}
