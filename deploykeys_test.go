package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeployKeysCreate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"id":    float64(123),
			"label": "deploy-key",
			"key":   "ssh-rsa AAAA...",
		})
	})
	defer server.Close()

	opts := &DeployKeyOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Label:    "deploy-key",
		Key:      "ssh-rsa AAAA...",
	}
	key, err := client.Repositories.DeployKeys.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/deploy-keys", receivedPath)
	assert.Equal(t, "deploy-key", receivedBody["label"])
	assert.Equal(t, "ssh-rsa AAAA...", receivedBody["key"])
	assert.Equal(t, 123, key.Id)
	assert.Equal(t, "deploy-key", key.Label)
	assert.Equal(t, "ssh-rsa AAAA...", key.Key)
}

func TestDeployKeysCreate_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	defer server.Close()

	opts := &DeployKeyOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Label:    "bad-key",
		Key:      "invalid",
	}
	result, err := client.Repositories.DeployKeys.Create(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDeployKeysGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id":    float64(123),
			"label": "deploy-key",
			"key":   "ssh-rsa AAAA...",
		})
	})
	defer server.Close()

	opts := &DeployKeyOptions{Owner: "owner", RepoSlug: "repo", Id: 123}
	key, err := client.Repositories.DeployKeys.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/deploy-keys/123", receivedPath)
	assert.Equal(t, 123, key.Id)
}

func TestDeployKeysDelete_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &DeployKeyOptions{Owner: "owner", RepoSlug: "repo", Id: 123}
	result, err := client.Repositories.DeployKeys.Delete(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/deploy-keys/123", receivedPath)
	assert.Nil(t, result)
}

func TestDeployKeysDelete_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	opts := &DeployKeyOptions{Owner: "owner", RepoSlug: "repo", Id: 999}
	result, err := client.Repositories.DeployKeys.Delete(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDeployKeysList_Success(t *testing.T) {
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
					"id":    float64(123),
					"label": "key1",
					"key":   "ssh-rsa A",
				},
			},
		})
	})
	defer server.Close()

	opts := &DeployKeyOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.DeployKeys.List(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/deploy-keys", receivedPath)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, 123, result.Items[0].Id)
}

func TestDecodeDeployKey_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"id":      float64(42),
		"label":   "my-key",
		"key":     "ssh-rsa AAAA...",
		"comment": "test@host",
	}

	key, err := decodeDeployKey(response)

	require.NoError(t, err)
	assert.Equal(t, 42, key.Id)
	assert.Equal(t, "my-key", key.Label)
	assert.Equal(t, "ssh-rsa AAAA...", key.Key)
	assert.Equal(t, "test@host", key.Comment)
}

func TestDecodeDeployKey_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type": "error",
		"error": map[string]interface{}{
			"message": "not found",
		},
	}

	_, err := decodeDeployKey(response)

	assert.Error(t, err)
}

func TestDecodeDeployKeys_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page":    float64(1),
		"pagelen": float64(10),
		"size":    float64(2),
		"values": []interface{}{
			map[string]interface{}{"id": float64(1), "label": "key1", "key": "ssh-rsa A", "comment": "comment1"},
			map[string]interface{}{"id": float64(2), "label": "key2", "key": "ssh-rsa B", "comment": "comment2"},
		},
	}

	result, err := decodeDeployKeys(response)

	require.NoError(t, err)
	assert.Equal(t, int32(1), result.Page)
	assert.Equal(t, int32(10), result.Pagelen)
	assert.Equal(t, int32(2), result.Size)
	require.Len(t, result.Items, 2)
	// Verify individual item fields
	assert.Equal(t, 1, result.Items[0].Id)
	assert.Equal(t, "key1", result.Items[0].Label)
	assert.Equal(t, "ssh-rsa A", result.Items[0].Key)
	assert.Equal(t, "comment1", result.Items[0].Comment)
	assert.Equal(t, 2, result.Items[1].Id)
	assert.Equal(t, "key2", result.Items[1].Label)
	assert.Equal(t, "ssh-rsa B", result.Items[1].Key)
	assert.Equal(t, "comment2", result.Items[1].Comment)
}

func TestDecodeDeployKeys_InvalidFormat(t *testing.T) {
	t.Parallel()
	_, err := decodeDeployKeys("invalid")

	assert.Error(t, err)
}

func TestBuildDeployKeysBody(t *testing.T) {
	t.Parallel()
	opts := &DeployKeyOptions{
		Label: "deploy-key",
		Key:   "ssh-rsa AAAA...",
	}

	data, err := buildDeployKeysBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	err = json.Unmarshal([]byte(data), &body)
	require.NoError(t, err)
	assert.Equal(t, "deploy-key", body["label"])
	assert.Equal(t, "ssh-rsa AAAA...", body["key"])
}
