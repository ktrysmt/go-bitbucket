package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHKeysCreate_Success(t *testing.T) {
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
			"uuid":    "{key-uuid}",
			"label":   "my-key",
			"key":     "ssh-rsa AAAA...",
			"comment": "",
		})
	})
	defer server.Close()

	opts := &SSHKeyOptions{
		Owner: "testuser",
		Label: "my-key",
		Key:   "ssh-rsa AAAA...",
	}
	key, err := client.Users.SSHKeys.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/users/testuser/ssh-keys", receivedPath)
	assert.Equal(t, "my-key", receivedBody["label"])
	assert.Equal(t, "{key-uuid}", key.Uuid)
	assert.Equal(t, "my-key", key.Label)
}

func TestSSHKeysGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"uuid":  "{key-uuid}",
			"label": "my-key",
			"key":   "ssh-rsa AAAA...",
		})
	})
	defer server.Close()

	opts := &SSHKeyOptions{Owner: "testuser", Uuid: "{key-uuid}"}
	key, err := client.Users.SSHKeys.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/users/testuser/ssh-keys/{key-uuid}", receivedPath)
	assert.Equal(t, "{key-uuid}", key.Uuid)
}

func TestSSHKeysDelete_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &SSHKeyOptions{Owner: "testuser", Uuid: "{key-uuid}"}
	_, err := client.Users.SSHKeys.Delete(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/users/testuser/ssh-keys/{key-uuid}", receivedPath)
}

func TestSSHKeysGet_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &SSHKeyOptions{Owner: "testuser", Uuid: "bad-uuid"}
	_, err := client.Users.SSHKeys.Get(opts)

	assert.Error(t, err)
}

func TestDecodeSSHKey_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"uuid":  "{key-uuid}",
		"label": "my-key",
		"key":   "ssh-rsa AAAA...",
	}

	key, err := decodeSSHKey(response)

	require.NoError(t, err)
	assert.Equal(t, "{key-uuid}", key.Uuid)
	assert.Equal(t, "my-key", key.Label)
}

func TestDecodeSSHKey_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type": "error",
		"error": map[string]interface{}{
			"message": "not found",
		},
	}

	_, err := decodeSSHKey(response)

	assert.Error(t, err)
}

func TestDecodeSSHKeys_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page":    float64(1),
		"pagelen": float64(10),
		"size":    float64(2),
		"values": []interface{}{
			map[string]interface{}{"uuid": "{key-1}", "label": "key1", "key": "ssh-rsa A"},
			map[string]interface{}{"uuid": "{key-2}", "label": "key2", "key": "ssh-rsa B"},
		},
	}

	result, err := decodeSSHKeys(response)

	require.NoError(t, err)
	assert.Equal(t, int32(1), result.Page)
	assert.Equal(t, int32(2), result.Size)
	assert.Len(t, result.Items, 2)
}

func TestDecodeSSHKeys_InvalidFormat(t *testing.T) {
	t.Parallel()
	_, err := decodeSSHKeys("invalid")

	assert.Error(t, err)
}

func TestBuildSSHKeysBody(t *testing.T) {
	t.Parallel()
	opts := &SSHKeyOptions{
		Label: "my-key",
		Key:   "ssh-rsa AAAA...",
	}

	data, err := buildSSHKeysBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "my-key", body["label"])
	assert.Equal(t, "ssh-rsa AAAA...", body["key"])
}
