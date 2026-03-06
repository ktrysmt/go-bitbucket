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
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"page": float64(1), "pagelen": float64(10), "size": float64(1),
			"values": []interface{}{
				map[string]interface{}{
					"slug":      "team-repo",
					"full_name": "team/team-repo",
					"type":      "repository",
				},
			},
		})
	})
	defer server.Close()

	opts := &RepositoriesOptions{Owner: "team"}
	result, err := client.Repositories.ListForTeam(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/team", receivedPath)
	require.Len(t, result.Items, 1)
	assert.Equal(t, "team-repo", result.Items[0].Slug)
	assert.Equal(t, "team/team-repo", result.Items[0].Full_name)
	assert.Equal(t, int32(1), result.Page)
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
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"page": float64(1), "pagelen": float64(20), "size": float64(2),
			"values": []interface{}{
				map[string]interface{}{
					"slug":        "public-repo1",
					"full_name":   "community/public-repo1",
					"description": "A public repository",
					"is_private":  false,
					"fork_policy": "allow_forks",
					"language":    "go",
					"type":        "repository",
				},
				map[string]interface{}{
					"slug":        "public-repo2",
					"full_name":   "community/public-repo2",
					"description": "Another public repo",
					"is_private":  false,
					"fork_policy": "no_public_forks",
					"language":    "python",
					"type":        "repository",
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Repositories.ListPublic()

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/", receivedPath)
	require.NotNil(t, result)
	assert.Equal(t, int32(1), result.Page)
	assert.Equal(t, int32(20), result.Pagelen)
	assert.Equal(t, int32(2), result.Size)
	require.Len(t, result.Items, 2)
	assert.Equal(t, "public-repo1", result.Items[0].Slug)
	assert.Equal(t, "community/public-repo1", result.Items[0].Full_name)
	assert.Equal(t, "A public repository", result.Items[0].Description)
	assert.Equal(t, false, result.Items[0].Is_private)
	assert.Equal(t, "allow_forks", result.Items[0].Fork_policy)
	assert.Equal(t, "go", result.Items[0].Language)
}

func TestRepositoriesListPublic_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "internal server error",
			},
		})
	})
	defer server.Close()

	result, err := client.Repositories.ListPublic()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDecodeRepositories_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page":    float64(1),
		"pagelen": float64(10),
		"size":    float64(2),
		"values": []interface{}{
			map[string]interface{}{
				"slug":        "repo1",
				"full_name":   "owner/repo1",
				"description": "First repository",
				"is_private":  true,
				"fork_policy": "no_forks",
				"language":    "go",
				"type":        "repository",
			},
			map[string]interface{}{
				"slug":        "repo2",
				"full_name":   "owner/repo2",
				"description": "Second repository",
				"is_private":  false,
				"fork_policy": "allow_forks",
				"language":    "python",
				"type":        "repository",
			},
		},
	}

	result, err := decodeRepositories(response)

	require.NoError(t, err)
	assert.Equal(t, int32(1), result.Page)
	assert.Equal(t, int32(10), result.Pagelen)
	assert.Equal(t, int32(2), result.Size)
	require.Len(t, result.Items, 2)

	repo1 := result.Items[0]
	assert.Equal(t, "repo1", repo1.Slug)
	assert.Equal(t, "owner/repo1", repo1.Full_name)
	assert.Equal(t, "First repository", repo1.Description)
	assert.Equal(t, true, repo1.Is_private)
	assert.Equal(t, "no_forks", repo1.Fork_policy)
	assert.Equal(t, "go", repo1.Language)

	repo2 := result.Items[1]
	assert.Equal(t, "repo2", repo2.Slug)
	assert.Equal(t, "owner/repo2", repo2.Full_name)
	assert.Equal(t, "Second repository", repo2.Description)
	assert.Equal(t, false, repo2.Is_private)
	assert.Equal(t, "allow_forks", repo2.Fork_policy)
	assert.Equal(t, "python", repo2.Language)
}

func TestDecodeRepositories_InvalidFormat(t *testing.T) {
	t.Parallel()
	_, err := decodeRepositories("invalid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Not a valid format")
}
