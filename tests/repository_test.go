package tests

import (
	"os"
	"testing"

	_ "github.com/k0kubun/pp"
	"github.com/ktrysmt/go-bitbucket"
)

func TestGetRepositoryRepositories(t *testing.T) {

	user := os.Getenv("BITBUCKET_TEST_USERNAME")
	pass := os.Getenv("BITBUCKET_TEST_PASSWORD")
	owner := os.Getenv("BITBUCKET_TEST_OWNER")
	repo := os.Getenv("BITBUCKET_TEST_REPOSLUG")

	if user == "" {
		t.Error("BITBUCKET_TEST_USERNAME is empty.")
	}
	if pass == "" {
		t.Error("BITBUCKET_TEST_PASSWORD is empty.")
	}
	if owner == "" {
		t.Error("BITBUCKET_TEST_OWNER is empty.")
	}
	if repo == "" {
		t.Error("BITBUCKET_TEST_REPOSLUG is empty.")
	}

	c := bitbucket.NewBasicAuth(user, pass)

	opt := &bitbucket.RepositoryOptions{
		Owner:    owner,
		RepoSlug: repo,
	}

	res, err := c.Repositories.Repository.Get(opt)
	if err != nil {
		t.Error("The repository is not found.")
	}

	if res.Full_name != owner+"/"+repo {
		t.Error("Cannot catch repos full name.")
	}
}

func TestGetRepositoryPipelineVariables(t *testing.T) {

	user := os.Getenv("BITBUCKET_TEST_USERNAME")
	pass := os.Getenv("BITBUCKET_TEST_PASSWORD")
	owner := os.Getenv("BITBUCKET_TEST_OWNER")
	repo := os.Getenv("BITBUCKET_TEST_REPOSLUG")

	if user == "" {
		t.Error("BITBUCKET_TEST_USERNAME is empty.")
	}
	if pass == "" {
		t.Error("BITBUCKET_TEST_PASSWORD is empty.")
	}
	if owner == "" {
		t.Error("BITBUCKET_TEST_OWNER is empty.")
	}
	if repo == "" {
		t.Error("BITBUCKET_TEST_REPOSLUG is empty.")
	}

	c := bitbucket.NewBasicAuth(user, pass)

	opt := &bitbucket.RepositoryPipelineVariablesOptions{
		Owner:    owner,
		RepoSlug: repo,
	}

	res, err := c.Repositories.Repository.ListPipelineVariables(opt)
	if err != nil {
		t.Error(err)
	}

	if res == nil {
		t.Error("Cannot list pipeline variables")
	}
}

func TestDeleteRepositoryPipelineVariables(t *testing.T) {

	user := os.Getenv("BITBUCKET_TEST_USERNAME")
	pass := os.Getenv("BITBUCKET_TEST_PASSWORD")
	owner := os.Getenv("BITBUCKET_TEST_OWNER")
	repo := os.Getenv("BITBUCKET_TEST_REPOSLUG")

	if user == "" {
		t.Error("BITBUCKET_TEST_USERNAME is empty.")
	}
	if pass == "" {
		t.Error("BITBUCKET_TEST_PASSWORD is empty.")
	}
	if owner == "" {
		t.Error("BITBUCKET_TEST_OWNER is empty.")
	}
	if repo == "" {
		t.Error("BITBUCKET_TEST_REPOSLUG is empty.")
	}

	c := bitbucket.NewBasicAuth(user, pass)

	variable := &bitbucket.RepositoryPipelineVariableOptions{
		Owner:    owner,
		RepoSlug: repo,
		Key:      "test_key_to_delete",
		Value:    "Some value to delete",
		Secured:  false,
	}

	res, err := c.Repositories.Repository.AddPipelineVariable(variable)
	if err != nil {
		t.Error(err)
	}

	opt := &bitbucket.RepositoryPipelineVariableDeleteOptions{
		Owner:    owner,
		RepoSlug: repo,
		Uuid:     res.Uuid,
	}

	// On success the delete API doesn't return any content (HTTP status 204)
	_, err = c.Repositories.Repository.DeletePipelineVariable(opt)
	if err != nil {
		t.Error(err)
	}
}

func TestGetRepositoryRefs(t *testing.T) {

	user := os.Getenv("BITBUCKET_TEST_USERNAME")
	pass := os.Getenv("BITBUCKET_TEST_PASSWORD")
	owner := os.Getenv("BITBUCKET_TEST_OWNER")
	repo := os.Getenv("BITBUCKET_TEST_REPOSLUG")

	if user == "" {
		t.Error("BITBUCKET_TEST_USERNAME is empty.")
	}
	if pass == "" {
		t.Error("BITBUCKET_TEST_PASSWORD is empty.")
	}
	if owner == "" {
		t.Error("BITBUCKET_TEST_OWNER is empty.")
	}
	if repo == "" {
		t.Error("BITBUCKET_TEST_REPOSLUG is empty.")
	}

	c := bitbucket.NewBasicAuth(user, pass)

	opt := &bitbucket.RepositoryBranchCreationOptions{
		Owner:    owner,
		RepoSlug: repo,
		Name:     "TestGetRepoRefsBranch",
		Target:   bitbucket.RepositoryBranchTarget{Hash: "master"},
	}

	_, err := c.Repositories.Repository.CreateBranch(opt)
	if err != nil {
		t.Error("Could not create new branch", err)
	}

	refOpts := &bitbucket.RepositoryRefOptions{
		Owner:    owner,
		RepoSlug: repo,
	}

	resRefs, err := c.Repositories.Repository.ListRefs(refOpts)
	if err != nil {
		t.Error("The refs is not found.")
	}

	expected := struct {
		n string
		t string
	}{}

	for _, ref := range resRefs.Refs {
		for k, v := range ref {
			// kCopy := k
			vCopy := v
			if val, ok := vCopy.(string); ok {
				if k == "name" && val == "TestGetRepoRefsBranch" {
					expected.n = val
				}
			}
			if val, ok := vCopy.(string); ok {
				if k == "type" && val == "branch" {
					expected.t = val
				}
			}
		}
	}

	if !(expected.n == "TestGetRepoRefsBranch" && expected.t == "branch") {
		t.Error("Could not list refs/branch that was created in test setup")
	}
}
