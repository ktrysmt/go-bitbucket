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
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"uuid":         "{pipe-1}",
				"build_number": 10,
				"state":        map[string]interface{}{"name": "COMPLETED"},
				"target":       map[string]interface{}{"ref_name": "main", "type": "pipeline_ref_target"},
				"trigger":      map[string]interface{}{"type": "push"},
				"created_on":   "2025-01-15T10:00:00.000000+00:00",
				"completed_on": "2025-01-15T10:05:00.000000+00:00",
			},
		}))
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.Pipelines.List(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pipelines/", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)

	pipeline := values[0].(map[string]interface{})
	assert.Equal(t, "{pipe-1}", pipeline["uuid"])
	assert.Equal(t, float64(10), pipeline["build_number"])
	state := pipeline["state"].(map[string]interface{})
	assert.Equal(t, "COMPLETED", state["name"])
	target := pipeline["target"].(map[string]interface{})
	assert.Equal(t, "main", target["ref_name"])
	assert.Equal(t, "2025-01-15T10:00:00.000000+00:00", pipeline["created_on"])
	assert.Equal(t, "2025-01-15T10:05:00.000000+00:00", pipeline["completed_on"])
}

func TestPipelinesList_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"error": map[string]interface{}{"message": "unauthorized"},
		})
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Pipelines.List(opts)

	assert.Error(t, err)
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
	result, err := client.Repositories.Pipelines.List(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "sort=-created_on")
	assert.Contains(t, receivedQuery, "page=2")
	assert.Contains(t, receivedQuery, "q=target.branch")
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 0)
}

func TestPipelinesGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"uuid":         "{pipe-uuid}",
			"build_number": 42,
			"state":        map[string]interface{}{"name": "COMPLETED", "result": map[string]interface{}{"name": "SUCCESSFUL"}},
			"target":       map[string]interface{}{"ref_name": "develop", "type": "pipeline_ref_target"},
			"trigger":      map[string]interface{}{"type": "push"},
			"created_on":   "2025-02-01T12:00:00.000000+00:00",
			"completed_on": "2025-02-01T12:10:00.000000+00:00",
		})
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "{pipe-uuid}"}
	result, err := client.Repositories.Pipelines.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pipelines/{pipe-uuid}", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, "{pipe-uuid}", resultMap["uuid"])
	assert.Equal(t, float64(42), resultMap["build_number"])
	state := resultMap["state"].(map[string]interface{})
	assert.Equal(t, "COMPLETED", state["name"])
	assert.Equal(t, "2025-02-01T12:00:00.000000+00:00", resultMap["created_on"])
	assert.Equal(t, "2025-02-01T12:10:00.000000+00:00", resultMap["completed_on"])
}

func TestPipelinesGet_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "Pipeline not found"},
		})
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "nonexistent"}
	_, err := client.Repositories.Pipelines.Get(opts)

	assert.Error(t, err)
}

func TestPipelinesListSteps_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"uuid":       "{step-1}",
				"state":      map[string]interface{}{"name": "COMPLETED", "result": map[string]interface{}{"name": "SUCCESSFUL"}},
				"started_on": "2025-02-01T12:01:00.000000+00:00",
				"script_commands": []interface{}{
					map[string]interface{}{"command": "npm install", "name": "Install"},
				},
			},
		}))
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "{pipe-uuid}"}
	result, err := client.Repositories.Pipelines.ListSteps(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pipelines/{pipe-uuid}/steps/", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	step := values[0].(map[string]interface{})
	assert.Equal(t, "{step-1}", step["uuid"])
	state := step["state"].(map[string]interface{})
	assert.Equal(t, "COMPLETED", state["name"])
}

func TestPipelinesListSteps_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "Pipeline not found"},
		})
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "nonexistent"}
	_, err := client.Repositories.Pipelines.ListSteps(opts)

	assert.Error(t, err)
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
	result, err := client.Repositories.Pipelines.ListSteps(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "sort=started_on")
	assert.Contains(t, receivedQuery, "page=3")
	assert.Contains(t, receivedQuery, "q=state.name")
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 0)
}

func TestPipelinesGetStep_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"uuid":         "{step-uuid}",
			"state":        map[string]interface{}{"name": "COMPLETED", "result": map[string]interface{}{"name": "SUCCESSFUL"}},
			"started_on":   "2025-02-01T12:01:00.000000+00:00",
			"completed_on": "2025-02-01T12:03:00.000000+00:00",
			"setup_commands": []interface{}{
				map[string]interface{}{"command": "docker pull", "name": "Setup"},
			},
		})
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "{pipe-uuid}", StepUuid: "{step-uuid}"}
	result, err := client.Repositories.Pipelines.GetStep(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pipelines/{pipe-uuid}/steps/{step-uuid}", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, "{step-uuid}", resultMap["uuid"])
	state := resultMap["state"].(map[string]interface{})
	assert.Equal(t, "COMPLETED", state["name"])
	assert.Equal(t, "2025-02-01T12:01:00.000000+00:00", resultMap["started_on"])
	assert.Equal(t, "2025-02-01T12:03:00.000000+00:00", resultMap["completed_on"])
}

func TestPipelinesGetStep_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "Step not found"},
		})
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "{pipe-uuid}", StepUuid: "nonexistent"}
	_, err := client.Repositories.Pipelines.GetStep(opts)

	assert.Error(t, err)
}

func TestPipelinesGetLog_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("build log line 1\nbuild log line 2\n"))
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "{pipe-uuid}", StepUuid: "{step-uuid}"}
	logContent, err := client.Repositories.Pipelines.GetLog(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pipelines/{pipe-uuid}/steps/{step-uuid}/log", receivedPath)
	assert.Contains(t, logContent, "build log line 1")
	assert.Contains(t, logContent, "build log line 2")
}

func TestPipelinesGetLog_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	})
	defer server.Close()

	opts := &PipelinesOptions{Owner: "owner", RepoSlug: "repo", IDOrUuid: "bad", StepUuid: "bad"}
	_, err := client.Repositories.Pipelines.GetLog(opts)

	assert.Error(t, err)
}
