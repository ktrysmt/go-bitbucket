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
			map[string]interface{}{"name": "release-1.0.zip", "size": float64(1024)},
			map[string]interface{}{"name": "release-2.0.zip", "size": float64(2048)},
		}))
	})
	defer server.Close()

	opts := &DownloadsOptions{Owner: "owner", RepoSlug: "repo"}
	result, err := client.Repositories.Downloads.List(opts)

	require.NoError(t, err)
	assert.Equal(t, "GET", receivedMethod)
	assert.Equal(t, "/2.0/repositories/owner/repo/downloads", receivedPath)

	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "result should be a map")
	assert.Equal(t, float64(1), resultMap["page"])
	assert.Equal(t, float64(10), resultMap["pagelen"])
	assert.Equal(t, float64(2), resultMap["size"])

	values, ok := resultMap["values"].([]interface{})
	require.True(t, ok, "values should be a slice")
	require.Len(t, values, 2)

	first, ok := values[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "release-1.0.zip", first["name"])
	assert.Equal(t, float64(1024), first["size"])

	second, ok := values[1].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "release-2.0.zip", second["name"])
	assert.Equal(t, float64(2048), second["size"])
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
	var receivedPath string
	var uploadedContent string
	var uploadedFilename string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		receivedContentType = r.Header.Get("Content-Type")

		err := r.ParseMultipartForm(10 << 20)
		if err == nil {
			file, header, ferr := r.FormFile("upload.txt")
			if ferr == nil {
				defer func() { _ = file.Close() }()
				uploadedFilename = header.Filename
				buf := make([]byte, 1024)
				n, _ := file.Read(buf)
				uploadedContent = string(buf[:n])
			}
		}
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
	assert.Equal(t, "/2.0/repositories/owner/repo/downloads", receivedPath)
	assert.Contains(t, receivedContentType, "multipart/form-data")
	assert.Equal(t, "upload.txt", uploadedFilename)
	assert.Equal(t, "test file content", uploadedContent)
}

func TestDownloadsCreate_ServerError(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": map[string]interface{}{"message": "internal server error"},
		})
	})
	defer server.Close()

	tmpFile, err := os.CreateTemp("", "test-download-error-*.txt")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, err = tmpFile.WriteString("some content")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	opts := &DownloadsOptions{
		Owner:    "owner",
		RepoSlug: "repo",
		Files:    []File{{Path: tmpFile.Name(), Name: "upload.txt"}},
	}
	result, createErr := client.Repositories.Downloads.Create(opts)

	assert.Nil(t, result)
	require.Error(t, createErr)
	unexpectedErr, ok := createErr.(*UnexpectedResponseStatusError)
	require.True(t, ok, "error should be *UnexpectedResponseStatusError")
	assert.Equal(t, "500 Internal Server Error", unexpectedErr.Status)
	assert.Contains(t, string(unexpectedErr.Body), "internal server error")
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
