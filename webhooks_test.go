package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhooksList_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{
				map[string]interface{}{
					"uuid":        "{hook-1}",
					"description": "webhook 1",
					"url":         "https://example.com/hook1",
					"active":      true,
					"events":      []interface{}{"repo:push"},
				},
			},
		})
	})
	defer server.Close()

	opts := &WebhooksOptions{Owner: "owner", RepoSlug: "repo"}
	webhooks, err := client.Repositories.Webhooks.List(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/hooks/", receivedPath)
	assert.Len(t, webhooks, 1)
	assert.Equal(t, "{hook-1}", webhooks[0].Uuid)
	assert.Equal(t, "https://example.com/hook1", webhooks[0].Url)
}

func TestWebhooksGets_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &WebhooksOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Webhooks.Gets(opts)

	require.NoError(t, err)
}

func TestWebhooksCreate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"uuid":        "{new-hook}",
			"description": "my webhook",
			"url":         "https://example.com/hook",
			"active":      true,
			"events":      []interface{}{"repo:push"},
		})
	})
	defer server.Close()

	opts := &WebhooksOptions{
		Owner:       "owner",
		RepoSlug:    "repo",
		Description: "my webhook",
		Url:         "https://example.com/hook",
		Active:      true,
		Events:      []string{"repo:push"},
	}
	webhook, err := client.Repositories.Webhooks.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "{new-hook}", webhook.Uuid)
	assert.Equal(t, "my webhook", receivedBody["description"])
	assert.Equal(t, "https://example.com/hook", receivedBody["url"])
}

func TestWebhooksGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"uuid":        "{hook-uuid}",
			"description": "webhook",
			"url":         "https://example.com/hook",
			"active":      true,
		})
	})
	defer server.Close()

	opts := &WebhooksOptions{Owner: "owner", RepoSlug: "repo", Uuid: "{hook-uuid}"}
	webhook, err := client.Repositories.Webhooks.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/hooks/{hook-uuid}", receivedPath)
	assert.Equal(t, "{hook-uuid}", webhook.Uuid)
}

func TestWebhooksUpdate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"uuid":        "{hook-uuid}",
			"description": "updated",
			"url":         "https://example.com/hook-new",
		})
	})
	defer server.Close()

	opts := &WebhooksOptions{
		Owner:       "owner",
		RepoSlug:    "repo",
		Uuid:        "{hook-uuid}",
		Description: "updated",
		Url:         "https://example.com/hook-new",
		Events:      []string{"repo:push"},
	}
	webhook, err := client.Repositories.Webhooks.Update(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "updated", webhook.Description)
}

func TestWebhooksDelete_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &WebhooksOptions{Owner: "owner", RepoSlug: "repo", Uuid: "{hook-uuid}"}
	_, err := client.Repositories.Webhooks.Delete(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/hooks/{hook-uuid}", receivedPath)
}

func TestDecodeWebhook_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"uuid":        "{hook}",
		"description": "test",
		"url":         "https://example.com",
		"active":      true,
		"events":      []interface{}{"repo:push"},
	}

	webhook, err := decodeWebhook(response)

	require.NoError(t, err)
	assert.Equal(t, "{hook}", webhook.Uuid)
	assert.Equal(t, "test", webhook.Description)
	assert.True(t, webhook.Active)
}

func TestDecodeWebhook_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type": "error",
		"error": map[string]interface{}{
			"message": "webhook not found",
		},
	}

	_, err := decodeWebhook(response)

	assert.Error(t, err)
}

func TestDecodeWebhooks_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"values": []interface{}{
			map[string]interface{}{
				"uuid":        "{hook-1}",
				"description": "hook 1",
				"url":         "https://example.com/1",
			},
			map[string]interface{}{
				"uuid":        "{hook-2}",
				"description": "hook 2",
				"url":         "https://example.com/2",
			},
		},
	}

	webhooks, err := decodeWebhooks(response)

	require.NoError(t, err)
	assert.Len(t, webhooks, 2)
	assert.Equal(t, "{hook-1}", webhooks[0].Uuid)
	assert.Equal(t, "{hook-2}", webhooks[1].Uuid)
}

func TestBuildWebhooksBody(t *testing.T) {
	t.Parallel()
	webhooks := &Webhooks{}
	opts := &WebhooksOptions{
		Description: "test hook",
		Url:         "https://example.com/hook",
		Active:      true,
		Secret:      "mysecret",
		Events:      []string{"repo:push", "pullrequest:created"},
	}

	data, err := webhooks.buildWebhooksBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))

	assert.Equal(t, "test hook", body["description"])
	assert.Equal(t, "https://example.com/hook", body["url"])
	assert.Equal(t, true, body["active"])
	assert.Equal(t, "mysecret", body["secret"])
	events := body["events"].([]interface{})
	assert.Len(t, events, 2)
}
