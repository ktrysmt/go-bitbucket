package bitbucket

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadsList_Success(t *testing.T) {
	t.Parallel()
	var receivedPath string
	var receivedMethod string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"name": "release-1.0.zip", "size": 1024},
		}))
	})
	defer server.Close()

	opts := &DownloadsOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.Downloads.List(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/downloads", receivedPath)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 1)
}

func TestDownloadsList_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusForbidden, map[string]interface{}{
			"error": map[string]interface{}{"message": "forbidden"},
		})
	})
	defer server.Close()

	opts := &DownloadsOptions{Owner: "owner", RepoSlug: "repo"}
	_, err := client.Repositories.Downloads.List(opts)

	assert.Error(t, err)
}

func TestDownloadsCreate_BothFilesAndFilenameError(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	opts := &DownloadsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		FileName: "file.zip",
		Files:    []File{{Path: "other.zip", Name: "other.zip"}},
	}
	_, err := client.Repositories.Downloads.Create(opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can't specify both files and filename")
}

func TestDownloadsCreate_WithFile(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedContentType string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusCreated)
	})
	defer server.Close()

	tmpFile, err := os.CreateTemp("", "test-download-*.txt")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, err = tmpFile.WriteString("test file content")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	opts := &DownloadsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Files:    []File{{Path: tmpFile.Name(), Name: "upload.txt"}},
	}
	_, err = client.Repositories.Downloads.Create(opts)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Contains(t, receivedContentType, "multipart/form-data")
}

func TestDownloadsCreate_FileNotFound(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	opts := &DownloadsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		FileName: "/nonexistent/path/file.txt",
	}
	_, err := client.Repositories.Downloads.Create(opts)

	assert.Error(t, err)
}
