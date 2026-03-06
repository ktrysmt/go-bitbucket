package bitbucket

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPipelinesList_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"uuid": "{pipe-1}", "state": map[string]interface{}{"name": "COMPLETED"}},
		}))
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.Pipelines.List(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pipelines/", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
}

func TestPipelinesList_WithQueryAndSort(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &PipelinesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Query:    "target.branch=\"main\"",
		Sort:     "-created_on",
		Page:     2,
	}
	_, err := client.Repositories.Pipelines.List(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "sort=-created_on")
	assert.Contains(t, receivedQuery, "page=2")
}

func TestPipelinesGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"uuid": "{pipe-uuid}", "build_number": 42,
		})
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "{pipe-uuid}"}
	_, err := client.Repositories.Pipelines.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pipelines/{pipe-uuid}", receivedPath)
}

func TestPipelinesListSteps_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"uuid": "{step-1}"},
		}))
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "{pipe-uuid}"}
	_, err := client.Repositories.Pipelines.ListSteps(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pipelines/{pipe-uuid}/steps/", receivedPath)
}

func TestPipelinesListSteps_WithQueryAndSort(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &PipelinesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		IDOrUuid: "{pipe-uuid}",
		Query:    "state.name=\"COMPLETED\"",
		Sort:     "started_on",
		Page:     3,
	}
	_, err := client.Repositories.Pipelines.ListSteps(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "sort=started_on")
	assert.Contains(t, receivedQuery, "page=3")
}

func TestPipelinesGetStep_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"uuid": "{step-uuid}",
		})
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "{pipe-uuid}", StepUuid: "{step-uuid}"}
	_, err := client.Repositories.Pipelines.GetStep(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pipelines/{pipe-uuid}/steps/{step-uuid}", receivedPath)
}

func TestPipelinesGetLog_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("build log line 1\nbuild log line 2\n"))
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "{pipe-uuid}", StepUuid: "{step-uuid}"}
	logContent, err := client.Repositories.Pipelines.GetLog(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pipelines/{pipe-uuid}/steps/{step-uuid}/log", receivedPath)
	assert.Contains(t, logContent, "build log line 1")
	assert.Contains(t, logContent, "build log line 2")
}

func TestPipelinesGetLog_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "bad", StepUuid: "bad"}
	_, err := client.Repositories.Pipelines.GetLog(opts)

	assert.Error(t, err)
}
