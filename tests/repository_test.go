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

func TestCreateRepositoryRepositories(t *testing.T) {
	user := os.Getenv("BITBUCKET_TEST_USERNAME")
	pass := os.Getenv("BITBUCKET_TEST_PASSWORD")
	owner := os.Getenv("BITBUCKET_TEST_OWNER")

	if user == "" {
		t.Error("BITBUCKET_TEST_USERNAME is empty.")
	}
	if pass == "" {
		t.Error("BITBUCKET_TEST_PASSWORD is empty.")
	}
	if owner == "" {
		t.Error("BITBUCKET_TEST_OWNER is empty.")
	}

	c := bitbucket.NewBasicAuth(user, pass)

	// Create project - needed prior to creating repo
	projOpt := &bitbucket.ProjectOptions{
		Owner:     owner,
		Name:      "go-bitbucket-test-project",
		Key:       "GO_BB_TEST_PROJECT",
		IsPrivate: true,
	}
	project, err := c.Workspaces.CreateProject(projOpt)
	if err != nil {
		t.Error("The project could not be created.", err)
	}

	repoSlug := "go-bb-test-repo-create"
	forkPolicy := "no_forks"
	repoOpt := &bitbucket.RepositoryOptions{
		Owner:      owner,
		RepoSlug:   repoSlug,
		ForkPolicy: forkPolicy,
		Project:    project.Key,
		IsPrivate:  "true",
	}

	res, err := c.Repositories.Repository.Create(repoOpt)
	if err != nil {
		t.Error("The project could not be created.", err)
	}

	if res.Full_name != owner+"/"+repoSlug {
		t.Error("The repository `Full_name` attribute does not match the expected value.")
	}
	if res.Fork_policy != forkPolicy {
		t.Error("The repository `Fork_policy` attribute does not match the expected value.")
	}

	// Clean up
	_, err = c.Repositories.Repository.Delete(repoOpt)
	if err != nil {
		t.Error("The repository could not be deleted.", err)
	}

	_, err = c.Workspaces.DeleteProject(projOpt)
	if err != nil {
		t.Error("The project could not be deleted.", err)
	}
}

func TestRepositoryUpdateForkPolicy(t *testing.T) {
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
		t.Error("The repository is not found.", err)
	}

	forkPolicy := "allow_forks"
	opt = &bitbucket.RepositoryOptions{
		Uuid:       res.Uuid,
		Owner:      owner,
		RepoSlug:   res.Slug,
		ForkPolicy: forkPolicy,
	}
	res, err = c.Repositories.Repository.Update(opt)
	if err != nil {
		t.Error("The repository could not be updated.", err)
	}

	if res.Fork_policy != forkPolicy {
		t.Errorf("The repository's fork_policy did not match the expected: '%s'.", forkPolicy)
	}

	forkPolicy = "no_public_forks"
	opt = &bitbucket.RepositoryOptions{
		Uuid:       res.Uuid,
		Owner:      owner,
		RepoSlug:   res.Slug,
		ForkPolicy: forkPolicy,
	}
	res, err = c.Repositories.Repository.Update(opt)
	if err != nil {
		t.Error("The repository could not be updated.", err)
	}

	if res.Fork_policy != forkPolicy {
		t.Errorf("The repository's fork_policy did not match the expected: '%s'.", forkPolicy)
	}

	forkPolicy = "no_forks"
	opt = &bitbucket.RepositoryOptions{
		Uuid:       res.Uuid,
		Owner:      owner,
		RepoSlug:   res.Slug,
		ForkPolicy: forkPolicy,
	}
	res, err = c.Repositories.Repository.Update(opt)
	if err != nil {
		t.Error("The repository could not be updated.", err)
	}

	if res.Fork_policy != forkPolicy {
		t.Errorf("The repository's fork_policy did not match the expected: '%s'.", forkPolicy)
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

// This test tests the CreateBranch, CreateTag, ListRefs, DeleteBranch, and DeleteTag repo methods
func TestGetRepositoryRefs(t *testing.T) {
	// Create Branch, List Refs and search for Branch that was created,
	// then test successful Branch deletion
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
		Target:   bitbucket.RepositoryRefTarget{Hash: "master"},
	}

	_, err := c.Repositories.Repository.CreateBranch(opt)
	if err != nil {
		t.Error("Could not create new branch", err)
	}

	refBranchOpts := &bitbucket.RepositoryRefOptions{
		Owner:    owner,
		RepoSlug: repo,
	}

	resBranchRefs, err := c.Repositories.Repository.ListRefs(refBranchOpts)
	if err != nil {
		t.Error("The refs/branch is not found.")
	}

	branchExpected := struct {
		n string
		t string
	}{}

	for _, ref := range resBranchRefs.Refs {
		for k, v := range ref {
			// kCopy := k
			vCopy := v
			if val, ok := vCopy.(string); ok {
				if k == "name" && val == "TestGetRepoRefsBranch" {
					branchExpected.n = val
				}
			}
			if val, ok := vCopy.(string); ok {
				if k == "type" && val == "branch" {
					branchExpected.t = val
				}
			}
		}
	}

	if !(branchExpected.n == "TestGetRepoRefsBranch" && branchExpected.t == "branch") {
		t.Error("Could not list refs/branch that was created in test setup")
	}

	delBranchOpt := &bitbucket.RepositoryBranchOptions{
		Owner:    owner,
		RepoSlug: repo,
		Name:     "TestGetRepoRefsBranch",
	}
	resBranchDel, err := c.Repositories.Repository.DeleteBranch(delBranchOpt)
	if err != nil {
		t.Error("Could not delete branch")
	}

	if mp, ok := resBranchDel.(map[string]interface{}); ok {
		for k, v := range mp {
			if k == "type" && v == "error" {
				t.Error("Delete branch returned an error, when it should have successfully" +
					"delete the test branch created during test setup")
			}
		}
	}

	// Create Tag, List Refs and search for Tag that was created,
	// then test successful Tag deletion
	tagOpt := &bitbucket.RepositoryTagCreationOptions{
		Owner:    owner,
		RepoSlug: repo,
		Name:     "TestGetRepoRefsTag",
		Target:   bitbucket.RepositoryRefTarget{Hash: "master"},
	}

	_, err = c.Repositories.Repository.CreateTag(tagOpt)
	if err != nil {
		t.Error("Could not create new tag", err)
	}

	refTagOpts := &bitbucket.RepositoryRefOptions{
		Owner:    owner,
		RepoSlug: repo,
	}

	resTagRefs, err := c.Repositories.Repository.ListRefs(refTagOpts)
	if err != nil {
		t.Error("The ref/tag is not found.")
	}

	tagExpected := struct {
		n string
		t string
	}{}

	for _, ref := range resTagRefs.Refs {
		for k, v := range ref {
			// kCopy := k
			vCopy := v
			if val, ok := vCopy.(string); ok {
				if k == "name" && val == "TestGetRepoRefsTag" {
					tagExpected.n = val
				}
			}
			if val, ok := vCopy.(string); ok {
				if k == "type" && val == "tag" {
					tagExpected.t = val
				}
			}
		}
	}

	if !(tagExpected.n == "TestGetRepoRefsTag" && tagExpected.t == "tag") {
		t.Error("Could not list refs/tag that was created in test setup")
	}

	delTagOpt := &bitbucket.RepositoryTagOptions{
		Owner:    owner,
		RepoSlug: repo,
		Name:     "TestGetRepoRefsTag",
	}
	resTagDel, err := c.Repositories.Repository.DeleteTag(delTagOpt)
	if err != nil {
		t.Error("Could not delete tag")
	}

	if mp, ok := resTagDel.(map[string]interface{}); ok {
		for k, v := range mp {
			if k == "type" && v == "error" {
				t.Error("Delete tag returned an error, when it should have successfully" +
					"delete the test tag created during test setup")
			}
		}
	}
}
