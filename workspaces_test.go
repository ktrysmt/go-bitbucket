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
				map[string]interface{}{
					"workspace": map[string]interface{}{
						"slug":       "workspace1",
						"name":       "Workspace One",
						"uuid":       "{ws-1}",
						"type":       "workspace",
						"is_private": false,
					},
				},
				map[string]interface{}{
					"workspace": map[string]interface{}{
						"slug":       "workspace2",
						"name":       "Workspace Two",
						"uuid":       "{ws-2}",
						"type":       "workspace",
						"is_private": true,
					},
				},
			},
		})
	})
	defer server.Close()

	result, err := client.Workspaces.List()

	require.NoError(t, err)
	assert.Equal(t, "/2.0/user/workspaces", receivedPath)
	assert.Len(t, result.Workspaces, 2)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Pagelen)
	assert.Equal(t, 2, result.Size)

	ws0 := result.Workspaces[0]
	assert.Equal(t, "workspace1", ws0.Slug)
	assert.Equal(t, "Workspace One", ws0.Name)
	assert.Equal(t, "{ws-1}", ws0.UUID)
	assert.Equal(t, "workspace", ws0.Type)
	assert.Equal(t, false, ws0.Is_Private)

	ws1 := result.Workspaces[1]
	assert.Equal(t, "workspace2", ws1.Slug)
	assert.Equal(t, "Workspace Two", ws1.Name)
	assert.Equal(t, "{ws-2}", ws1.UUID)
	assert.Equal(t, "workspace", ws1.Type)
	assert.Equal(t, true, ws1.Is_Private)
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
	assert.Equal(t, "{ws-uuid}", result.UUID)
	assert.Equal(t, "workspace", result.Type)
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

	member := result.Members[0]
	assert.Equal(t, "member1", member.Username)
	assert.Equal(t, "Member One", member.DisplayName)
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
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedQuery = r.URL.RawQuery
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
	assert.Equal(t, "/2.0/workspaces/myorg/permissions", receivedPath)
	assert.Contains(t, receivedQuery, "user.nickname")
	assert.Contains(t, receivedQuery, "testuser")
	require.NotNil(t, result)
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
				"workspace": map[string]interface{}{
					"slug": "ws1",
					"name": "Workspace 1",
					"type": "workspace",
				},
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

func TestDecodeProjects_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page":    float64(2),
		"pagelen": float64(25),
		"size":    float64(1),
		"values": []interface{}{
			map[string]interface{}{
				"uuid":        "{proj-uuid-1}",
				"key":         "PROJ",
				"name":        "My Project",
				"description": "A test project",
				"is_private":  true,
			},
		},
	}

	result, err := decodeProjects(response)

	require.NoError(t, err)
	assert.Equal(t, int32(2), result.Page)
	assert.Equal(t, int32(25), result.Pagelen)
	assert.Equal(t, int32(1), result.Size)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, "PROJ", result.Items[0].Key)
	assert.Equal(t, "My Project", result.Items[0].Name)
	assert.Equal(t, "A test project", result.Items[0].Description)
	assert.Equal(t, true, result.Items[0].Is_private)
	assert.Equal(t, "{proj-uuid-1}", result.Items[0].Uuid)
}

// TestDecodeProjects_MaxDepthBug documents a known bug in decodeProjects:
// line 214 of workspaces.go reads "max_width" instead of "max_depth",
// so MaxDepth is always 0 even when "max_depth" is provided in the response.
func TestDecodeProjects_MaxDepthBug(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page":      float64(1),
		"pagelen":   float64(10),
		"size":      float64(0),
		"max_depth": float64(5),
		"values":    []interface{}{},
	}

	result, err := decodeProjects(response)

	require.NoError(t, err)
	// BUG: MaxDepth should be 5 but the code reads "max_width" instead of "max_depth"
	assert.Equal(t, int32(0), result.MaxDepth, "MaxDepth is always 0 because code reads 'max_width' instead of 'max_depth'")
}

func TestDecodeProjects_InvalidFormat(t *testing.T) {
	t.Parallel()
	_, err := decodeProjects("not a map")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Not a valid format")
}

// TestDecodeMembers_PaginationFieldsAreZero documents that decodeMembers
// expects page/pagelen/size as int, but JSON unmarshal produces float64,
// so pagination fields are always 0 in practice.
func TestDecodeMembers_PaginationFieldsAreZero(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page":    float64(3),
		"pagelen": float64(25),
		"size":    float64(1),
		"values": []interface{}{
			map[string]interface{}{
				"user": map[string]interface{}{
					"type":         "user",
					"username":     "testuser",
					"display_name": "Test User",
					"uuid":         "{user-uuid}",
				},
			},
		},
	}

	result, err := decodeMembers(response)

	require.NoError(t, err)
	assert.Len(t, result.Members, 1)
	assert.Equal(t, "testuser", result.Members[0].Username)
	assert.Equal(t, "Test User", result.Members[0].DisplayName)
	assert.Equal(t, "{user-uuid}", result.Members[0].Uuid)
	// BUG: decodeMembers uses type assertion .(int) but JSON produces float64,
	// so the assertion always fails and pagination fields default to 0.
	assert.Equal(t, 0, result.Page, "Page is 0 because .(int) type assertion fails on float64")
	assert.Equal(t, 0, result.Pagelen, "Pagelen is 0 because .(int) type assertion fails on float64")
	assert.Equal(t, 0, result.Size, "Size is 0 because .(int) type assertion fails on float64")
}
