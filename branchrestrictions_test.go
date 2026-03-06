package bitbucket

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranchRestrictionsGets_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"id": 1, "kind": "push", "pattern": "main"},
		}))
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.BranchRestrictions.Gets(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/branch-restrictions", receivedPath)
}

func TestBranchRestrictionsCreate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"id": 1, "kind": "push", "pattern": "main",
		})
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Kind:     "push",
		Pattern:  "main",
	}
	result, err := client.Repositories.BranchRestrictions.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "push", result.Kind)
	assert.Equal(t, "main", result.Pattern)
}

func TestBranchRestrictionsGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": 1, "kind": "push", "pattern": "main",
		})
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.BranchRestrictions.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/branch-restrictions/1", receivedPath)
	assert.Equal(t, "push", result.Kind)
}

func TestBranchRestrictionsUpdate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": 1, "kind": "push", "pattern": "develop",
		})
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		ID:       "1",
		Kind:     "push",
		Pattern:  "develop",
	}
	_, err := client.Repositories.BranchRestrictions.Update(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
}

func TestBranchRestrictionsDelete_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.BranchRestrictions.Delete(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
}

func TestBuildBranchRestrictionsBody(t *testing.T) {
	t.Parallel()
	br := &BranchRestrictions{}
	opts := &BranchRestrictionsOptions{
		Kind:    "push",
		Pattern: "main",
		Users:   []string{"user1", "user2"},
		Groups:  map[string]string{"group1": "group1"},
	}

	data, err := br.buildBranchRestrictionsBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	err = json.Unmarshal([]byte(data), &body)
	require.NoError(t, err)

	assert.Equal(t, "push", body["kind"])
	assert.Equal(t, "main", body["pattern"])
	users := body["users"].([]interface{})
	assert.Len(t, users, 2)
}

func TestDecodeBranchRestriction_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"id":      float64(1),
		"kind":    "push",
		"pattern": "main",
	}

	result, err := decodeBranchRestriction(response)

	require.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "push", result.Kind)
	assert.Equal(t, "main", result.Pattern)
}

func TestDecodeBranchRestriction_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type": "error",
		"error": map[string]interface{}{
			"message": "not found",
		},
	}

	_, err := decodeBranchRestriction(response)

	assert.Error(t, err)
}
