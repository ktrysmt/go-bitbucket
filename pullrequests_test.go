package bitbucket

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func realisticPRResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":    float64(1),
		"title": "Add feature X",
		"state": "OPEN",
		"author": map[string]interface{}{
			"display_name": "Jane Doe",
			"uuid":         "{user-uuid-1234}",
		},
		"source": map[string]interface{}{
			"branch": map[string]interface{}{"name": "feature-x"},
			"repository": map[string]interface{}{
				"full_name": "owner/repo",
			},
		},
		"destination": map[string]interface{}{
			"branch": map[string]interface{}{"name": "main"},
			"repository": map[string]interface{}{
				"full_name": "owner/repo",
			},
		},
		"reviewers": []interface{}{
			map[string]interface{}{"uuid": "{reviewer-uuid-1}", "display_name": "Alice"},
			map[string]interface{}{"uuid": "{reviewer-uuid-2}", "display_name": "Bob"},
		},
		"close_source_branch": true,
		"created_on":          "2025-01-15T10:30:00.000000+00:00",
		"updated_on":          "2025-01-15T12:00:00.000000+00:00",
		"comment_count":       float64(3),
		"task_count":          float64(1),
	}
}

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
		respondJSON(w, http.StatusCreated, realisticPRResponse())
	})
	defer server.Close()

	opts := &PullRequestsOptions{
		Owner:             "owner",
		RepoSlug:          "repo",
		Title:             "Add feature X",
		Description:       "Implements feature X",
		SourceBranch:      "feature-x",
		SourceRepository:  "owner/repo",
		DestinationBranch: "main",
		Reviewers:         []string{"{reviewer-uuid-1}", "{reviewer-uuid-2}"},
		CloseSourceBranch: true,
	}
	result, err := client.Repositories.PullRequests.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/", receivedPath)

	assert.Equal(t, "Add feature X", receivedBody["title"])
	assert.Equal(t, "Implements feature X", receivedBody["description"])
	assert.Equal(t, true, receivedBody["close_source_branch"])

	source := receivedBody["source"].(map[string]interface{})
	sourceBranch := source["branch"].(map[string]interface{})
	assert.Equal(t, "feature-x", sourceBranch["name"])
	sourceRepo := source["repository"].(map[string]interface{})
	assert.Equal(t, "owner/repo", sourceRepo["full_name"])

	dest := receivedBody["destination"].(map[string]interface{})
	destBranch := dest["branch"].(map[string]interface{})
	assert.Equal(t, "main", destBranch["name"])

	reviewers := receivedBody["reviewers"].([]interface{})
	assert.Len(t, reviewers, 2)
	firstReviewer := reviewers[0].(map[string]interface{})
	assert.Equal(t, "{reviewer-uuid-1}", firstReviewer["uuid"])

	resultMap := result.(map[string]interface{})
	assert.Equal(t, float64(1), resultMap["id"])
	assert.Equal(t, "Add feature X", resultMap["title"])
	assert.Equal(t, "OPEN", resultMap["state"])
}

func TestPullRequestsUpdate_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string
	var receivedBody map[string]interface{}

	updatedPR := realisticPRResponse()
	updatedPR["title"] = "Updated PR title"
	updatedPR["description"] = "Updated description"
	updatedPR["updated_on"] = "2025-01-16T08:00:00.000000+00:00"

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		bodyBytes, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bodyBytes, &receivedBody)
		respondJSON(w, http.StatusOK, updatedPR)
	})
	defer server.Close()

	opts := &PullRequestsOptions{
		Owner:             "owner",
		RepoSlug:          "repo",
		ID:                "1",
		Title:             "Updated PR title",
		Description:       "Updated description",
		DestinationBranch: "main",
	}
	result, err := client.Repositories.PullRequests.Update(opts)

	require.NoError(t, err)
	assert.Equal(t, "PUT", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1", receivedPath)

	assert.Equal(t, "Updated PR title", receivedBody["title"])
	assert.Equal(t, "Updated description", receivedBody["description"])
	dest := receivedBody["destination"].(map[string]interface{})
	destBranch := dest["branch"].(map[string]interface{})
	assert.Equal(t, "main", destBranch["name"])

	resultMap := result.(map[string]interface{})
	assert.Equal(t, float64(1), resultMap["id"])
	assert.Equal(t, "Updated PR title", resultMap["title"])
	assert.Equal(t, "2025-01-16T08:00:00.000000+00:00", resultMap["updated_on"])
}

func TestPullRequestsList_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	pr1 := realisticPRResponse()
	pr1["id"] = float64(1)
	pr1["title"] = "First PR"
	pr2 := realisticPRResponse()
	pr2["id"] = float64(2)
	pr2["title"] = "Second PR"
	pr2["state"] = "MERGED"

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{pr1, pr2}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.PullRequests.List(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 2)

	first := values[0].(map[string]interface{})
	assert.Equal(t, float64(1), first["id"])
	assert.Equal(t, "First PR", first["title"])
	assert.Equal(t, "OPEN", first["state"])
	author := first["author"].(map[string]interface{})
	assert.Equal(t, "Jane Doe", author["display_name"])

	second := values[1].(map[string]interface{})
	assert.Equal(t, float64(2), second["id"])
	assert.Equal(t, "Second PR", second["title"])
	assert.Equal(t, "MERGED", second["state"])
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
	var receivedMethod string

	prResponse := realisticPRResponse()
	prResponse["id"] = float64(42)
	prResponse["title"] = "Implement widget API"

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, prResponse)
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "42"}
	result, err := client.Repositories.PullRequests.Get(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/42", receivedPath)

	resultMap := result.(map[string]interface{})
	assert.Equal(t, float64(42), resultMap["id"])
	assert.Equal(t, "Implement widget API", resultMap["title"])
	assert.Equal(t, "OPEN", resultMap["state"])
	assert.Equal(t, true, resultMap["close_source_branch"])
	assert.Equal(t, float64(3), resultMap["comment_count"])
	assert.Equal(t, float64(1), resultMap["task_count"])
	assert.Equal(t, "2025-01-15T10:30:00.000000+00:00", resultMap["created_on"])

	author := resultMap["author"].(map[string]interface{})
	assert.Equal(t, "Jane Doe", author["display_name"])
	assert.Equal(t, "{user-uuid-1234}", author["uuid"])

	source := resultMap["source"].(map[string]interface{})
	sourceBranch := source["branch"].(map[string]interface{})
	assert.Equal(t, "feature-x", sourceBranch["name"])

	destination := resultMap["destination"].(map[string]interface{})
	destBranch := destination["branch"].(map[string]interface{})
	assert.Equal(t, "main", destBranch["name"])

	reviewers := resultMap["reviewers"].([]interface{})
	assert.Len(t, reviewers, 2)
	firstReviewer := reviewers[0].(map[string]interface{})
	assert.Equal(t, "{reviewer-uuid-1}", firstReviewer["uuid"])
	assert.Equal(t, "Alice", firstReviewer["display_name"])
}

func TestPullRequestsGetByCommit_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	pr := realisticPRResponse()

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{pr}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", Commit: "abc123"}
	result, err := client.Repositories.PullRequests.GetByCommit(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/commit/abc123/pullrequests/", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	first := values[0].(map[string]interface{})
	assert.Equal(t, float64(1), first["id"])
	assert.Equal(t, "OPEN", first["state"])
}

func TestPullRequestsGetCommits_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"hash":    "abc123def456",
				"message": "Initial commit",
				"author":  map[string]interface{}{"raw": "Jane Doe <jane@example.com>"},
			},
			map[string]interface{}{
				"hash":    "789ghi012jkl",
				"message": "Add tests",
				"author":  map[string]interface{}{"raw": "Jane Doe <jane@example.com>"},
			},
		}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.PullRequests.GetCommits(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/commits/", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 2)
	firstCommit := values[0].(map[string]interface{})
	assert.Equal(t, "abc123def456", firstCommit["hash"])
	assert.Equal(t, "Initial commit", firstCommit["message"])
}

func TestPullRequestsActivities_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"update": map[string]interface{}{
					"state":  "OPEN",
					"title":  "Add feature X",
					"author": map[string]interface{}{"display_name": "Jane Doe"},
				},
			},
		}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.PullRequests.Activities(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/activity", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	activity := values[0].(map[string]interface{})
	update := activity["update"].(map[string]interface{})
	assert.Equal(t, "OPEN", update["state"])
}

func TestPullRequestsActivity_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	activityResponse := map[string]interface{}{
		"values": []interface{}{
			map[string]interface{}{
				"comment": map[string]interface{}{
					"id":      float64(100),
					"content": map[string]interface{}{"raw": "Looks good"},
					"user":    map[string]interface{}{"display_name": "Alice"},
				},
			},
		},
	}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, activityResponse)
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.PullRequests.Activity(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/activity", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	first := values[0].(map[string]interface{})
	comment := first["comment"].(map[string]interface{})
	assert.Equal(t, float64(100), comment["id"])
}

func TestPullRequestsCommits_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{
				"hash":    "deadbeef1234",
				"message": "Fix bug in handler",
				"date":    "2025-01-15T11:00:00+00:00",
				"author":  map[string]interface{}{"raw": "Bob <bob@example.com>"},
			},
		}))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.PullRequests.Commits(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/commits", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
	commit := values[0].(map[string]interface{})
	assert.Equal(t, "deadbeef1234", commit["hash"])
	assert.Equal(t, "Fix bug in handler", commit["message"])
}

func TestPullRequestsMerge_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	mergedPR := realisticPRResponse()
	mergedPR["state"] = "MERGED"
	mergedPR["merge_commit"] = map[string]interface{}{"hash": "merged123abc"}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, mergedPR)
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1", Message: "Merging PR"}
	result, err := client.Repositories.PullRequests.Merge(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/merge", receivedPath)

	resultMap := result.(map[string]interface{})
	assert.Equal(t, "MERGED", resultMap["state"])
	assert.Equal(t, float64(1), resultMap["id"])
	mergeCommit := resultMap["merge_commit"].(map[string]interface{})
	assert.Equal(t, "merged123abc", mergeCommit["hash"])
}

func TestPullRequestsDecline_Success(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedPath string

	declinedPR := realisticPRResponse()
	declinedPR["state"] = "DECLINED"

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		respondJSON(w, http.StatusOK, declinedPR)
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.PullRequests.Decline(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/decline", receivedPath)

	resultMap := result.(map[string]interface{})
	assert.Equal(t, "DECLINED", resultMap["state"])
	assert.Equal(t, float64(1), resultMap["id"])
	assert.Equal(t, "Add feature X", resultMap["title"])
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
	patchContent := "diff --git a/file.txt b/file.txt\nindex 1234567..abcdefg 100644\n--- a/file.txt\n+++ b/file.txt\n@@ -1,3 +1,4 @@\n line1\n+new line\n line2\n line3"

	var receivedPath string
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(patchContent))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.PullRequests.Patch(opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/patch", receivedPath)

	rc := result.(io.ReadCloser)
	defer func() { _ = rc.Close() }()
	body, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Contains(t, string(body), "diff --git a/file.txt")
	assert.Contains(t, string(body), "+new line")
}

func TestPullRequestsDiff_Success(t *testing.T) {
	t.Parallel()
	diffContent := "diff --git a/main.go b/main.go\n--- a/main.go\n+++ b/main.go\n@@ -10,6 +10,7 @@\n func main() {\n+\tfmt.Println(\"hello\")\n }"

	var receivedPath string
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(diffContent))
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.PullRequests.Diff(opts)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "/2.0/repositories/owner/repo/pullrequests/1/diff", receivedPath)

	rc := result.(io.ReadCloser)
	defer func() { _ = rc.Close() }()
	body, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.Contains(t, string(body), "diff --git a/main.go")
	assert.Contains(t, string(body), "fmt.Println")
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

func TestPullRequestsCreate_ServerError(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Internal server error",
			},
		})
	})
	defer server.Close()

	opts := &PullRequestsOptions{
		Owner:        "owner",
		RepoSlug:     "repo",
		Title:        "My PR",
		SourceBranch: "feature",
	}
	result, err := client.Repositories.PullRequests.Create(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPullRequestsGet_NotFound(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Pull request not found",
			},
		})
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "99999"}
	result, err := client.Repositories.PullRequests.Get(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPullRequestsUpdate_Forbidden(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "You don't have permission to update this pull request",
			},
		})
	})
	defer server.Close()

	opts := &PullRequestsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		ID:       "1",
		Title:    "Attempted update",
	}
	result, err := client.Repositories.PullRequests.Update(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPullRequestsMerge_Conflict(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusConflict, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Pull request has conflicts that must be resolved",
			},
		})
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "1"}
	result, err := client.Repositories.PullRequests.Merge(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPullRequestsDecline_NotFound(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Pull request not found",
			},
		})
	})
	defer server.Close()

	opts := &PullRequestsOptions{Owner: "owner", RepoSlug: "repo", ID: "99999"}
	result, err := client.Repositories.PullRequests.Decline(opts)

	assert.Error(t, err)
	assert.Nil(t, result)
}
