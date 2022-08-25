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

func TestGetRepositoryPipelineConfig(t *testing.T) {
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

	opt := &bitbucket.RepositoryPipelineOptions{
		Owner:    owner,
		RepoSlug: repo,
	}

	res, err := c.Repositories.Repository.GetPipelineConfig(opt)
	if err != nil {
		t.Error(err)
	}

	if res == nil {
		t.Error("Cannot get pipeline config")
	}
	if res.Enabled != false {
		t.Error("Got wrong pipelines config data")
	}
}

func TestUpdateRepositoryPipelineConfig(t *testing.T) {
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

	opt := &bitbucket.RepositoryPipelineOptions{
		Owner:    owner,
		RepoSlug: repo,
		Enabled:  true,
	}

	res, err := c.Repositories.Repository.UpdatePipelineConfig(opt)
	if err != nil {
		t.Error(err)
	}

	if res == nil {
		t.Error("Cannot update pipeline config")
	}
	if res.Enabled != true {
		t.Error("Got wrong pipelines config data")
	}

	opt = &bitbucket.RepositoryPipelineOptions{
		Owner:    owner,
		RepoSlug: repo,
		Enabled:  false,
	}

	res, err = c.Repositories.Repository.UpdatePipelineConfig(opt)
	if err != nil {
		t.Error(err)
	}

	if res == nil {
		t.Error("Cannot update pipeline config")
	}
	if res.Enabled != false {
		t.Error("Got wrong pipelines config data")
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

func TestListRepositoryGroupPermissions(t *testing.T) {

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

	res, err := c.Repositories.Repository.ListGroupPermissions(opt)

	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("Cannot list repository group permissions")
	}

}

func TestSetRepositoryGroupPermissions(t *testing.T) {

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

	opt := &bitbucket.RepositoryGroupPermissionsOptions{
		Owner:      owner,
		RepoSlug:   repo,
		Group:      "developers",
		Permission: "read",
	}

	res, err := c.Repositories.Repository.SetGroupPermissions(opt)
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("Cannot set repository group permissions")
	}

}

func TestDeleteRepositoryGroupPermissions(t *testing.T) {

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

	opt := &bitbucket.RepositoryGroupPermissionsOptions{
		Owner:    owner,
		RepoSlug: repo,
		Group:    "developers",
	}

	_, err := c.Repositories.Repository.DeleteGroupPermissions(opt)
	if err != nil {
		t.Error(err)
	}
}

func TestGetRepositoryGroupPermissions(t *testing.T) {

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

	opt := &bitbucket.RepositoryGroupPermissionsOptions{
		Owner:    owner,
		RepoSlug: repo,
		Group:    "developers",
	}

	res, err := c.Repositories.Repository.GetGroupPermissions(opt)
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("Cannot get repository group permissions")
	}
}

func TestListRepositoryUserPermissions(t *testing.T) {

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

	res, err := c.Repositories.Repository.ListUserPermissions(opt)

	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("Cannot list repository user permissions")
	}

}
func TestGetRepositoryUserPermissions(t *testing.T) {

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

	opt := &bitbucket.RepositoryUserPermissionsOptions{
		Owner:    owner,
		RepoSlug: repo,
		User:     "5db9ab0fa766000da47cd9e4",
	}

	res, err := c.Repositories.Repository.GetUserPermissions(opt)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
	if res == nil {
		t.Error("Cannot get repository user permissions")
	}
}

func TestDeleteRepositoryUserPermissions(t *testing.T) {

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

	opt := &bitbucket.RepositoryUserPermissionsOptions{
		Owner:    owner,
		RepoSlug: repo,
		User:     "5db9ab0fa766000da47cd9e4",
	}

	_, err := c.Repositories.Repository.DeleteUserPermissions(opt)
	if err != nil {
		t.Error(err)
	}
}

func TestSetRepositoryUserPermissions(t *testing.T) {

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

	opt := &bitbucket.RepositoryUserPermissionsOptions{
		Owner:      owner,
		RepoSlug:   repo,
		User:       "5e78617ce6f50f0c3a00e583",
		Permission: "admin",
	}

	res, err := c.Repositories.Repository.SetUserPermissions(opt)
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Error("Cannot set repository user permissions")
	}

}
