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

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("diff --git a/file.txt b/file.txt\n"))
	})
	defer server.Close()

	opts := &DiffOptions{Owner: "owner", RepoSlug: "repo", Spec: "abc123"}
	_, err := client.Repositories.Diff.GetDiff(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/diff/abc123", receivedPath)
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

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("patch content"))
	})
	defer server.Close()

	opts := &DiffOptions{Owner: "owner", RepoSlug: "repo", Spec: "abc123"}
	_, err := client.Repositories.Diff.GetPatch(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/patch/abc123", receivedPath)
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
				},
			},
		})
	})
	defer server.Close()

	opts := &DiffStatOptions{Owner: "owner", RepoSlug: "repo", Spec: "abc123"}
	result, err := client.Repositories.Diff.GetDiffStat(opts)

	require.NoError(t, err)
	assert.Equal(t, "/2.0/repositories/owner/repo/diffstat/abc123", receivedPath)
	assert.Equal(t, 2, result.Size)
	assert.Len(t, result.DiffStats, 1)
	assert.Equal(t, "modified", result.DiffStats[0].Status)
	assert.Equal(t, 5, result.DiffStats[0].LinesRemoved)
	assert.Equal(t, 10, result.DiffStats[0].LinedAdded)
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
	input := `{"page":1,"pagelen":10,"size":1,"values":[{"type":"diffstat","status":"added","lines_removed":0,"lines_added":5}]}`

	result, err := decodeDiffStat(input)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Len(t, result.DiffStats, 1)
	assert.Equal(t, "added", result.DiffStats[0].Status)
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
