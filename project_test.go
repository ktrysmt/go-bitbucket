package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProject_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"key":         "PROJ",
			"name":        "My Project",
			"description": "A project",
			"is_private":  false,
			"uuid":        "{proj-uuid}",
		})
	})
	defer server.Close()

	opts := &ProjectOptions{Owner: "owner", Key: "PROJ"}
	project, err := client.Workspaces.GetProject(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/workspaces/owner/projects/PROJ", receivedPath)
	assert.Equal(t, "PROJ", project.Key)
	assert.Equal(t, "My Project", project.Name)
}

func TestGetProject_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "project not found"},
		})
	})
	defer server.Close()

	opts := &ProjectOptions{Owner: "owner", Key: "BAD"}
	_, err := client.Workspaces.GetProject(opts)

	assert.Error(t, err)
}

func TestCreateProject_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"key":         "NEW",
			"name":        "New Project",
			"description": "desc",
		})
	})
	defer server.Close()

	opts := &ProjectOptions{
		Owner:       "owner",
		Key:         "NEW",
		Name:        "New Project",
		Description: "desc",
		IsPrivate:   true,
	}
	project, err := client.Workspaces.CreateProject(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "NEW", project.Key)
	assert.Equal(t, "NEW", receivedBody["key"])
	assert.Equal(t, "New Project", receivedBody["name"])
	assert.Equal(t, true, receivedBody["is_private"])
}

func TestDeleteProject_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &ProjectOptions{Owner: "owner", Key: "PROJ"}
	_, err := client.Workspaces.DeleteProject(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/workspaces/owner/projects/PROJ", receivedPath)
}

func TestUpdateProject_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"key":         "PROJ",
			"name":        "Updated",
			"description": "updated desc",
		})
	})
	defer server.Close()

	opts := &ProjectOptions{
		Owner:       "owner",
		Key:         "PROJ",
		Name:        "Updated",
		Description: "updated desc",
	}
	project, err := client.Workspaces.UpdateProject(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "Updated", project.Name)
	assert.Equal(t, "Updated", receivedBody["name"])
}

func TestDecodeProject_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"key":         "PROJ",
		"name":        "Test Project",
		"description": "A test project",
		"is_private":  true,
		"uuid":        "{proj-uuid}",
	}

	project, err := decodeProject(response)

	require.NoError(t, err)
	assert.Equal(t, "PROJ", project.Key)
	assert.Equal(t, "Test Project", project.Name)
	assert.Equal(t, true, project.Is_private)
}

func TestDecodeProject_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type": "error",
		"error": map[string]interface{}{
			"message": "project not found",
		},
	}

	_, err := decodeProject(response)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project not found")
}

func TestBuildProjectBody(t *testing.T) {
	t.Parallel()
	ws := &Workspace{}
	opts := &ProjectOptions{
		Key:         "PROJ",
		Name:        "My Project",
		Description: "desc",
		IsPrivate:   true,
	}

	data, err := ws.buildProjectBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "PROJ", body["key"])
	assert.Equal(t, "My Project", body["name"])
	assert.Equal(t, "desc", body["description"])
	assert.Equal(t, true, body["is_private"])
}

func TestBuildProjectBody_MinimalFields(t *testing.T) {
	t.Parallel()
	ws := &Workspace{}
	opts := &ProjectOptions{}

	data, err := ws.buildProjectBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Nil(t, body["key"], "key should be absent when empty")
	assert.Nil(t, body["name"], "name should be absent when empty")
	assert.Nil(t, body["description"], "description should be absent when empty")
	assert.Equal(t, false, body["is_private"], "is_private should default to false")
	assert.Len(t, body, 1, "only is_private should be present")
}

func TestDecodeProject_AllFields(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"uuid":        "{full-proj-uuid}",
		"key":         "ALLF",
		"name":        "All Fields Project",
		"description": "A fully populated project",
		"is_private":  true,
	}

	project, err := decodeProject(response)

	require.NoError(t, err)
	assert.Equal(t, "{full-proj-uuid}", project.Uuid)
	assert.Equal(t, "ALLF", project.Key)
	assert.Equal(t, "All Fields Project", project.Name)
	assert.Equal(t, "A fully populated project", project.Description)
	assert.Equal(t, true, project.Is_private)
}
