package bitbucket

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceList_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"page":    float64(1),
			"pagelen": float64(10),
			"size":    float64(2),
			"values": []interface{}{
				map[string]interface{}{"slug": "workspace1", "name": "Workspace One", "uuid": "{ws-1}"},
				map[string]interface{}{"slug": "workspace2", "name": "Workspace Two", "uuid": "{ws-2}"},
			},
		})
	})
	defer server.Close()

	result, err := client.Workspaces.List()

	require.NoError(t, err)
	assert.Equal(t, "/2.0/workspaces", receivedPath)
	assert.Len(t, result.Workspaces, 2)
	assert.Equal(t, "workspace1", result.Workspaces[0].Slug)
	assert.Equal(t, "workspace2", result.Workspaces[1].Slug)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 2, result.Size)
}

func TestWorkspaceGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"slug": "myworkspace",
			"name": "My Workspace",
			"uuid": "{ws-uuid}",
			"type": "workspace",
		})
	})
	defer server.Close()

	result, err := client.Workspaces.Get("myworkspace")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/workspaces/myworkspace", receivedPath)
	assert.Equal(t, "myworkspace", result.Slug)
	assert.Equal(t, "My Workspace", result.Name)
}

func TestWorkspaceGet_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"type": "error",
			"error": map[string]interface{}{
				"message": "workspace not found",
			},
		})
	})
	defer server.Close()

	_, err := client.Workspaces.Get("nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "workspace not found")
}

func TestWorkspaceMembers_Success(t *testing.T) {
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
					"user": map[string]interface{}{
						"type":         "user",
						"username":     "member1",
						"display_name": "Member One",
					},
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Workspaces.Members("myworkspace")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/workspaces/myworkspace/members", receivedPath)
	assert.Len(t, result.Members, 1)
	assert.Equal(t, "member1", result.Members[0].Username)
}

func TestWorkspaceProjects_Success(t *testing.T) {
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
					"key":         "PROJ",
					"name":        "My Project",
					"description": "A project",
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Workspaces.Projects("myworkspace")

	require.NoError(t, err)
	assert.Equal(t, "/2.0/workspaces/myworkspace/projects/", receivedPath)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "PROJ", result.Items[0].Key)
}

func TestPermissionGetUserPermissions_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"size": float64(1),
			"values": []interface{}{
				map[string]interface{}{
					"permission": "owner",
					"user":       map[string]interface{}{"nickname": "testuser"},
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Workspaces.Permissions.GetUserPermissions("myorg", "testuser")

	require.NoError(t, err)
	assert.Contains(t, receivedPath, "/2.0/workspaces/myorg/permissions")
	assert.Equal(t, "owner", result.Type)
}

func TestPermissionGetUserPermissions_NoPermission(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"size":   float64(0),
			"values": []interface{}{},
		})
	})
	defer server.Close()

	result, err := client.Workspaces.Permissions.GetUserPermissions("myorg", "unknown")

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestPermissionGetUserPermissionsByUuid_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"size": float64(1),
			"values": []interface{}{
				map[string]interface{}{
					"permission": "admin",
					"user":       map[string]interface{}{"uuid": "{user-uuid}"},
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Workspaces.Permissions.GetUserPermissionsByUuid("myorg", "{user-uuid}")

	require.NoError(t, err)
	assert.Equal(t, "admin", result.Type)
}

func TestDecodeWorkspace_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"slug": "test-ws",
		"name": "Test Workspace",
		"uuid": "{uuid}",
		"type": "workspace",
	}

	ws, err := decodeWorkspace(response)

	require.NoError(t, err)
	assert.Equal(t, "test-ws", ws.Slug)
	assert.Equal(t, "Test Workspace", ws.Name)
}

func TestDecodeWorkspace_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type": "error",
		"error": map[string]interface{}{
			"message": "not found",
		},
	}

	_, err := decodeWorkspace(response)

	assert.Error(t, err)
}

func TestDecodeWorkspaceList_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page":    float64(1),
		"pagelen": float64(10),
		"size":    float64(1),
		"values": []interface{}{
			map[string]interface{}{
				"slug": "ws1",
				"name": "Workspace 1",
				"type": "workspace",
			},
		},
	}

	result, err := decodeWorkspaceList(response)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Pagelen)
	assert.Equal(t, 1, result.Size)
	assert.Len(t, result.Workspaces, 1)
}

func TestDecodePermission_WithPermission(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"size": float64(1),
		"values": []interface{}{
			map[string]interface{}{
				"permission": "write",
			},
		},
	}

	result := decodePermission(response)

	assert.NotNil(t, result)
	assert.Equal(t, "write", result.Type)
}

func TestDecodePermission_Empty(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"size":   float64(0),
		"values": []interface{}{},
	}

	result := decodePermission(response)

	assert.Nil(t, result)
}
