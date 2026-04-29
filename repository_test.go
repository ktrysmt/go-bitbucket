package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- CRUD ---

func TestRepositoryCreate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod, receivedPath string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"slug":      "my-repo",
			"full_name": "owner/my-repo",
			"name":      "my-repo",
			"scm":       "git",
		})
	})
	defer server.Close()

	opts := &RepositoryOptions{
		Owner:    "owner",
		RepoSlug: "my-repo",
		Scm:      "git",
	}
	repo, err := client.Repositories.Repository.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/my-repo", receivedPath)
	assert.Equal(t, "my-repo", repo.Slug)
	assert.Equal(t, "git", receivedBody["scm"])
}

func TestRepositoryGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"slug":      "my-repo",
			"full_name": "owner/my-repo",
			"name":      "my-repo",
		})
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "my-repo"}
	repo, err := client.Repositories.Repository.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/my-repo", receivedPath)
	assert.Equal(t, "my-repo", repo.Slug)
}

func TestRepositoryUpdate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"slug":        "my-repo",
			"description": "updated",
		})
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "my-repo", Description: "updated"}
	repo, err := client.Repositories.Repository.Update(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "updated", repo.Description)
	assert.Equal(t, "my-repo", repo.Slug)
	assert.Equal(t, "updated", receivedBody["description"])
	assert.Equal(t, "my-repo", receivedBody["name"])
}

func TestRepositoryUpdate_WithUuid(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{"slug": "my-repo"})
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "my-repo", Uuid: "{repo-uuid}"}
	repo, err := client.Repositories.Repository.Update(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/{repo-uuid}", receivedPath)
	assert.Equal(t, "my-repo", repo.Slug)
}

func TestRepositoryDelete_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod, receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "my-repo"}
	result, err := client.Repositories.Repository.Delete(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/my-repo", receivedPath)
	assert.Nil(t, result)
}

func TestRepositoryFork_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod, receivedPath string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"slug":      "forked-repo",
			"full_name": "new-owner/forked-repo",
		})
	})
	defer server.Close()

	opts := &RepositoryForkOptions{
		FromOwner: "orig-owner",
		FromSlug:  "orig-repo",
		Owner:     "new-owner",
		Name:      "forked-repo",
	}
	repo, err := client.Repositories.Repository.Fork(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/orig-owner/orig-repo/forks", receivedPath)
	assert.Equal(t, "forked-repo", repo.Slug)
	assert.Equal(t, "new-owner/forked-repo", repo.Full_name)
	assert.Equal(t, "forked-repo", receivedBody["name"])
	workspace := receivedBody["workspace"].(map[string]interface{})
	assert.Equal(t, "new-owner", workspace["slug"])
}

func TestRepositoryListWatchers_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"nickname": "watcher1"},
		}))
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.Repository.ListWatchers(opts)

	require.NoError(t, err)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	watcher := values[0].(map[string]interface{})
	assert.Equal(t, "watcher1", watcher["nickname"])
}

func TestRepositoryListForks_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"slug": "fork1"},
		}))
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.Repository.ListForks(opts)

	require.NoError(t, err)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	fork := values[0].(map[string]interface{})
	assert.Equal(t, "fork1", fork["slug"])
}

func TestRepositoryGet_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "bad"}
	repo, err := client.Repositories.Repository.Get(opts)

	require.Error(t, err)
	assert.Nil(t, repo)
	var unexpectedErr *UnexpectedResponseStatusError
	require.ErrorAs(t, err, &unexpectedErr)
	assert.Equal(t, "404 Not Found", unexpectedErr.Status)
	assert.Contains(t, string(unexpectedErr.Body), "not found")
}

func TestRepositoryDelete_WithUuid(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo", Uuid: "{repo-uuid}"}
	result, err := client.Repositories.Repository.Delete(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/{repo-uuid}", receivedPath)
	assert.Nil(t, result)
}

// --- File Operations ---

func TestBuildContentsURL(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	opts := &RepositoryFilesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Ref:      "main",
		Path:     "src/file.go",
		MaxDepth: 5,
	}
	result, err := client.Repositories.Repository.buildContentsURL(opts)

	require.NoError(t, err)
	assert.Contains(t, result, "/2.0/repositories/owner/repo/src/main/src/file.go")
	assert.Contains(t, result, "max_depth=5")
}

func TestGetFileContent_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("file content here"))
	})
	defer server.Close()

	opts := &RepositoryFilesOptions{Owner: "owner", RepoSlug: "repo", Ref: "main", Path: "README.md"}
	content, err := client.Repositories.Repository.GetFileContent(opts)

	require.NoError(t, err)
	assert.Equal(t, "file content here", string(content))
}

func TestListFiles_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"path": "file1.go", "type": "commit_file", "size": 100},
			map[string]interface{}{"path": "file2.go", "type": "commit_file", "size": 200},
		}))
	})
	defer server.Close()

	opts := &RepositoryFilesOptions{Owner: "owner", RepoSlug: "repo", Ref: "main", Path: ""}
	files, err := client.Repositories.Repository.ListFiles(opts)

	require.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Equal(t, "file1.go", files[0].Path)
	assert.Equal(t, "commit_file", files[0].Type)
	assert.Equal(t, 100, files[0].Size)
	assert.Equal(t, "file2.go", files[1].Path)
	assert.Equal(t, 200, files[1].Size)
}

func TestGetFileBlob_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("blob content"))
	})
	defer server.Close()

	opts := &RepositoryBlobOptions{Owner: "owner", RepoSlug: "repo", Ref: "main", Path: "file.txt"}
	blob, err := client.Repositories.Repository.GetFileBlob(opts)

	require.NoError(t, err)
	assert.Equal(t, "blob content", string(blob.Content))
}

// --- Branches ---

func TestListBranches_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"page":    1,
			"pagelen": 10,
			"size":    2,
			"values": []interface{}{
				map[string]interface{}{"name": "main", "type": "branch"},
				map[string]interface{}{"name": "develop", "type": "branch"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	opts := &RepositoryBranchOptions{Owner: "owner", RepoSlug: "repo", Query: "main"}
	branches, err := client.Repositories.Repository.ListBranches(opts)

	require.NoError(t, err)
	assert.Len(t, branches.Branches, 2)
	assert.Equal(t, "main", branches.Branches[0].Name)
}

func TestListBranches_WithParams(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		resp := map[string]interface{}{
			"values": []interface{}{},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	opts := &RepositoryBranchOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Sort:     "name",
		PageNum:  2,
		Pagelen:  25,
		MaxDepth: 3,
	}
	_, err := client.Repositories.Repository.ListBranches(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "sort=name")
	assert.Contains(t, receivedQuery, "page=2")
	assert.Contains(t, receivedQuery, "pagelen=25")
	assert.Contains(t, receivedQuery, "max_depth=3")
}

func TestGetBranch_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		resp := map[string]interface{}{
			"name": "main",
			"type": "branch",
			"target": map[string]interface{}{
				"hash": "abc123",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	opts := &RepositoryBranchOptions{Owner: "owner", RepoSlug: "repo", BranchName: "main"}
	branch, err := client.Repositories.Repository.GetBranch(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/refs/branches/main", receivedPath)
	assert.Equal(t, "main", branch.Name)
}

func TestGetBranch_EmptyName(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	opts := &RepositoryBranchOptions{Owner: "owner", RepoSlug: "repo", BranchName: ""}
	_, err := client.Repositories.Repository.GetBranch(opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Branch Name is empty")
}

func TestDeleteBranch_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod, receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryBranchDeleteOptions{Owner: "owner", RepoSlug: "repo", RefName: "feature"}
	err := client.Repositories.Repository.DeleteBranch(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/refs/branches/feature", receivedPath)
}

func TestDeleteBranch_WithUUIDs(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryBranchDeleteOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		RepoUUID: "{repo-uuid}",
		RefName:  "feature",
		RefUUID:  "{ref-uuid}",
	}
	err := client.Repositories.Repository.DeleteBranch(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/{repo-uuid}/refs/branches/{ref-uuid}", receivedPath)
}

func TestCreateBranch_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		resp := map[string]interface{}{
			"name": "new-branch",
			"type": "branch",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	opts := &RepositoryBranchCreationOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Name:     "new-branch",
		Target:   RepositoryBranchTarget{Hash: "abc123"},
	}
	branch, err := client.Repositories.Repository.CreateBranch(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "new-branch", branch.Name)
	assert.Equal(t, "new-branch", receivedBody["name"])
}

// --- Tags ---

func TestListTags_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"page":    1,
			"pagelen": 10,
			"size":    1,
			"values": []interface{}{
				map[string]interface{}{"name": "v1.0", "type": "tag"},
			},
		})
	})
	defer server.Close()

	opts := &RepositoryTagOptions{Owner: "owner", RepoSlug: "repo"}
	tags, err := client.Repositories.Repository.ListTags(opts)

	require.NoError(t, err)
	assert.Len(t, tags.Tags, 1)
	assert.Equal(t, "v1.0", tags.Tags[0].Name)
}

func TestCreateTag_Success(t *testing.T) {
	t.Parallel()
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		resp := map[string]interface{}{
			"name": "v2.0",
			"type": "tag",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	opts := &RepositoryTagCreationOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Name:     "v2.0",
		Target:   RepositoryTagTarget{Hash: "def456"},
	}
	tag, err := client.Repositories.Repository.CreateTag(opts)

	require.NoError(t, err)
	assert.Equal(t, "v2.0", tag.Name)
	assert.Equal(t, "v2.0", receivedBody["name"])
}

// --- Refs ---

func TestListRefs_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"page":    1,
			"pagelen": 10,
			"size":    2,
			"values": []interface{}{
				map[string]interface{}{"name": "main", "type": "branch"},
				map[string]interface{}{"name": "v1.0", "type": "tag"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	opts := &RepositoryRefOptions{Owner: "owner", RepoSlug: "repo"}
	refs, err := client.Repositories.Repository.ListRefs(opts)

	require.NoError(t, err)
	assert.Len(t, refs.Refs, 2)
	assert.Equal(t, 1, refs.Page)
	assert.Equal(t, 10, refs.Pagelen)
	assert.Equal(t, 2, refs.Size)
	assert.Equal(t, "main", refs.Refs[0]["name"])
	assert.Equal(t, "branch", refs.Refs[0]["type"])
	assert.Equal(t, "v1.0", refs.Refs[1]["name"])
	assert.Equal(t, "tag", refs.Refs[1]["type"])
}

func TestListRefs_WithParams(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		resp := map[string]interface{}{"values": []interface{}{}}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	opts := &RepositoryRefOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Query:    "name~main",
		Sort:     "name",
		PageNum:  2,
		Pagelen:  25,
		MaxDepth: 3,
	}
	_, err := client.Repositories.Repository.ListRefs(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "sort=name")
	assert.Contains(t, receivedQuery, "page=2")
	assert.Contains(t, receivedQuery, "pagelen=25")
	assert.Contains(t, receivedQuery, "max_depth=3")
}

// --- Default Reviewers ---

func TestListDefaultReviewers_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"nickname": "reviewer1", "type": "user"},
		}))
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo"}
	reviewers, err := client.Repositories.Repository.ListDefaultReviewers(opts)

	require.NoError(t, err)
	assert.Len(t, reviewers.DefaultReviewers, 1)
	assert.Equal(t, "reviewer1", reviewers.DefaultReviewers[0].Nickname)
}

func TestGetDefaultReviewer_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"nickname": "reviewer1",
			"type":     "user",
			"uuid":     "{user-uuid}",
		})
	})
	defer server.Close()

	opts := &RepositoryDefaultReviewerOptions{Owner: "owner", RepoSlug: "repo", Username: "reviewer1"}
	reviewer, err := client.Repositories.Repository.GetDefaultReviewer(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/default-reviewers/reviewer1", receivedPath)
	assert.Equal(t, "reviewer1", reviewer.Nickname)
}

func TestAddDefaultReviewer_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"nickname": "reviewer1",
			"type":     "user",
		})
	})
	defer server.Close()

	opts := &RepositoryDefaultReviewerOptions{Owner: "owner", RepoSlug: "repo", Username: "reviewer1"}
	reviewer, err := client.Repositories.Repository.AddDefaultReviewer(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "reviewer1", reviewer.Nickname)
}

func TestDeleteDefaultReviewer_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod, receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryDefaultReviewerOptions{Owner: "owner", RepoSlug: "repo", Username: "reviewer1"}
	result, err := client.Repositories.Repository.DeleteDefaultReviewer(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/default-reviewers/reviewer1", receivedPath)
	assert.Nil(t, result)
}

func TestListEffectiveDefaultReviewers_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"type":          "effective_default_reviewer",
				"reviewer_type": "required",
				"user":          map[string]interface{}{"nickname": "r1"},
			},
		}))
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo"}
	reviewers, err := client.Repositories.Repository.ListEffectiveDefaultReviewers(opts)

	require.NoError(t, err)
	assert.Len(t, reviewers.EffectiveDefaultReviewers, 1)
}

// --- Pipeline ---

func TestGetPipelineConfig_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"type":    "pipeline_config",
			"enabled": true,
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineOptions{Owner: "owner", RepoSlug: "repo"}
	pipeline, err := client.Repositories.Repository.GetPipelineConfig(opts)

	require.NoError(t, err)
	assert.True(t, pipeline.Enabled)
}

func TestUpdatePipelineConfig_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"type":    "pipeline_config",
			"enabled": true,
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineOptions{Owner: "owner", RepoSlug: "repo", Enabled: true}
	pipeline, err := client.Repositories.Repository.UpdatePipelineConfig(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.True(t, pipeline.Enabled)
}

// --- Pipeline Variables ---

func TestListPipelineVariables_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"page":    1,
			"pagelen": 10,
			"size":    1,
			"values": []interface{}{
				map[string]interface{}{"key": "VAR1", "value": "val1", "secured": false, "uuid": "{var-uuid}"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	opts := &RepositoryPipelineVariablesOptions{Owner: "owner", RepoSlug: "repo"}
	vars, err := client.Repositories.Repository.ListPipelineVariables(opts)

	require.NoError(t, err)
	assert.Len(t, vars.Variables, 1)
	assert.Equal(t, "VAR1", vars.Variables[0].Key)
}

func TestAddPipelineVariable_Success(t *testing.T) {
	t.Parallel()
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"key":     "NEW_VAR",
			"value":   "new_val",
			"secured": true,
			"uuid":    "{new-uuid}",
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineVariableOptions{
		Owner: "owner", RepoSlug: "repo",
		Key: "NEW_VAR", Value: "new_val", Secured: true,
	}
	v, err := client.Repositories.Repository.AddPipelineVariable(opts)

	require.NoError(t, err)
	assert.Equal(t, "NEW_VAR", v.Key)
	assert.Equal(t, true, receivedBody["secured"])
}

func TestGetPipelineVariable_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"key": "VAR1", "value": "val1", "uuid": "{var-uuid}",
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineVariableOptions{Owner: "owner", RepoSlug: "repo", Uuid: "{var-uuid}"}
	v, err := client.Repositories.Repository.GetPipelineVariable(opts)

	require.NoError(t, err)
	assert.Equal(t, "VAR1", v.Key)
}

func TestUpdatePipelineVariable_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"key": "VAR1", "value": "updated", "uuid": "{var-uuid}",
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineVariableOptions{
		Owner: "owner", RepoSlug: "repo", Uuid: "{var-uuid}",
		Key: "VAR1", Value: "updated",
	}
	v, err := client.Repositories.Repository.UpdatePipelineVariable(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "updated", v.Value)
}

func TestDeletePipelineVariable_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod, receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryPipelineVariableDeleteOptions{Owner: "owner", RepoSlug: "repo", Uuid: "{var-uuid}"}
	result, err := client.Repositories.Repository.DeletePipelineVariable(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Contains(t, receivedPath, "/pipelines_config/variables/{var-uuid}")
	assert.Nil(t, result)
}

// --- Pipeline KeyPair ---

func TestGetPipelineKeyPair_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"type":       "pipeline_ssh_key_pair",
			"public_key": "ssh-rsa AAAA...",
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineKeyPairOptions{Owner: "owner", RepoSlug: "repo"}
	kp, err := client.Repositories.Repository.GetPipelineKeyPair(opts)

	require.NoError(t, err)
	assert.Equal(t, "ssh-rsa AAAA...", kp.Public_key)
}

func TestAddPipelineKeyPair_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"type":       "pipeline_ssh_key_pair",
			"public_key": "ssh-rsa BBBB...",
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineKeyPairOptions{
		Owner: "owner", RepoSlug: "repo",
		PrivateKey: "private", PublicKey: "ssh-rsa BBBB...",
	}
	kp, err := client.Repositories.Repository.AddPipelineKeyPair(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "ssh-rsa BBBB...", kp.Public_key)
}

func TestDeletePipelineKeyPair_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryPipelineKeyPairOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.Repository.DeletePipelineKeyPair(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Nil(t, result)
}

// --- Pipeline Build Number ---

func TestUpdatePipelineBuildNumber_Success(t *testing.T) {
	t.Parallel()
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"type": "pipeline_build_number",
			"next": 100,
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineBuildNumberOptions{Owner: "owner", RepoSlug: "repo", Next: 100}
	bn, err := client.Repositories.Repository.UpdatePipelineBuildNumber(opts)

	require.NoError(t, err)
	assert.Equal(t, 100, bn.Next)
	assert.Equal(t, float64(100), receivedBody["next"])
}

// --- Branching Model ---

func TestBranchingModel_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"type": "branching_model",
			"branch_types": []interface{}{
				map[string]interface{}{"kind": "feature", "prefix": "feature/"},
			},
			"development": map[string]interface{}{
				"name":           "develop",
				"use_mainbranch": false,
			},
			"production": map[string]interface{}{
				"name":           "main",
				"use_mainbranch": true,
			},
		})
	})
	defer server.Close()

	opts := &RepositoryBranchingModelOptions{Owner: "owner", RepoSlug: "repo"}
	model, err := client.Repositories.Repository.BranchingModel(opts)

	require.NoError(t, err)
	assert.Len(t, model.Branch_Types, 1)
	assert.Equal(t, "feature", model.Branch_Types[0].Kind)
}

// --- Environments ---

func TestListEnvironments_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"page":    1,
			"pagelen": 10,
			"size":    1,
			"values": []interface{}{
				map[string]interface{}{
					"uuid": "{env-uuid}",
					"name": "production",
					"type": "deployment_environment",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	opts := &RepositoryEnvironmentsOptions{Owner: "owner", RepoSlug: "repo"}
	envs, err := client.Repositories.Repository.ListEnvironments(opts)

	require.NoError(t, err)
	assert.Len(t, envs.Environments, 1)
	assert.Equal(t, "production", envs.Environments[0].Name)
}

func TestAddEnvironment_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"uuid": "{new-env}",
			"name": "staging",
			"type": "deployment_environment",
		})
	})
	defer server.Close()

	opts := &RepositoryEnvironmentOptions{
		Owner: "owner", RepoSlug: "repo",
		Name:            "staging",
		EnvironmentType: Staging,
		Rank:            1,
	}
	env, err := client.Repositories.Repository.AddEnvironment(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "staging", env.Name)
	assert.Equal(t, "{new-env}", env.Uuid)
	assert.Equal(t, "deployment_environment", env.Type)
}

func TestGetEnvironment_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"uuid": "{env-uuid}",
			"name": "production",
			"type": "deployment_environment",
		})
	})
	defer server.Close()

	opts := &RepositoryEnvironmentOptions{Owner: "owner", RepoSlug: "repo", Uuid: "{env-uuid}"}
	env, err := client.Repositories.Repository.GetEnvironment(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedPath, "/environments/{env-uuid}")
	assert.Equal(t, "production", env.Name)
	assert.Equal(t, "{env-uuid}", env.Uuid)
	assert.Equal(t, "deployment_environment", env.Type)
}

func TestDeleteEnvironment_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryEnvironmentDeleteOptions{Owner: "owner", RepoSlug: "repo", Uuid: "{env-uuid}"}
	result, err := client.Repositories.Repository.DeleteEnvironment(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Nil(t, result)
}

// --- Deployment Variables ---

func TestListDeploymentVariables_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"page":    1,
			"pagelen": 10,
			"size":    1,
			"values": []interface{}{
				map[string]interface{}{"key": "DEPLOY_KEY", "value": "val", "uuid": "{dv-uuid}"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	opts := &RepositoryDeploymentVariablesOptions{
		Owner: "owner", RepoSlug: "repo",
		Environment: &Environment{Uuid: "{env-uuid}"},
	}
	vars, err := client.Repositories.Repository.ListDeploymentVariables(opts)

	require.NoError(t, err)
	assert.Len(t, vars.Variables, 1)
	assert.Equal(t, "DEPLOY_KEY", vars.Variables[0].Key)
}

func TestAddDeploymentVariable_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"key": "NEW_KEY", "value": "new_val", "uuid": "{dv-uuid}",
		})
	})
	defer server.Close()

	opts := &RepositoryDeploymentVariableOptions{
		Owner: "owner", RepoSlug: "repo",
		Environment: &Environment{Uuid: "{env-uuid}"},
		Key:         "NEW_KEY",
		Value:       "new_val",
	}
	v, err := client.Repositories.Repository.AddDeploymentVariable(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedPath, "/deployments_config/environments/{env-uuid}/variables")
	assert.Equal(t, "NEW_KEY", v.Key)
}

func TestUpdateDeploymentVariable_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"key": "KEY1", "value": "updated", "uuid": "{dv-uuid}",
		})
	})
	defer server.Close()

	opts := &RepositoryDeploymentVariableOptions{
		Owner: "owner", RepoSlug: "repo",
		Environment: &Environment{Uuid: "{env-uuid}"},
		Uuid:        "{dv-uuid}",
		Key:         "KEY1",
		Value:       "updated",
	}
	v, err := client.Repositories.Repository.UpdateDeploymentVariable(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "updated", v.Value)
}

func TestDeleteDeploymentVariable_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod, receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryDeploymentVariableDeleteOptions{
		Owner: "owner", RepoSlug: "repo",
		Environment: &Environment{Uuid: "{env-uuid}"},
		Uuid:        "{dv-uuid}",
	}
	result, err := client.Repositories.Repository.DeleteDeploymentVariable(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Contains(t, receivedPath, "/deployments_config/environments/{env-uuid}/variables/{dv-uuid}")
	assert.Nil(t, result)
}

// --- Group Permissions ---

func TestListGroupPermissions_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"type":       "repository_group_permission",
				"permission": "write",
				"group":      map[string]interface{}{"slug": "devs"},
			},
		}))
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo"}
	perms, err := client.Repositories.Repository.ListGroupPermissions(opts)

	require.NoError(t, err)
	assert.Len(t, perms.GroupPermissions, 1)
}

func TestSetGroupPermissions_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod, receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"permission": "admin",
			"group":      map[string]interface{}{"slug": "devs"},
		})
	})
	defer server.Close()

	opts := &RepositoryGroupPermissionsOptions{
		Owner: "owner", RepoSlug: "repo",
		Group: "devs", Permission: "admin",
	}
	perm, err := client.Repositories.Repository.SetGroupPermissions(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Contains(t, receivedPath, "/permissions-config/groups/devs")
	assert.Equal(t, "admin", perm.Permission)
}

func TestGetGroupPermissions_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"permission": "read",
				"group":      map[string]interface{}{"slug": "readers"},
			},
		}))
	})
	defer server.Close()

	opts := &RepositoryGroupPermissionsOptions{Owner: "owner", RepoSlug: "repo", Group: "readers"}
	// Note: GetGroupPermissions returns the paginated response decoded as a single GroupPermission
	_, err := client.Repositories.Repository.GetGroupPermissions(opts)

	require.NoError(t, err)
}

func TestDeleteGroupPermissions_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryGroupPermissionsOptions{Owner: "owner", RepoSlug: "repo", Group: "devs"}
	result, err := client.Repositories.Repository.DeleteGroupPermissions(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Nil(t, result)
}

// --- User Permissions ---

func TestListUserPermissions_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"type":       "repository_user_permission",
				"permission": "write",
				"user":       map[string]interface{}{"nickname": "jdoe"},
			},
		}))
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo"}
	perms, err := client.Repositories.Repository.ListUserPermissions(opts)

	require.NoError(t, err)
	assert.Len(t, perms.UserPermissions, 1)
}

func TestSetUserPermissions_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"permission": "admin",
			"user":       map[string]interface{}{"nickname": "jdoe"},
		})
	})
	defer server.Close()

	opts := &RepositoryUserPermissionsOptions{
		Owner: "owner", RepoSlug: "repo",
		User: "jdoe", Permission: "admin",
	}
	perm, err := client.Repositories.Repository.SetUserPermissions(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "admin", perm.Permission)
}

func TestGetUserPermissions_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"permission": "read",
				"user":       map[string]interface{}{"nickname": "jdoe"},
			},
		}))
	})
	defer server.Close()

	opts := &RepositoryUserPermissionsOptions{Owner: "owner", RepoSlug: "repo", User: "jdoe"}
	_, err := client.Repositories.Repository.GetUserPermissions(opts)

	require.NoError(t, err)
}

func TestDeleteUserPermissions_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &RepositoryUserPermissionsOptions{Owner: "owner", RepoSlug: "repo", User: "jdoe"}
	result, err := client.Repositories.Repository.DeleteUserPermissions(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Nil(t, result)
}

// --- Build Body Functions ---

func TestBuildRepositoryBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryOptions{
		Owner:       "owner",
		RepoSlug:    "my-repo",
		Scm:         "git",
		IsPrivate:   "true",
		Description: "desc",
		ForkPolicy:  "allow_forks",
		Language:    "go",
		HasIssues:   "true",
		HasWiki:     "true",
		Project:     "PROJ",
	}

	data, err := repo.buildRepositoryBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "my-repo", body["name"])
	assert.Equal(t, "git", body["scm"])
	assert.Equal(t, true, body["is_private"])
	assert.Equal(t, "allow_forks", body["fork_policy"])
	assert.Equal(t, false, body["no_forks"])
	assert.Equal(t, false, body["no_public_forks"])
}

func TestBuildRepositoryBody_ForkPolicyNoForks(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryOptions{ForkPolicy: "no_forks"}

	data, err := repo.buildRepositoryBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, true, body["no_forks"])
	assert.Equal(t, true, body["no_public_forks"])
}

func TestBuildRepositoryBody_ForkPolicyNoPublicForks(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryOptions{ForkPolicy: "no_public_forks"}

	data, err := repo.buildRepositoryBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, false, body["no_forks"])
	assert.Equal(t, true, body["no_public_forks"])
}

func TestBuildForkBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryForkOptions{
		Owner:       "new-owner",
		Name:        "forked",
		IsPrivate:   "false",
		Description: "fork desc",
		ForkPolicy:  "allow_forks",
		Project:     "PROJ",
	}

	data, err := repo.buildForkBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "forked", body["name"])
	assert.Equal(t, false, body["is_private"])
	workspace := body["workspace"].(map[string]interface{})
	assert.Equal(t, "new-owner", workspace["slug"])
}

func TestBuildPipelineBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryPipelineOptions{Enabled: true}

	data, err := repo.buildPipelineBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, true, body["enabled"])
}

func TestBuildPipelineVariableBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryPipelineVariableOptions{
		Uuid: "{uuid}", Key: "KEY", Value: "VAL", Secured: true,
	}

	data, err := repo.buildPipelineVariableBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "{uuid}", body["uuid"])
	assert.Equal(t, "KEY", body["key"])
	assert.Equal(t, "VAL", body["value"])
	assert.Equal(t, true, body["secured"])
}

func TestBuildPipelineKeyPairBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryPipelineKeyPairOptions{PrivateKey: "priv", PublicKey: "pub"}

	data, err := repo.buildPipelineKeyPairBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "priv", body["private_key"])
	assert.Equal(t, "pub", body["public_key"])
}

func TestBuildPipelineBuildNumberBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryPipelineBuildNumberOptions{Next: 42}

	data, err := repo.buildPipelineBuildNumberBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, float64(42), body["next"])
}

func TestBuildBranchBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryBranchCreationOptions{
		Name:   "feature",
		Target: RepositoryBranchTarget{Hash: "abc123"},
	}

	data, err := repo.buildBranchBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "feature", body["name"])
	target := body["target"].(map[string]interface{})
	assert.Equal(t, "abc123", target["hash"])
}

func TestBuildTagBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryTagCreationOptions{
		Name:   "v1.0",
		Target: RepositoryTagTarget{Hash: "def456"},
	}

	data, err := repo.buildTagBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "v1.0", body["name"])
	target := body["target"].(map[string]interface{})
	assert.Equal(t, "def456", target["hash"])
}

func TestBuildEnvironmentBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryEnvironmentOptions{
		Name:            "staging",
		EnvironmentType: Staging,
		Rank:            1,
		Uuid:            "{env-uuid}",
	}

	data, err := repo.buildEnvironmentBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "staging", body["name"])
	assert.Equal(t, "{env-uuid}", body["uuid"])
	envType := body["environment_type"].(map[string]interface{})
	assert.Equal(t, "Staging", envType["name"])
}

func TestBuildDeploymentVariableBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryDeploymentVariableOptions{
		Uuid: "{dv-uuid}", Key: "KEY", Value: "VAL", Secured: true,
	}

	data, err := repo.buildDeploymentVariableBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "{dv-uuid}", body["uuid"])
	assert.Equal(t, "KEY", body["key"])
	assert.Equal(t, true, body["secured"])
}

func TestBuildRepositoryGroupPermissionBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryGroupPermissionsOptions{Permission: "write"}

	data, err := repo.buildRepositoryGroupPermissionBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "write", body["permission"])
}

func TestBuildRepositoryUserPermissionBody(t *testing.T) {
	t.Parallel()
	repo := &Repository{}
	opts := &RepositoryUserPermissionsOptions{Permission: "admin"}

	data, err := repo.buildRepositoryUserPermissionBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))
	assert.Equal(t, "admin", body["permission"])
}

// --- Regression: returned *Repository must carry a usable client ---
// https://github.com/ktrysmt/go-bitbucket/issues/347

func TestRepositoryGet_ReturnedRepoHasClient(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/2.0/repositories/owner/my-repo":
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"slug": "my-repo", "full_name": "owner/my-repo", "name": "my-repo",
			})
		case "/2.0/repositories/owner/my-repo/src/main/README.md":
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("hello"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	repo, err := client.Repositories.Repository.Get(&RepositoryOptions{Owner: "owner", RepoSlug: "my-repo"})
	require.NoError(t, err)
	require.NotNil(t, repo)

	content, err := repo.GetFileContent(&RepositoryFilesOptions{
		Owner: "owner", RepoSlug: "my-repo", Ref: "main", Path: "README.md",
	})
	require.NoError(t, err)
	assert.Equal(t, "hello", string(content))
}

func TestRepositoryCreate_ReturnedRepoHasClient(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{"slug": "my-repo"})
	})
	defer server.Close()

	repo, err := client.Repositories.Repository.Create(&RepositoryOptions{Owner: "owner", RepoSlug: "my-repo"})
	require.NoError(t, err)
	require.NotNil(t, repo)
	assert.Same(t, client, repo.c)
}

func TestRepositoryFork_ReturnedRepoHasClient(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{"slug": "forked"})
	})
	defer server.Close()

	repo, err := client.Repositories.Repository.Fork(&RepositoryForkOptions{
		FromOwner: "orig", FromSlug: "orig-repo", Owner: "new", Name: "forked",
	})
	require.NoError(t, err)
	require.NotNil(t, repo)
	assert.Same(t, client, repo.c)
}

func TestRepositoryUpdate_ReturnedRepoHasClient(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{"slug": "my-repo"})
	})
	defer server.Close()

	repo, err := client.Repositories.Repository.Update(&RepositoryOptions{Owner: "owner", RepoSlug: "my-repo"})
	require.NoError(t, err)
	require.NotNil(t, repo)
	assert.Same(t, client, repo.c)
}

func TestRepositoriesList_ItemsHaveClient(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"slug": "repo1", "full_name": "owner/repo1"},
			map[string]interface{}{"slug": "repo2", "full_name": "owner/repo2"},
		}))
	})
	defer server.Close()

	res, err := client.Repositories.ListForAccount(&RepositoriesOptions{Owner: "owner"})
	require.NoError(t, err)
	require.Len(t, res.Items, 2)
	for i := range res.Items {
		assert.Same(t, client, res.Items[i].c, "item %d should have client wired", i)
	}
}

// --- Decode Functions ---

func TestDecodeRepository_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"slug":        "my-repo",
		"full_name":   "owner/my-repo",
		"name":        "my-repo",
		"description": "a repo",
		"is_private":  true,
		"scm":         "git",
		"language":    "go",
	}

	repo, err := decodeRepository(response)

	require.NoError(t, err)
	assert.Equal(t, "my-repo", repo.Slug)
	assert.Equal(t, true, repo.Is_private)
	assert.Equal(t, "go", repo.Language)
}

func TestDecodeRepository_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type": "error",
		"error": map[string]interface{}{
			"message": "repo not found",
		},
	}

	_, err := decodeRepository(response)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repo not found")
}

func TestDecodeRepositoryFiles_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"values": []interface{}{
			map[string]interface{}{"path": "file1.go", "type": "commit_file", "size": 100},
			map[string]interface{}{"path": "file2.go", "type": "commit_file", "size": 200},
		},
	}

	files, err := decodeRepositoryFiles(response)

	require.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Equal(t, "file1.go", files[0].Path)
}

func TestDecodeRepositoryRefs_Success(t *testing.T) {
	t.Parallel()
	jsonStr := `{
		"page": 1, "pagelen": 10, "size": 2, "next": "http://next",
		"values": [
			{"name": "main", "type": "branch"},
			{"name": "v1.0", "type": "tag"}
		]
	}`

	refs, err := decodeRepositoryRefs(jsonStr)

	require.NoError(t, err)
	assert.Equal(t, 1, refs.Page)
	assert.Equal(t, 10, refs.Pagelen)
	assert.Equal(t, 2, refs.Size)
	assert.Equal(t, "http://next", refs.Next)
	assert.Len(t, refs.Refs, 2)
}

func TestDecodeRepositoryBranches_Success(t *testing.T) {
	t.Parallel()
	jsonStr := `{
		"page": 1, "pagelen": 10, "size": 1,
		"values": [
			{"name": "main", "type": "branch"}
		]
	}`

	branches, err := decodeRepositoryBranches(jsonStr)

	require.NoError(t, err)
	assert.Len(t, branches.Branches, 1)
	assert.Equal(t, "main", branches.Branches[0].Name)
}

func TestDecodeRepositoryBranch_Success(t *testing.T) {
	t.Parallel()
	jsonStr := `{"name": "develop", "type": "branch", "target": {"hash": "abc"}}`

	branch, err := decodeRepositoryBranch(jsonStr)

	require.NoError(t, err)
	assert.Equal(t, "develop", branch.Name)
}

func TestDecodeRepositoryBranchCreated_Success(t *testing.T) {
	t.Parallel()
	jsonStr := `{"name": "new-branch", "type": "branch"}`

	branch, err := decodeRepositoryBranchCreated(jsonStr)

	require.NoError(t, err)
	assert.Equal(t, "new-branch", branch.Name)
}

func TestDecodeRepositoryTagCreated_Success(t *testing.T) {
	t.Parallel()
	jsonStr := `{"name": "v2.0", "type": "tag"}`

	tag, err := decodeRepositoryTagCreated(jsonStr)

	require.NoError(t, err)
	assert.Equal(t, "v2.0", tag.Name)
}

func TestDecodeRepositoryTags_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page": float64(1), "pagelen": float64(10), "size": float64(1),
		"values": []interface{}{
			map[string]interface{}{"name": "v1.0", "type": "tag"},
		},
	}

	tags, err := decodeRepositoryTags(response)

	require.NoError(t, err)
	assert.Len(t, tags.Tags, 1)
	assert.Equal(t, "v1.0", tags.Tags[0].Name)
	assert.Equal(t, 1, tags.Page)
}

func TestDecodePipelineRepository_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type":    "pipeline_config",
		"enabled": true,
	}

	pipeline, err := decodePipelineRepository(response)

	require.NoError(t, err)
	assert.True(t, pipeline.Enabled)
}

func TestDecodePipelineRepository_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type":  "error",
		"error": map[string]interface{}{"message": "not found"},
	}

	_, err := decodePipelineRepository(response)

	assert.Error(t, err)
}

func TestDecodePipelineVariables_Success(t *testing.T) {
	t.Parallel()
	jsonStr := `{
		"page": 1, "pagelen": 10, "size": 1,
		"values": [
			{"key": "VAR1", "value": "val1", "secured": false, "uuid": "{uuid}"}
		]
	}`

	vars, err := decodePipelineVariables(jsonStr)

	require.NoError(t, err)
	assert.Len(t, vars.Variables, 1)
	assert.Equal(t, "VAR1", vars.Variables[0].Key)
}

func TestDecodePipelineVariableRepository_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"key": "VAR1", "value": "val1", "uuid": "{uuid}",
	}

	v, err := decodePipelineVariableRepository(response)

	require.NoError(t, err)
	assert.Equal(t, "VAR1", v.Key)
}

func TestDecodePipelineKeyPairRepository_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type":       "pipeline_ssh_key_pair",
		"public_key": "ssh-rsa AAAA...",
	}

	kp, err := decodePipelineKeyPairRepository(response)

	require.NoError(t, err)
	assert.Equal(t, "ssh-rsa AAAA...", kp.Public_key)
}

func TestDecodePipelineBuildNumberRepository_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type": "pipeline_build_number",
		"next": 42,
	}

	bn, err := decodePipelineBuildNumberRepository(response)

	require.NoError(t, err)
	assert.Equal(t, 42, bn.Next)
}

func TestDecodeBranchingModel_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type": "branching_model",
		"branch_types": []interface{}{
			map[string]interface{}{"kind": "feature", "prefix": "feature/"},
		},
	}

	model, err := decodeBranchingModel(response)

	require.NoError(t, err)
	assert.Len(t, model.Branch_Types, 1)
}

func TestDecodeBranchingModel_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type":  "error",
		"error": map[string]interface{}{"message": "not found"},
	}

	_, err := decodeBranchingModel(response)

	assert.Error(t, err)
}

func TestDecodeEnvironments_Success(t *testing.T) {
	t.Parallel()
	jsonStr := `{
		"page": 1, "pagelen": 10, "size": 1,
		"values": [
			{"uuid": "{env-uuid}", "name": "production", "type": "deployment_environment"}
		]
	}`

	envs, err := decodeEnvironments(jsonStr)

	require.NoError(t, err)
	assert.Len(t, envs.Environments, 1)
	assert.Equal(t, "production", envs.Environments[0].Name)
}

func TestDecodeEnvironment_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"uuid": "{env-uuid}",
		"name": "staging",
		"type": "deployment_environment",
	}

	env, err := decodeEnvironment(response)

	require.NoError(t, err)
	assert.Equal(t, "staging", env.Name)
}

func TestDecodeEnvironment_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type":  "error",
		"error": map[string]interface{}{"message": "not found"},
	}

	_, err := decodeEnvironment(response)

	assert.Error(t, err)
}

func TestDecodeDeploymentVariables_Success(t *testing.T) {
	t.Parallel()
	jsonStr := `{
		"page": 1, "pagelen": 10, "size": 1,
		"values": [
			{"key": "DV1", "value": "val", "uuid": "{dv-uuid}"}
		]
	}`

	vars, err := decodeDeploymentVariables(jsonStr)

	require.NoError(t, err)
	assert.Len(t, vars.Variables, 1)
	assert.Equal(t, "DV1", vars.Variables[0].Key)
}

func TestDecodeDeploymentVariable_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"key": "DV1", "value": "val", "uuid": "{dv-uuid}",
	}

	v, err := decodeDeploymentVariable(response)

	require.NoError(t, err)
	assert.Equal(t, "DV1", v.Key)
}

func TestDecodeDeploymentVariable_ErrorType(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type":  "error",
		"error": map[string]interface{}{"message": "not found"},
	}

	_, err := decodeDeploymentVariable(response)

	assert.Error(t, err)
}

func TestDecodeDefaultReviewer_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"nickname":     "reviewer1",
		"display_name": "Reviewer One",
		"type":         "user",
		"uuid":         "{user-uuid}",
		"account_id":   "123",
	}

	reviewer, err := decodeDefaultReviewer(response)

	require.NoError(t, err)
	assert.Equal(t, "reviewer1", reviewer.Nickname)
	assert.Equal(t, "Reviewer One", reviewer.DisplayName)
	assert.Equal(t, "123", reviewer.AccountId)
}

func TestDecodeDefaultReviewers_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page": float64(1), "pagelen": float64(10), "size": float64(1),
		"values": []interface{}{
			map[string]interface{}{"nickname": "r1", "type": "user"},
		},
	}

	reviewers, err := decodeDefaultReviewers(response)

	require.NoError(t, err)
	assert.Len(t, reviewers.DefaultReviewers, 1)
}

func TestDecodeEffectiveDefaultReviewers_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page": float64(1), "pagelen": float64(10), "size": float64(1),
		"values": []interface{}{
			map[string]interface{}{
				"type":          "effective_default_reviewer",
				"reviewer_type": "required",
				"user":          map[string]interface{}{"nickname": "r1"},
			},
		},
	}

	reviewers, err := decodeEffectiveDefaultReviewers(response)

	require.NoError(t, err)
	assert.Len(t, reviewers.EffectiveDefaultReviewers, 1)
}

func TestDecodeGroupPermissions_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type":       "repository_group_permission",
		"permission": "write",
		"group":      map[string]interface{}{"slug": "devs"},
	}

	perm, err := decodeGroupPermissions(response)

	require.NoError(t, err)
	assert.Equal(t, "write", perm.Permission)
}

func TestDecodeGroupsPermissions_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page": float64(1), "pagelen": float64(10), "size": float64(1),
		"values": []interface{}{
			map[string]interface{}{
				"permission": "write",
				"group":      map[string]interface{}{"slug": "devs"},
			},
		},
	}

	perms, err := decodeGroupsPermissions(response)

	require.NoError(t, err)
	assert.Len(t, perms.GroupPermissions, 1)
}

func TestDecodeUserPermissions_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"type":       "repository_user_permission",
		"permission": "admin",
		"user":       map[string]interface{}{"nickname": "jdoe"},
	}

	perm, err := decodeUserPermissions(response)

	require.NoError(t, err)
	assert.Equal(t, "admin", perm.Permission)
}

func TestDecodeUsersPermissions_Success(t *testing.T) {
	t.Parallel()
	response := map[string]interface{}{
		"page": float64(1), "pagelen": float64(10), "size": float64(1),
		"values": []interface{}{
			map[string]interface{}{
				"permission": "write",
				"user":       map[string]interface{}{"nickname": "jdoe"},
			},
		},
	}

	perms, err := decodeUsersPermissions(response)

	require.NoError(t, err)
	assert.Len(t, perms.UserPermissions, 1)
}

// --- String methods ---

func TestRepositoryFile_String(t *testing.T) {
	t.Parallel()
	rf := RepositoryFile{Path: "src/main.go"}
	assert.Equal(t, "src/main.go", rf.String())
}

func TestRepositoryBlob_String(t *testing.T) {
	t.Parallel()
	rb := RepositoryBlob{Content: []byte("hello")}
	assert.Equal(t, "hello", rb.String())
}

// --- WriteFileBlob ---

func TestWriteFileBlob_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedContentType string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusCreated)
	})
	defer server.Close()

	tmpFile, err := os.CreateTemp("", "test-blob-*.txt")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, err = tmpFile.WriteString("file content")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	opts := &RepositoryBlobWriteOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Files:    []File{{Path: tmpFile.Name(), Name: "test.txt"}},
		Message:  "add file",
		Branch:   "main",
		Author:   "Test <test@example.com>",
	}
	err = client.Repositories.Repository.WriteFileBlob(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Contains(t, receivedContentType, "multipart/form-data")
}

func TestWriteFileBlob_BothFilesAndFilename(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	opts := &RepositoryBlobWriteOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		FileName: "file.txt",
		Files:    []File{{Path: "other.txt", Name: "other.txt"}},
	}
	err := client.Repositories.Repository.WriteFileBlob(opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can't specify both files and filename")
}

func TestWriteFileBlob_FileNotFound(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	opts := &RepositoryBlobWriteOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		FileName: "/nonexistent/path.txt",
	}
	err := client.Repositories.Repository.WriteFileBlob(opts)

	assert.Error(t, err)
}

// --- Error paths for Pipeline, Environment, etc. ---

func TestGetPipelineConfig_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineOptions{Owner: "owner", RepoSlug: "repo"}
	pipeline, err := client.Repositories.Repository.GetPipelineConfig(opts)

	require.Error(t, err)
	assert.Nil(t, pipeline)
	var unexpectedErr *UnexpectedResponseStatusError
	require.ErrorAs(t, err, &unexpectedErr)
	assert.Equal(t, "404 Not Found", unexpectedErr.Status)
}

func TestGetPipelineKeyPair_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineKeyPairOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Repository.GetPipelineKeyPair(opts)

	assert.Error(t, err)
}

func TestGetPipelineVariable_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineVariableOptions{Owner: "owner", RepoSlug: "repo", Uuid: "{bad}"}
	_, err := client.Repositories.Repository.GetPipelineVariable(opts)

	assert.Error(t, err)
}

func TestBranchingModel_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryBranchingModelOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Repository.BranchingModel(opts)

	assert.Error(t, err)
}

func TestGetEnvironment_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryEnvironmentOptions{Owner: "owner", RepoSlug: "repo", Uuid: "{bad}"}
	env, err := client.Repositories.Repository.GetEnvironment(opts)

	require.Error(t, err)
	assert.Nil(t, env)
	var unexpectedErr *UnexpectedResponseStatusError
	require.ErrorAs(t, err, &unexpectedErr)
	assert.Equal(t, "404 Not Found", unexpectedErr.Status)
}

func TestGetDefaultReviewer_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryDefaultReviewerOptions{Owner: "owner", RepoSlug: "repo", Username: "bad"}
	_, err := client.Repositories.Repository.GetDefaultReviewer(opts)

	assert.Error(t, err)
}

func TestListDefaultReviewers_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Repository.ListDefaultReviewers(opts)

	assert.Error(t, err)
}

func TestListEffectiveDefaultReviewers_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Repository.ListEffectiveDefaultReviewers(opts)

	assert.Error(t, err)
}

func TestAddDefaultReviewer_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "user not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryDefaultReviewerOptions{Owner: "owner", RepoSlug: "repo", Username: "bad"}
	_, err := client.Repositories.Repository.AddDefaultReviewer(opts)

	assert.Error(t, err)
}

func TestListGroupPermissions_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo"}
	perms, err := client.Repositories.Repository.ListGroupPermissions(opts)

	require.Error(t, err)
	assert.Nil(t, perms)
	var unexpectedErr *UnexpectedResponseStatusError
	require.ErrorAs(t, err, &unexpectedErr)
	assert.Equal(t, "403 Forbidden", unexpectedErr.Status)
}

func TestListUserPermissions_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	opts := &RepositoryOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Repository.ListUserPermissions(opts)

	assert.Error(t, err)
}

func TestGetFileBlob_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "file not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryBlobOptions{Owner: "owner", RepoSlug: "repo", Ref: "main", Path: "bad"}
	_, err := client.Repositories.Repository.GetFileBlob(opts)

	assert.Error(t, err)
}

func TestGetFileContent_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryFilesOptions{Owner: "owner", RepoSlug: "repo", Ref: "main", Path: "bad"}
	_, err := client.Repositories.Repository.GetFileContent(opts)

	assert.Error(t, err)
}

func TestListRefs_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	opts := &RepositoryRefOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Repository.ListRefs(opts)

	assert.Error(t, err)
}

func TestListBranches_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	opts := &RepositoryBranchOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Repository.ListBranches(opts)

	assert.Error(t, err)
}

func TestGetBranch_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &RepositoryBranchOptions{Owner: "owner", RepoSlug: "repo", BranchName: "bad"}
	_, err := client.Repositories.Repository.GetBranch(opts)

	assert.Error(t, err)
}

func TestCreateBranch_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{"message": "bad request"},
		})
	})
	defer server.Close()

	opts := &RepositoryBranchCreationOptions{
		Owner: "owner", RepoSlug: "repo",
		Name: "bad", Target: RepositoryBranchTarget{Hash: "invalid"},
	}
	branch, err := client.Repositories.Repository.CreateBranch(opts)

	require.Error(t, err)
	assert.Nil(t, branch)
	var unexpectedErr *UnexpectedResponseStatusError
	require.ErrorAs(t, err, &unexpectedErr)
	assert.Equal(t, "400 Bad Request", unexpectedErr.Status)
	assert.Contains(t, string(unexpectedErr.Body), "bad request")
}

func TestCreateTag_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{"message": "bad request"},
		})
	})
	defer server.Close()

	opts := &RepositoryTagCreationOptions{
		Owner: "owner", RepoSlug: "repo",
		Name: "bad", Target: RepositoryTagTarget{Hash: "invalid"},
	}
	_, err := client.Repositories.Repository.CreateTag(opts)

	assert.Error(t, err)
}

func TestListEnvironments_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	opts := &RepositoryEnvironmentsOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Repository.ListEnvironments(opts)

	assert.Error(t, err)
}

func TestListPipelineVariables_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	opts := &RepositoryPipelineVariablesOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Repository.ListPipelineVariables(opts)

	assert.Error(t, err)
}

func TestListDeploymentVariables_Error(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	opts := &RepositoryDeploymentVariablesOptions{
		Owner: "owner", RepoSlug: "repo",
		Environment: &Environment{Uuid: "{env}"},
	}
	_, err := client.Repositories.Repository.ListDeploymentVariables(opts)

	assert.Error(t, err)
}
