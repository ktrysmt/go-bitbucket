package bitbucket

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDiff_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string
	diffContent := "diff --git a/file.txt b/file.txt\nindex 1234567..abcdefg 100644\n--- a/file.txt\n+++ b/file.txt\n@@ -1,3 +1,4 @@\n+new line\n existing\n"

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(diffContent))
	})
	defer server.Close()

	opts := &DiffOptions{Owner: "owner", RepoSlug: "repo", Spec: "abc123"}
	result, err := client.Repositories.Diff.GetDiff(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/diff/abc123", receivedPath)
	assert.NotNil(t, result)
}

func TestGetDiff_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	})
	defer server.Close()

	opts := &DiffOptions{Owner: "owner", RepoSlug: "repo", Spec: "bad-spec"}
	_, err := client.Repositories.Diff.GetDiff(opts)

	assert.Error(t, err)
}

func TestGetDiff_WithOptions(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("diff content"))
	})
	defer server.Close()

	opts := &DiffOptions{
		Owner:      "owner",
		RepoSlug:   "repo",
		Spec:       "abc123",
		Context:    5,
		Path:       "src/main.go",
		Whitespace: true,
		Topic:      true,
	}
	_, err := client.Repositories.Diff.GetDiff(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "context=5")
	assert.Contains(t, receivedQuery, "path=src%2Fmain.go")
	assert.Contains(t, receivedQuery, "ignore_whitespace=true")
	assert.Contains(t, receivedQuery, "topic=true")
}

func TestGetPatch_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	patchContent := "From abc123\nSubject: [PATCH] fix bug\n---\n file.txt | 2 +-\n 1 file changed\n"

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(patchContent))
	})
	defer server.Close()

	opts := &DiffOptions{Owner: "owner", RepoSlug: "repo", Spec: "abc123"}
	result, err := client.Repositories.Diff.GetPatch(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/patch/abc123", receivedPath)
	assert.NotNil(t, result)
}

func TestGetPatch_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	})
	defer server.Close()

	opts := &DiffOptions{Owner: "owner", RepoSlug: "repo", Spec: "bad-spec"}
	_, err := client.Repositories.Diff.GetPatch(opts)

	assert.Error(t, err)
}

func TestGetDiffStat_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"page":    1,
			"pagelen": 10,
			"size":    2,
			"values": []interface{}{
				map[string]interface{}{
					"type":          "diffstat",
					"status":        "modified",
					"lines_removed": 5,
					"lines_added":   10,
					"old": map[string]interface{}{
						"path": "src/main.go",
						"type": "commit_file",
					},
					"new": map[string]interface{}{
						"path": "src/main.go",
						"type": "commit_file",
					},
				},
			},
		})
	})
	defer server.Close()

	opts := &DiffStatOptions{Owner: "owner", RepoSlug: "repo", Spec: "abc123"}
	result, err := client.Repositories.Diff.GetDiffStat(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/diffstat/abc123", receivedPath)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Pagelen)
	assert.Equal(t, 2, result.Size)
	assert.Len(t, result.DiffStats, 1)

	ds := result.DiffStats[0]
	assert.Equal(t, "diffstat", ds.Type)
	assert.Equal(t, "modified", ds.Status)
	assert.Equal(t, 5, ds.LinesRemoved)
	assert.Equal(t, 10, ds.LinedAdded)
	assert.Equal(t, "src/main.go", ds.Old["path"])
	assert.Equal(t, "commit_file", ds.Old["type"])
	assert.Equal(t, "src/main.go", ds.New["path"])
	assert.Equal(t, "commit_file", ds.New["type"])
}

func TestGetDiffStat_WithOptions(t *testing.T) {
	t.Parallel()
	var receivedQuery string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"values": []interface{}{},
		})
	})
	defer server.Close()

	opts := &DiffStatOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Spec:     "abc123",
		PageNum:  2,
		Pagelen:  50,
		MaxDepth: 3,
		Fields:   []string{"values.status", "values.lines_added"},
	}
	_, err := client.Repositories.Diff.GetDiffStat(opts)

	require.NoError(t, err)
	assert.Contains(t, receivedQuery, "page=2")
	assert.Contains(t, receivedQuery, "pagelen=50")
	assert.Contains(t, receivedQuery, "max_depth=3")
	assert.Contains(t, receivedQuery, "fields=values.status")
}

func TestGetDiffStat_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	})
	defer server.Close()

	opts := &DiffStatOptions{Owner: "owner", RepoSlug: "repo", Spec: "bad"}
	_, err := client.Repositories.Diff.GetDiffStat(opts)

	assert.Error(t, err)
}

func TestDecodeDiffStat_Success(t *testing.T) {
	t.Parallel()
	input := `{"page":1,"pagelen":10,"size":1,"next":"https://api.bitbucket.org/next","previous":"https://api.bitbucket.org/prev","values":[{"type":"diffstat","status":"added","lines_removed":0,"lines_added":5,"old":{"path":"old.go","type":"commit_file"},"new":{"path":"new.go","type":"commit_file"}}]}`

	result, err := decodeDiffStat(input)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Pagelen)
	assert.Equal(t, 1, result.Size)
	assert.Equal(t, "https://api.bitbucket.org/next", result.Next)
	assert.Equal(t, "https://api.bitbucket.org/prev", result.Previous)
	require.Len(t, result.DiffStats, 1)

	ds := result.DiffStats[0]
	assert.Equal(t, "diffstat", ds.Type)
	assert.Equal(t, "added", ds.Status)
	assert.Equal(t, 0, ds.LinesRemoved)
	assert.Equal(t, 5, ds.LinedAdded)
	assert.Equal(t, "old.go", ds.Old["path"])
	assert.Equal(t, "commit_file", ds.Old["type"])
	assert.Equal(t, "new.go", ds.New["path"])
	assert.Equal(t, "commit_file", ds.New["type"])
}

func TestDecodeDiffStat_InvalidJSON(t *testing.T) {
	t.Parallel()
	_, err := decodeDiffStat("not json")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DiffStat decode error")
}

func TestCleanFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{"single field", []string{"values.status"}, "values.status"},
		{"multiple fields", []string{"values.status", "values.lines_added"}, "values.status,values.lines_added"},
		{"with spaces", []string{"values.status ", " values.lines_added"}, "values.status,values.lines_added"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := cleanFields(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
