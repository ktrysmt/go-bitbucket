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
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
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
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/hooks/", receivedPath)
	assert.Len(t, webhooks, 1)
	assert.Equal(t, "{hook-1}", webhooks[0].Uuid)
	assert.Equal(t, "https://example.com/hook1", webhooks[0].Url)
}

func TestWebhooksGets_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"uuid":        "{hook-1}",
				"description": "test hook",
				"url":         "https://example.com",
				"active":      true,
			},
		}))
	})
	defer server.Close()

	opts := &WebhooksOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.Webhooks.Gets(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/hooks/", receivedPath)
	require.NotNil(t, result)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "result should be a map")
	values, ok := resultMap["values"].([]interface{})
	require.True(t, ok, "result should contain values array")
	assert.Len(t, values, 1)
	firstItem := values[0].(map[string]interface{})
	assert.Equal(t, "{hook-1}", firstItem["uuid"])
}

func TestWebhooksGets_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	opts := &WebhooksOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.Webhooks.Gets(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
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
	assert.Equal(t, "my webhook", webhook.Description)
	assert.Equal(t, "https://example.com/hook", webhook.Url)
	assert.True(t, webhook.Active)
	require.Len(t, webhook.Events, 1)
	assert.Equal(t, "repo:push", webhook.Events[0])

	// Verify request body serialization
	assert.Equal(t, "my webhook", receivedBody["description"])
	assert.Equal(t, "https://example.com/hook", receivedBody["url"])
	assert.Equal(t, true, receivedBody["active"])
	events := receivedBody["events"].([]interface{})
	require.Len(t, events, 1)
	assert.Equal(t, "repo:push", events[0])
}

func TestWebhooksCreate_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	defer server.Close()

	opts := &WebhooksOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Url:      "https://example.com/hook",
		Events:   []string{"repo:push"},
	}
	result, err := client.Repositories.Webhooks.Create(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestWebhooksGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
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
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/hooks/{hook-uuid}", receivedPath)
	assert.Equal(t, "{hook-uuid}", webhook.Uuid)
}

func TestWebhooksUpdate_Success(t *testing.T) {
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
			"uuid":        "{hook-uuid}",
			"description": "updated",
			"url":         "https://example.com/hook-new",
			"active":      true,
			"events":      []interface{}{"repo:push", "issue:created"},
		})
	})
	defer server.Close()

	opts := &WebhooksOptions{
		Owner:       "owner",
		RepoSlug:    "repo",
		Uuid:        "{hook-uuid}",
		Description: "updated",
		Url:         "https://example.com/hook-new",
		Active:      true,
		Events:      []string{"repo:push", "issue:created"},
	}
	webhook, err := client.Repositories.Webhooks.Update(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/hooks/{hook-uuid}", receivedPath)
	assert.Equal(t, "updated", webhook.Description)
	assert.Equal(t, "https://example.com/hook-new", webhook.Url)
	assert.True(t, webhook.Active)
	require.Len(t, webhook.Events, 2)
	assert.Equal(t, "repo:push", webhook.Events[0])
	assert.Equal(t, "issue:created", webhook.Events[1])

	// Verify request body Events array
	events := receivedBody["events"].([]interface{})
	require.Len(t, events, 2)
	assert.Equal(t, "repo:push", events[0])
	assert.Equal(t, "issue:created", events[1])
}

func TestWebhooksUpdate_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	opts := &WebhooksOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Uuid:     "bad-uuid",
		Url:      "https://example.com",
		Events:   []string{"repo:push"},
	}
	result, err := client.Repositories.Webhooks.Update(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
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
	result, err := client.Repositories.Webhooks.Delete(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/hooks/{hook-uuid}", receivedPath)
	assert.Nil(t, result)
}

func TestWebhooksDelete_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	opts := &WebhooksOptions{Owner: "owner", RepoSlug: "repo", Uuid: "bad-uuid"}
	result, err := client.Repositories.Webhooks.Delete(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
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
	assert.Equal(t, "repo:push", events[0])
	assert.Equal(t, "pullrequest:created", events[1])
}

func TestBuildWebhooksBody_EmptyOptionalFields(t *testing.T) {
	t.Parallel()
	webhooks := &Webhooks{}
	opts := &WebhooksOptions{
		Events: []string{"repo:push"},
		// Description, Url, Secret are all empty strings -> should be excluded
		// Active is false (zero value for bool) -> still included because of the always-true condition
	}

	data, err := webhooks.buildWebhooksBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))

	// Empty string fields should be excluded from the body
	_, hasDescription := body["description"]
	assert.False(t, hasDescription, "empty description should be excluded")
	_, hasUrl := body["url"]
	assert.False(t, hasUrl, "empty url should be excluded")
	_, hasSecret := body["secret"]
	assert.False(t, hasSecret, "empty secret should be excluded")

	// Active is always included (the condition `true || false` is always true)
	assert.Equal(t, false, body["active"])

	// Events should always be present
	events := body["events"].([]interface{})
	require.Len(t, events, 1)
	assert.Equal(t, "repo:push", events[0])
}

func TestWebhooksList_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	opts := &WebhooksOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.Webhooks.List(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestWebhooksGet_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	opts := &WebhooksOptions{Owner: "owner", RepoSlug: "repo", Uuid: "bad-uuid"}
	result, err := client.Repositories.Webhooks.Get(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}
