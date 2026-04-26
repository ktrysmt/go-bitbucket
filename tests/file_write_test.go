package tests

import (
	"os"
	"testing"
	"time"

	"github.com/ktrysmt/go-bitbucket"
)

var (
	fileTestUser  = os.Getenv("BITBUCKET_TEST_USERNAME")
	fileTestPass  = os.Getenv("BITBUCKET_TEST_PASSWORD")
	fileTestOwner = os.Getenv("BITBUCKET_TEST_OWNER")
	fileTestRepo  = os.Getenv("BITBUCKET_TEST_REPOSLUG")
)

func setupFileTest(t *testing.T) *bitbucket.Client {
	if fileTestUser == "" {
		t.Skip("BITBUCKET_TEST_USERNAME is empty.")
	}
	if fileTestPass == "" {
		t.Skip("BITBUCKET_TEST_PASSWORD is empty.")
	}
	if fileTestOwner == "" {
		t.Skip("BITBUCKET_TEST_OWNER is empty.")
	}
	if fileTestRepo == "" {
		t.Skip("BITBUCKET_TEST_REPOSLUG is empty.")
	}

	c, err := bitbucket.NewBasicAuth(fileTestUser, fileTestPass)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestWriteFileRaw(t *testing.T) {
	c := setupFileTest(t)

	timestamp := time.Now().Format("20060102150405")
	filePath := "test-raw-file-" + timestamp + ".txt"
	fileContent := "Hello from WriteFileRaw test at " + timestamp

	opt := &bitbucket.RepositoryRawFileWriteOptions{
		Owner:    fileTestOwner,
		RepoSlug: fileTestRepo,
		Files: []bitbucket.RepositoryRawFileContent{
			{
				Path:    filePath,
				Content: fileContent,
			},
		},
		Message: "Test WriteFileRaw: add " + filePath,
		Branch:  "master",
	}

	err := c.Repositories.Repository.WriteFileRaw(opt)
	if err != nil {
		t.Fatalf("WriteFileRaw failed: %v", err)
	}

	// Verify the file was written by reading it back
	blobOpt := &bitbucket.RepositoryBlobOptions{
		Owner:    fileTestOwner,
		RepoSlug: fileTestRepo,
		Ref:      "master",
		Path:     filePath,
	}

	blob, err := c.Repositories.Repository.GetFileBlob(blobOpt)
	if err != nil {
		t.Fatalf("GetFileBlob failed: %v", err)
	}

	if string(blob.Content) != fileContent {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", fileContent, string(blob.Content))
	}
}

func TestWriteFileRawMultipleFiles(t *testing.T) {
	c := setupFileTest(t)

	timestamp := time.Now().Format("20060102150405")
	file1Path := "test-multi-1-" + timestamp + ".txt"
	file1Content := "Content 1 at " + timestamp
	file2Path := "test-multi-2-" + timestamp + ".txt"
	file2Content := "Content 2 at " + timestamp

	opt := &bitbucket.RepositoryRawFileWriteOptions{
		Owner:    fileTestOwner,
		RepoSlug: fileTestRepo,
		Files: []bitbucket.RepositoryRawFileContent{
			{Path: file1Path, Content: file1Content},
			{Path: file2Path, Content: file2Content},
		},
		Message: "Test WriteFileRaw: add multiple files at " + timestamp,
		Branch:  "master",
	}

	err := c.Repositories.Repository.WriteFileRaw(opt)
	if err != nil {
		t.Fatalf("WriteFileRaw with multiple files failed: %v", err)
	}

	// Verify file 1
	blob1, err := c.Repositories.Repository.GetFileBlob(&bitbucket.RepositoryBlobOptions{
		Owner:    fileTestOwner,
		RepoSlug: fileTestRepo,
		Ref:      "master",
		Path:     file1Path,
	})
	if err != nil {
		t.Fatalf("GetFileBlob for file1 failed: %v", err)
	}
	if string(blob1.Content) != file1Content {
		t.Errorf("File1 content mismatch. Expected: %s, Got: %s", file1Content, string(blob1.Content))
	}

	// Verify file 2
	blob2, err := c.Repositories.Repository.GetFileBlob(&bitbucket.RepositoryBlobOptions{
		Owner:    fileTestOwner,
		RepoSlug: fileTestRepo,
		Ref:      "master",
		Path:     file2Path,
	})
	if err != nil {
		t.Fatalf("GetFileBlob for file2 failed: %v", err)
	}
	if string(blob2.Content) != file2Content {
		t.Errorf("File2 content mismatch. Expected: %s, Got: %s", file2Content, string(blob2.Content))
	}
}

func TestWriteFileRawWithDelete(t *testing.T) {
	c := setupFileTest(t)

	timestamp := time.Now().Format("20060102150405")

	// First, create a file to delete
	createPath := "test-to-delete-" + timestamp + ".txt"
	createOpt := &bitbucket.RepositoryRawFileWriteOptions{
		Owner:    fileTestOwner,
		RepoSlug: fileTestRepo,
		Files: []bitbucket.RepositoryRawFileContent{
			{Path: createPath, Content: "File to be deleted"},
		},
		Message: "Test: create file to delete at " + timestamp,
		Branch:  "master",
	}

	err := c.Repositories.Repository.WriteFileRaw(createOpt)
	if err != nil {
		t.Fatalf("WriteFileRaw (create) failed: %v", err)
	}

	// Now delete the file
	deleteOpt := &bitbucket.RepositoryRawFileWriteOptions{
		Owner:         fileTestOwner,
		RepoSlug:      fileTestRepo,
		FilesToDelete: []string{createPath},
		Message:       "Test: delete file at " + timestamp,
		Branch:        "master",
	}

	err = c.Repositories.Repository.WriteFileRaw(deleteOpt)
	if err != nil {
		t.Fatalf("WriteFileRaw (delete) failed: %v", err)
	}

	// Verify the file was deleted (should return an error)
	_, err = c.Repositories.Repository.GetFileBlob(&bitbucket.RepositoryBlobOptions{
		Owner:    fileTestOwner,
		RepoSlug: fileTestRepo,
		Ref:      "master",
		Path:     createPath,
	})
	if err == nil {
		t.Error("Expected error when getting deleted file, but got none")
	}
}
