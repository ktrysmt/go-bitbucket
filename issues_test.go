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
			map[string]interface{}{
				"id":       float64(1),
				"title":    "bug report",
				"state":    "open",
				"kind":     "bug",
				"priority": "major",
				"reporter": map[string]interface{}{
					"display_name": "John Doe",
					"uuid":         "{user-uuid-1}",
				},
				"content": map[string]interface{}{
					"raw":    "Found a bug",
					"markup": "markdown",
					"html":   "<p>Found a bug</p>",
				},
				"votes":      float64(3),
				"created_on": "2025-01-15T10:00:00Z",
				"updated_on": "2025-01-16T12:00:00Z",
			},
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
	require.Len(t, values, 1)
	issue := values[0].(map[string]interface{})
	assert.Equal(t, float64(1), issue["id"])
	assert.Equal(t, "bug report", issue["title"])
	assert.Equal(t, "open", issue["state"])
	assert.Equal(t, "bug", issue["kind"])
	assert.Equal(t, "major", issue["priority"])
	assert.Equal(t, float64(3), issue["votes"])
	reporter := issue["reporter"].(map[string]interface{})
	assert.Equal(t, "John Doe", reporter["display_name"])
	content := issue["content"].(map[string]interface{})
	assert.Equal(t, "Found a bug", content["raw"])
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
	assert.Contains(t, receivedQuery, "q=")
}

func TestIssuesGet_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id":       float64(42),
			"title":    "test issue",
			"state":    "open",
			"kind":     "enhancement",
			"priority": "critical",
			"votes":    float64(5),
			"reporter": map[string]interface{}{
				"display_name": "Jane Smith",
				"uuid":         "{reporter-uuid}",
			},
			"assignee": map[string]interface{}{
				"display_name": "Dev User",
				"uuid":         "{assignee-uuid}",
			},
			"content": map[string]interface{}{
				"raw":    "Detailed description",
				"markup": "markdown",
				"html":   "<p>Detailed description</p>",
			},
			"created_on": "2025-02-01T08:30:00Z",
			"updated_on": "2025-02-02T14:00:00Z",
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
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/42", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, float64(42), resultMap["id"])
	assert.Equal(t, "test issue", resultMap["title"])
	assert.Equal(t, "open", resultMap["state"])
	assert.Equal(t, "enhancement", resultMap["kind"])
	assert.Equal(t, "critical", resultMap["priority"])
	assert.Equal(t, float64(5), resultMap["votes"])
	reporter := resultMap["reporter"].(map[string]interface{})
	assert.Equal(t, "Jane Smith", reporter["display_name"])
	assert.Equal(t, "{reporter-uuid}", reporter["uuid"])
	assignee := resultMap["assignee"].(map[string]interface{})
	assert.Equal(t, "Dev User", assignee["display_name"])
	content := resultMap["content"].(map[string]interface{})
	assert.Equal(t, "Detailed description", content["raw"])
	assert.Equal(t, "markdown", content["markup"])
	assert.Equal(t, "2025-02-01T08:30:00Z", resultMap["created_on"])
	assert.Equal(t, "2025-02-02T14:00:00Z", resultMap["updated_on"])
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
			"id":         float64(42),
			"title":      "updated title",
			"state":      "closed",
			"updated_on": "2025-03-05T18:00:00Z",
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
	result, err := client.Repositories.Issues.Update(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "updated title", receivedBody["title"])
	assert.Equal(t, "closed", receivedBody["state"])
	resultMap := result.(map[string]interface{})
	assert.Equal(t, float64(42), resultMap["id"])
	assert.Equal(t, "updated title", resultMap["title"])
	assert.Equal(t, "closed", resultMap["state"])
	assert.Equal(t, "2025-03-05T18:00:00Z", resultMap["updated_on"])
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
			"id":       float64(1),
			"title":    "new issue",
			"state":    "new",
			"kind":     "bug",
			"priority": "critical",
			"reporter": map[string]interface{}{
				"display_name": "test-user",
				"uuid":         "{creator-uuid}",
			},
			"content": map[string]interface{}{
				"raw":    "description here",
				"markup": "markdown",
				"html":   "<p>description here</p>",
			},
			"votes":      float64(0),
			"created_on": "2025-03-05T10:00:00Z",
			"updated_on": "2025-03-05T10:00:00Z",
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
	result, err := client.Repositories.Issues.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "new issue", receivedBody["title"])
	assert.Equal(t, "bug", receivedBody["kind"])
	assert.Equal(t, "critical", receivedBody["priority"])
	reqContent := receivedBody["content"].(map[string]interface{})
	assert.Equal(t, "description here", reqContent["raw"])
	resultMap := result.(map[string]interface{})
	assert.Equal(t, float64(1), resultMap["id"])
	assert.Equal(t, "new issue", resultMap["title"])
	assert.Equal(t, "new", resultMap["state"])
	assert.Equal(t, "bug", resultMap["kind"])
	assert.Equal(t, "critical", resultMap["priority"])
	assert.Equal(t, float64(0), resultMap["votes"])
	reporter := resultMap["reporter"].(map[string]interface{})
	assert.Equal(t, "test-user", reporter["display_name"])
	resContent := resultMap["content"].(map[string]interface{})
	assert.Equal(t, "description here", resContent["raw"])
	assert.Equal(t, "2025-03-05T10:00:00Z", resultMap["created_on"])
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
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	err := client.Repositories.Issues.PutWatch(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/watch", receivedPath)
}

func TestIssuesDeleteWatch_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	opts := &IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	err := client.Repositories.Issues.DeleteWatch(opts)

	require.NoError(t, err)
	assert.Equal(t, "DELETE", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/watch", receivedPath)
}

func TestIssuesGetComments_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"values": []interface{}{
				map[string]interface{}{
					"id": float64(1),
					"content": map[string]interface{}{
						"raw":    "first comment",
						"markup": "markdown",
						"html":   "<p>first comment</p>",
					},
					"created_on": "2025-03-01T10:00:00Z",
				},
				map[string]interface{}{
					"id": float64(2),
					"content": map[string]interface{}{
						"raw": "second comment",
					},
				},
			},
		})
	})
	defer server.Close()

	opts := &IssueCommentsOptions{
		IssuesOptions: IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
	}
	result, err := client.Repositories.Issues.GetComments(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/comments", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	require.Len(t, values, 2)
	first := values[0].(map[string]interface{})
	assert.Equal(t, float64(1), first["id"])
	firstContent := first["content"].(map[string]interface{})
	assert.Equal(t, "first comment", firstContent["raw"])
	second := values[1].(map[string]interface{})
	assert.Equal(t, float64(2), second["id"])
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
			"id": float64(10),
			"content": map[string]interface{}{
				"raw":    "new comment",
				"markup": "markdown",
				"html":   "<p>new comment</p>",
			},
			"created_on": "2025-03-05T09:00:00Z",
		})
	})
	defer server.Close()

	opts := &IssueCommentsOptions{
		IssuesOptions:  IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
		CommentContent: "new comment",
	}
	result, err := client.Repositories.Issues.CreateComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	reqContent := receivedBody["content"].(map[string]interface{})
	assert.Equal(t, "new comment", reqContent["raw"])
	resultMap := result.(map[string]interface{})
	assert.Equal(t, float64(10), resultMap["id"])
	resContent := resultMap["content"].(map[string]interface{})
	assert.Equal(t, "new comment", resContent["raw"])
	assert.Equal(t, "2025-03-05T09:00:00Z", resultMap["created_on"])
}

func TestIssuesGetComment_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": float64(5),
			"content": map[string]interface{}{
				"raw":    "existing comment",
				"markup": "markdown",
				"html":   "<p>existing comment</p>",
			},
			"created_on": "2025-02-20T11:00:00Z",
			"updated_on": "2025-02-21T15:00:00Z",
		})
	})
	defer server.Close()

	opts := &IssueCommentsOptions{
		IssuesOptions: IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
		CommentID:     "5",
	}
	result, err := client.Repositories.Issues.GetComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/comments/5", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, float64(5), resultMap["id"])
	resContent := resultMap["content"].(map[string]interface{})
	assert.Equal(t, "existing comment", resContent["raw"])
	assert.Equal(t, "markdown", resContent["markup"])
	assert.Equal(t, "2025-02-20T11:00:00Z", resultMap["created_on"])
}

func TestIssuesUpdateComment_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedBody map[string]interface{}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": float64(5),
			"content": map[string]interface{}{
				"raw":  "updated comment",
				"html": "<p>updated comment</p>",
			},
			"updated_on": "2025-03-05T16:00:00Z",
		})
	})
	defer server.Close()

	opts := &IssueCommentsOptions{
		IssuesOptions:  IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
		CommentID:      "5",
		CommentContent: "updated comment",
	}
	result, err := client.Repositories.Issues.UpdateComment(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	reqContent := receivedBody["content"].(map[string]interface{})
	assert.Equal(t, "updated comment", reqContent["raw"])
	resultMap := result.(map[string]interface{})
	assert.Equal(t, float64(5), resultMap["id"])
	resContent := resultMap["content"].(map[string]interface{})
	assert.Equal(t, "updated comment", resContent["raw"])
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
			"values": []interface{}{
				map[string]interface{}{
					"id": float64(100),
					"changes": map[string]interface{}{
						"status": map[string]interface{}{
							"old": "open",
							"new": "closed",
						},
					},
					"message": map[string]interface{}{
						"raw": "Closing this issue",
					},
					"created_on": "2025-03-01T12:00:00Z",
				},
			},
		})
	})
	defer server.Close()

	opts := &IssueChangesOptions{
		IssuesOptions: IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
	}
	result, err := client.Repositories.Issues.GetChanges(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/changes", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	require.Len(t, values, 1)
	change := values[0].(map[string]interface{})
	assert.Equal(t, float64(100), change["id"])
	changes := change["changes"].(map[string]interface{})
	status := changes["status"].(map[string]interface{})
	assert.Equal(t, "open", status["old"])
	assert.Equal(t, "closed", status["new"])
	msg := change["message"].(map[string]interface{})
	assert.Equal(t, "Closing this issue", msg["raw"])
}

func TestIssuesGetChange_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"id": "change-1",
			"changes": map[string]interface{}{
				"priority": map[string]interface{}{
					"old": "major",
					"new": "critical",
				},
			},
			"message": map[string]interface{}{
				"raw": "Bumping priority",
			},
			"created_on": "2025-03-02T08:00:00Z",
		})
	})
	defer server.Close()

	opts := &IssueChangesOptions{
		IssuesOptions: IssuesOptions{Owner: "owner", RepoSlug: "repo", ID: "1"},
		ChangeID:      "change-1",
	}
	result, err := client.Repositories.Issues.GetChange(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/issues/1/changes/change-1", receivedPath)
	resultMap := result.(map[string]interface{})
	assert.Equal(t, "change-1", resultMap["id"])
	changes := resultMap["changes"].(map[string]interface{})
	priority := changes["priority"].(map[string]interface{})
	assert.Equal(t, "major", priority["old"])
	assert.Equal(t, "critical", priority["new"])
	msg := resultMap["message"].(map[string]interface{})
	assert.Equal(t, "Bumping priority", msg["raw"])
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

func TestIssuesCreate_ServerError(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{"message": "internal server error"},
		})
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Title:    "new issue",
	}
	result, err := client.Repositories.Issues.Create(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIssuesGet_NotFound(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "Issue not found"},
		})
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		ID:       "9999",
	}
	result, err := client.Repositories.Issues.Get(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "404")
}

func TestIssuesUpdate_Unauthorized(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "You don't have permission to update this issue"},
		})
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		ID:       "42",
		Title:    "should fail",
	}
	result, err := client.Repositories.Issues.Update(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "403")
}

func TestIssuesDelete_NotFound(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "Issue not found"},
		})
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		ID:       "9999",
	}
	result, err := client.Repositories.Issues.Delete(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "404")
}

func TestIssuesCreate_BadRequest(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{"message": "title: This field is required."},
		})
	})
	defer server.Close()

	opts := &IssuesOptions{
		Owner:    "owner",
		RepoSlug: "repo",
	}
	result, err := client.Repositories.Issues.Create(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "400")
}
