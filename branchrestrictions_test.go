package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranchRestrictionsGets_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"id": 1, "kind": "push", "pattern": "main"},
		}))
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.BranchRestrictions.Gets(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/branch-restrictions", receivedPath)
	require.NotNil(t, result)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "result should be a map")
	values, ok := resultMap["values"].([]interface{})
	require.True(t, ok, "result should contain values array")
	assert.Len(t, values, 1)
	firstItem := values[0].(map[string]interface{})
	assert.Equal(t, float64(1), firstItem["id"])
	assert.Equal(t, "push", firstItem["kind"])
	assert.Equal(t, "main", firstItem["pattern"])
}

func TestBranchRestrictionsGets_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.BranchRestrictions.Gets(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBranchRestrictionsCreate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
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
		Users:    []string{"user1"},
		Groups:   map[string]string{"group1": "group1"},
	}
	result, err := client.Repositories.BranchRestrictions.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "push", result.Kind)
	assert.Equal(t, "main", result.Pattern)

	// Verify the serialized request body
	assert.Equal(t, "push", receivedBody["kind"])
	assert.Equal(t, "main", receivedBody["pattern"])
	users := receivedBody["users"].([]interface{})
	require.Len(t, users, 1)
	assert.Equal(t, "user1", users[0].(map[string]interface{})["username"])
	groups := receivedBody["groups"].([]interface{})
	require.Len(t, groups, 1)
	assert.Equal(t, "group1", groups[0].(map[string]interface{})["slug"])
}

func TestBranchRestrictionsCreate_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Kind:     "push",
		Pattern:  "main",
	}
	result, err := client.Repositories.BranchRestrictions.Create(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
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
	var receivedPath string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
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
	result, err := client.Repositories.BranchRestrictions.Update(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/branch-restrictions/1", receivedPath)
	assert.Equal(t, "push", receivedBody["kind"])
	assert.Equal(t, "develop", receivedBody["pattern"])
	// Update returns interface{}, but it's decoded via decodeBranchRestriction
	brResult, ok := result.(*BranchRestrictions)
	require.True(t, ok, "result should be *BranchRestrictions")
	assert.Equal(t, 1, brResult.ID)
	assert.Equal(t, "push", brResult.Kind)
	assert.Equal(t, "develop", brResult.Pattern)
}

func TestBranchRestrictionsUpdate_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		ID:       "999",
		Kind:     "push",
		Pattern:  "develop",
	}
	result, err := client.Repositories.BranchRestrictions.Update(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBranchRestrictionsDelete_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.BranchRestrictions.Delete(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/branch-restrictions/1", receivedPath)
	// Delete with 204 No Content returns nil body
	assert.Nil(t, result)
}

func TestBranchRestrictionsDelete_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{Owner: "owner", RepoSlug: "repo", ID: "999"}
	result, err := client.Repositories.BranchRestrictions.Delete(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBranchRestrictionsGet_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	opts := &BranchRestrictionsOptions{Owner: "owner", RepoSlug: "repo", ID: "999"}
	result, err := client.Repositories.BranchRestrictions.Get(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBuildBranchRestrictionsBody(t *testing.T) {
	t.Parallel()
	br := &BranchRestrictions{}
	opts := &BranchRestrictionsOptions{
		Kind:    "push",
		Pattern: "main",
		Users:   []string{"user1", "user2"},
		Groups:  map[string]string{"group1": "group1", "group2": "group2"},
		Value:   42,
	}

	data, err := br.buildBranchRestrictionsBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	err = json.Unmarshal([]byte(data), &body)
	require.NoError(t, err)

	assert.Equal(t, "push", body["kind"])
	assert.Equal(t, "main", body["pattern"])
	assert.Equal(t, float64(42), body["value"])

	// Verify Users array: string[] -> struct[] with username field
	users := body["users"].([]interface{})
	require.Len(t, users, 2)
	user1 := users[0].(map[string]interface{})
	assert.Equal(t, "user1", user1["username"])
	user2 := users[1].(map[string]interface{})
	assert.Equal(t, "user2", user2["username"])

	// Verify Groups array: map[string]string -> struct[] with slug field
	groups := body["groups"].([]interface{})
	require.Len(t, groups, 2)
	groupSlugs := []string{}
	for _, g := range groups {
		groupSlugs = append(groupSlugs, g.(map[string]interface{})["slug"].(string))
	}
	assert.Contains(t, groupSlugs, "group1")
	assert.Contains(t, groupSlugs, "group2")
}

func TestBuildBranchRestrictionsBody_Empty(t *testing.T) {
	t.Parallel()
	br := &BranchRestrictions{}
	opts := &BranchRestrictionsOptions{
		Kind:    "push",
		Pattern: "main",
	}

	data, err := br.buildBranchRestrictionsBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	err = json.Unmarshal([]byte(data), &body)
	require.NoError(t, err)

	assert.Equal(t, "branchrestriction", body["type"])
	assert.Equal(t, "push", body["kind"])
	assert.Equal(t, "main", body["pattern"])
	assert.Equal(t, "glob", body["branch_match_kind"])
	// users/groups serialize as empty arrays so the swagger schema accepts them.
	assert.Equal(t, []interface{}{}, body["users"])
	assert.Equal(t, []interface{}{}, body["groups"])
}

func TestBuildBranchRestrictionsBody_BranchingModel(t *testing.T) {
	t.Parallel()
	br := &BranchRestrictions{}
	opts := &BranchRestrictionsOptions{
		Kind:            "require_passing_builds_to_merge",
		Pattern:         "",
		BranchMatchKind: "branching_model",
		BranchType:      "production",
		Value:           2,
	}

	data, err := br.buildBranchRestrictionsBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	err = json.Unmarshal([]byte(data), &body)
	require.NoError(t, err)

	assert.Equal(t, "branching_model", body["branch_match_kind"])
	assert.Equal(t, "production", body["branch_type"])
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
