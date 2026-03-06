package bitbucket

import (
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
			map[string]interface{}{"hash": "abc123", "message": "initial commit"},
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
	assert.Len(t, values, 1)
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
			"author":  map[string]interface{}{"raw": "Test User <test@example.com>"},
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
	_, err := client.Repositories.Commits.GetCommitComments(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/comments", receivedPath)
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
	_, err := client.Repositories.Commits.GetCommitComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/comments/42", receivedPath)
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
	_, err := client.Repositories.Commits.GetCommitStatuses(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/statuses", receivedPath)
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
	_, err := client.Repositories.Commits.GetCommitStatus(opts, "build-key")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/statuses/build/build-key", receivedPath)
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
	_, err := client.Repositories.Commits.GiveApprove(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/approve", receivedPath)
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
	_, err := client.Repositories.Commits.RemoveApprove(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
}

func TestCreateCommitStatus_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"key": "build", "state": "INPROGRESS",
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
	_, err := client.Repositories.Commits.CreateCommitStatus(cmo, cso)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/statuses/build", receivedPath)
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
