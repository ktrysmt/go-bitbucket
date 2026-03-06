package bitbucket

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoriesListForAccount_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"page":    float64(1),
			"pagelen": float64(10),
			"size":    float64(1),
			"values": []interface{}{
				map[string]interface{}{
					"slug":      "my-repo",
					"full_name": "owner/my-repo",
					"type":      "repository",
				},
			},
		})
	})
	defer server.Close()

	opts := &RepositoriesOptions{Owner: "owner"}
	result, err := client.Repositories.ListForAccount(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner", receivedPath)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "my-repo", result.Items[0].Slug)
}

func TestRepositoriesListForAccount_WithRole(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"page": float64(1), "pagelen": float64(10), "size": float64(0),
			"values": []interface{}{},
		})
	})
	defer server.Close()

	opts := &RepositoriesOptions{Owner: "owner", Role: "admin"}
	_, err := client.Repositories.ListForAccount(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "role=admin")
}

func TestRepositoriesListForAccount_WithKeyword(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"page": float64(1), "pagelen": float64(10), "size": float64(0),
			"values": []interface{}{},
		})
	})
	defer server.Close()

	keyword := "my-search"
	opts := &RepositoriesOptions{Owner: "owner", Keyword: &keyword}
	_, err := client.Repositories.ListForAccount(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "my-search")
}

func TestRepositoriesListForAccount_EmptyOwner(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	opts := &RepositoriesOptions{Owner: ""}
	_, err := client.Repositories.ListForAccount(opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "owner")
}

func TestRepositoriesListForTeam_DelegatesToListForAccount(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"page": float64(1), "pagelen": float64(10), "size": float64(0),
			"values": []interface{}{},
		})
	})
	defer server.Close()

	opts := &RepositoriesOptions{Owner: "team"}
	_, err := client.Repositories.ListForTeam(opts)

	require.NoError(t, err)
}

func TestRepositoriesListProject_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedQuery = r.URL.RawQuery
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"page": float64(1), "pagelen": float64(10), "size": float64(1),
			"values": []interface{}{
				map[string]interface{}{
					"slug": "proj-repo",
					"type": "repository",
				},
			},
		})
	})
	defer server.Close()

	opts := &RepositoriesOptions{Owner: "owner", Project: "PROJ"}
	result, err := client.Repositories.ListProject(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/", receivedPath)
	assert.Contains(t, receivedQuery, "project.key")
	assert.Len(t, result.Items, 1)
}

func TestRepositoriesListPublic_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"page": float64(1), "pagelen": float64(10), "size": float64(0),
			"values": []interface{}{},
		})
	})
	defer server.Close()

	result, err := client.Repositories.ListPublic()

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/", receivedPath)
	assert.NotNil(t, result)
}

func TestDecodeRepositories_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page":    float64(1),
		"pagelen": float64(10),
		"size":    float64(2),
		"values": []interface{}{
			map[string]interface{}{
				"slug":      "repo1",
				"full_name": "owner/repo1",
				"type":      "repository",
			},
			map[string]interface{}{
				"slug":      "repo2",
				"full_name": "owner/repo2",
				"type":      "repository",
			},
		},
	}

	result, err := decodeRepositories(response)

	require.NoError(t, err)
	assert.Equal(t, int32(1), result.Page)
	assert.Equal(t, int32(10), result.Pagelen)
	assert.Equal(t, int32(2), result.Size)
	assert.Len(t, result.Items, 2)
}

func TestDecodeRepositories_InvalidFormat(t *testing.T) {
	t.Parallel()
	_, err := decodeRepositories("invalid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Not a valid format")
}
