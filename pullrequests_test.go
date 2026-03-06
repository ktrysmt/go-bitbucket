package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullRequestsCreate_Success(t *testing.T) {
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
			"id": 1, "title": "My PR",
		})
	})
	defer server.Close()

	opts := &PullRequestsOptions{
		Owner:             "owner",
		RepoSlug:          "repo",
		Title:             "My PR",
		SourceBranch:      "feature",
		DestinationBranch: "main",
	}
	_, err := client.Repositories.PullRequests.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/", receivedPath)
	assert.Equal(t, "My PR", receivedBody["title"])
}

func TestPullRequestsUpdate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": 1, "title": "Updated PR",
		})
	})
	defer server.Close()

	opts := &PullRequestsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		ID:       "1",
		Title:    "Updated PR",
	}
	_, err := client.Repositories.PullRequests.Update(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1", receivedPath)
}

func TestPullRequestsList_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"id": 1, "title": "PR 1"},
			map[string]interface{}{"id": 2, "title": "PR 2"},
		}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.PullRequests.List(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 2)
}

func TestPullRequestsList_WithStates(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		States:   []string{"OPEN"},
	}
	_, err := client.Repositories.PullRequests.List(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "state=OPEN")
}

func TestPullRequestsList_WithQueryAndSort(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Query:    "state=\"OPEN\"",
		Sort:     "-created_on",
	}
	_, err := client.Repositories.PullRequests.List(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "sort=-created_on")
}

func TestPullRequestsGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": 42, "title": "Test PR", "state": "OPEN",
		})
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "42"}
	result, err := client.Repositories.PullRequests.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/42", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, "Test PR", resultMap["title"])
}

func TestPullRequestsGetByCommit_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"id": 1},
		}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", Commit: "abc123"}
	_, err := client.Repositories.PullRequests.GetByCommit(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/pullrequests/", receivedPath)
}

func TestPullRequestsGetCommits_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"hash": "abc123"},
		}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.GetCommits(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/commits/", receivedPath)
}

func TestPullRequestsActivities_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.PullRequests.Activities(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/activity", receivedPath)
}

func TestPullRequestsActivity_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{})
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.Activity(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/activity", receivedPath)
}

func TestPullRequestsCommits_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.Commits(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/commits", receivedPath)
}

func TestPullRequestsMerge_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{"state": "MERGED"})
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.Merge(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
}

func TestPullRequestsDecline_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{"state": "DECLINED"})
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.Decline(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
}

func TestPullRequestsApprove_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{"approved": true})
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.Approve(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/approve", receivedPath)
}

func TestPullRequestsUnApprove_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.UnApprove(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
}

func TestPullRequestsRequestChanges_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{})
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.RequestChanges(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/request-changes", receivedPath)
}

func TestPullRequestsUnRequestChanges_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.UnRequestChanges(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
}

func TestPullRequestsAddComment_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusCreated, map[string]interface{}{"id": 1})
	})
	defer server.Close()

	opts := &PullRequestCommentOptions{
		Owner:         "owner",
		RepoSlug:      "repo",
		PullRequestID: "1",
		Content:       "LGTM",
	}
	_, err := client.Repositories.PullRequests.AddComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	content := receivedBody["content"].(map[string]interface{})
	assert.Equal(t, "LGTM", content["raw"])
}

func TestPullRequestsUpdateComment_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{"id": 5})
	})
	defer server.Close()

	opts := &PullRequestCommentOptions{
		Owner:         "owner",
		RepoSlug:      "repo",
		PullRequestID: "1",
		CommentId:     "5",
		Content:       "Updated comment",
	}
	_, err := client.Repositories.PullRequests.UpdateComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/comments/5", receivedPath)
}

func TestPullRequestsDeleteComment_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &PullRequestCommentOptions{
		Owner:         "owner",
		RepoSlug:      "repo",
		PullRequestID: "1",
		CommentId:     "5",
	}
	_, err := client.Repositories.PullRequests.DeleteComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
}

func TestPullRequestsGetComments_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"id": 1, "content": map[string]interface{}{"raw": "comment"}},
		}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.GetComments(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/comments/", receivedPath)
}

func TestPullRequestsGetComment_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{"id": 5})
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1", CommentID: "5"}
	_, err := client.Repositories.PullRequests.GetComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/comments/5", receivedPath)
}

func TestPullRequestsStatuses_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	_, err := client.Repositories.PullRequests.Statuses(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/statuses", receivedPath)
}

func TestPullRequestsGets_DelegatesToList(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.PullRequests.Gets(opts)

	require.NoError(t, err)
}

func TestBuildPullRequestBody(t *testing.T) {
	t.Parallel()
	pr := &PullRequests{}
	opts := &PullRequestsOptions{
		Title:             "My PR",
		Description:       "PR description",
		SourceBranch:      "feature",
		DestinationBranch: "main",
		Reviewers:         []string{"{uuid-1}", "{uuid-2}"},
		CloseSourceBranch: true,
		Draft:             true,
	}

	data, err := pr.buildPullRequestBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))

	assert.Equal(t, "My PR", body["title"])
	assert.Equal(t, "PR description", body["description"])
	assert.Equal(t, true, body["close_source_branch"])
	assert.Equal(t, true, body["draft"])

	reviewers := body["reviewers"].([]interface{})
	assert.Len(t, reviewers, 2)

	source := body["source"].(map[string]interface{})
	sourceBranch := source["branch"].(map[string]interface{})
	assert.Equal(t, "feature", sourceBranch["name"])

	dest := body["destination"].(map[string]interface{})
	destBranch := dest["branch"].(map[string]interface{})
	assert.Equal(t, "main", destBranch["name"])
}

func TestBuildPullRequestCommentBody(t *testing.T) {
	t.Parallel()
	pr := &PullRequests{}

	parentID := 10
	opts := &PullRequestCommentOptions{
		Content: "Nice work!",
		Parent:  &parentID,
	}

	data, err := pr.buildPullRequestCommentBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))

	content := body["content"].(map[string]interface{})
	assert.Equal(t, "Nice work!", content["raw"])

	parent := body["parent"].(map[string]interface{})
	assert.NotNil(t, parent["id"])
}

func TestPullRequestsPatch_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("diff --git a/file.txt"))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.PullRequests.Patch(opts)

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestPullRequestsDiff_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("diff output"))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.PullRequests.Diff(opts)

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestBuildPullRequestCommentBody_NoParent(t *testing.T) {
	t.Parallel()
	pr := &PullRequests{}

	opts := &PullRequestCommentOptions{
		Content: "Top-level comment",
	}

	data, err := pr.buildPullRequestCommentBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))

	assert.Nil(t, body["parent"])
}
