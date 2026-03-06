package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssuesGets_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"id": 1, "title": "bug report"},
		}))
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
	}
	result, err := client.Repositories.Issues.Gets(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
}

func TestIssuesGets_WithStates(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		States:   []string{"open"},
	}
	_, err := client.Repositories.Issues.Gets(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "state=open")
}

func TestIssuesGets_WithQueryAndSort(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Query:    "priority = \"critical\"",
		Sort:     "-created_on",
	}
	_, err := client.Repositories.Issues.Gets(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "sort=-created_on")
}

func TestIssuesGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": 42, "title": "test issue", "state": "open",
		})
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		ID:       "42",
	}
	result, err := client.Repositories.Issues.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/42", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, "test issue", resultMap["title"])
}

func TestIssuesDelete_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		ID:       "42",
	}
	_, err := client.Repositories.Issues.Delete(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
}

func TestIssuesUpdate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": 42, "title": "updated title", "state": "closed",
		})
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		ID:       "42",
		Title:    "updated title",
		State:    "closed",
	}
	_, err := client.Repositories.Issues.Update(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "updated title", receivedBody["title"])
	assert.Equal(t, "closed", receivedBody["state"])
}

func TestIssuesCreate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"id": 1, "title": "new issue",
		})
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Title:    "new issue",
		Content:  "description here",
		Kind:     "bug",
		Priority: "critical",
	}
	_, err := client.Repositories.Issues.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "new issue", receivedBody["title"])
	assert.Equal(t, "bug", receivedBody["kind"])
	assert.Equal(t, "critical", receivedBody["priority"])
	content := receivedBody["content"].(map[string]interface{})
	assert.Equal(t, "description here", content["raw"])
}

func TestIssuesGetVote_Voted(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{})
	})
	defer server.Close()

	opts := &IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	voted, _, err := client.Repositories.Issues.GetVote(opts)

	require.NoError(t, err)
	assert.True(t, voted)
}

func TestIssuesGetVote_NotVoted(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	opts := &IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	voted, _, err := client.Repositories.Issues.GetVote(opts)

	require.NoError(t, err)
	assert.False(t, voted)
}

func TestIssuesPutVote_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	err := client.Repositories.Issues.PutVote(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
}

func TestIssuesDeleteVote_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	err := client.Repositories.Issues.DeleteVote(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
}

func TestIssuesGetWatch_Watching(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{})
	})
	defer server.Close()

	opts := &IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	watching, _, err := client.Repositories.Issues.GetWatch(opts)

	require.NoError(t, err)
	assert.True(t, watching)
}

func TestIssuesPutWatch_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	err := client.Repositories.Issues.PutWatch(opts)

	require.NoError(t, err)
}

func TestIssuesDeleteWatch_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	err := client.Repositories.Issues.DeleteWatch(opts)

	require.NoError(t, err)
}

func TestIssuesGetComments_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{
				map[string]interface{}{"id": 1, "content": map[string]interface{}{"raw": "comment"}},
			},
		})
	})
	defer server.Close()

	opts := &IssueCommentsOptions{
		IssuesOptions: IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
	}
	_, err := client.Repositories.Issues.GetComments(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/comments", receivedPath)
}

func TestIssuesCreateComment_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"id": 1, "content": map[string]interface{}{"raw": "new comment"},
		})
	})
	defer server.Close()

	opts := &IssueCommentsOptions{
		IssuesOptions:  IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
		CommentContent: "new comment",
	}
	_, err := client.Repositories.Issues.CreateComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	content := receivedBody["content"].(map[string]interface{})
	assert.Equal(t, "new comment", content["raw"])
}

func TestIssuesGetComment_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": 5, "content": map[string]interface{}{"raw": "existing comment"},
		})
	})
	defer server.Close()

	opts := &IssueCommentsOptions{
		IssuesOptions: IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
		CommentID:     "5",
	}
	_, err := client.Repositories.Issues.GetComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/comments/5", receivedPath)
}

func TestIssuesUpdateComment_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": 5, "content": map[string]interface{}{"raw": "updated"},
		})
	})
	defer server.Close()

	opts := &IssueCommentsOptions{
		IssuesOptions:  IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
		CommentID:      "5",
		CommentContent: "updated",
	}
	_, err := client.Repositories.Issues.UpdateComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
}

func TestIssuesDeleteComment_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &IssueCommentsOptions{
		IssuesOptions: IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
		CommentID:     "5",
	}
	_, err := client.Repositories.Issues.DeleteComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/comments/5", receivedPath)
}

func TestIssuesGetChanges_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	opts := &IssueChangesOptions{
		IssuesOptions: IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
	}
	_, err := client.Repositories.Issues.GetChanges(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/changes", receivedPath)
}

func TestIssuesGetChange_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": "change-1",
		})
	})
	defer server.Close()

	opts := &IssueChangesOptions{
		IssuesOptions: IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
		ChangeID:      "change-1",
	}
	_, err := client.Repositories.Issues.GetChange(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/changes/change-1", receivedPath)
}

func TestBuildIssueBody(t *testing.T) {
	t.Parallel()
	issues := &Issues{}

	opts := &IssuesOptions{
		Title:     "test issue",
		Content:   "issue content",
		State:     "open",
		Kind:      "bug",
		Priority:  "critical",
		Milestone: "v1.0",
		Component: "backend",
		Version:   "2.0",
		Assignee:  "{user-uuid}",
	}

	data, err := issues.buildIssueBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))

	assert.Equal(t, "test issue", body["title"])
	assert.Equal(t, "open", body["state"])
	assert.Equal(t, "bug", body["kind"])
	assert.Equal(t, "critical", body["priority"])

	content := body["content"].(map[string]interface{})
	assert.Equal(t, "issue content", content["raw"])

	milestone := body["milestone"].(map[string]interface{})
	assert.Equal(t, "v1.0", milestone["name"])

	component := body["component"].(map[string]interface{})
	assert.Equal(t, "backend", component["name"])

	version := body["version"].(map[string]interface{})
	assert.Equal(t, "2.0", version["name"])

	assignee := body["assignee"].(map[string]interface{})
	assert.Equal(t, "{user-uuid}", assignee["uuid"])
}

func TestBuildIssueBody_MinimalFields(t *testing.T) {
	t.Parallel()
	issues := &Issues{}

	opts := &IssuesOptions{
		Title: "minimal issue",
	}

	data, err := issues.buildIssueBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))

	assert.Equal(t, "minimal issue", body["title"])
	assert.Nil(t, body["state"])
	assert.Nil(t, body["kind"])
}

func TestIssuesCreateChange_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"id": 1, "message": map[string]interface{}{"raw": "changing status"},
		})
	})
	defer server.Close()

	opts := &IssueChangesOptions{
		IssuesOptions: IssuesOptions{
			Owner:    "owner",
			RepoSlug: "repo",
			ID:       "42",
		},
		Message: "changing status",
		Changes: []struct {
			Type     string
			NewValue string
		}{
			{Type: "status", NewValue: "closed"},
		},
	}
	_, err := client.Repositories.Issues.CreateChange(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.NotNil(t, receivedBody["changes"])
}

func TestBuildCommentBody(t *testing.T) {
	t.Parallel()
	issues := &Issues{}

	opts := &IssueCommentsOptions{
		CommentContent: "test comment content",
	}

	data, err := issues.buildCommentBody(opts)

	require.NoError(t, err)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(data), &body))

	content := body["content"].(map[string]interface{})
	assert.Equal(t, "test comment content", content["raw"])
}
