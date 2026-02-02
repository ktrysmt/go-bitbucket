package tests

import (
	"os"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/stretchr/testify/assert"
)

/*
Note: Run test with test_utils.go for util function dependencies
go test tests/repository_access_token_test.go tests/test_utils.go
*/

func TestAddGetandDeletePipelineVariableAccess(t *testing.T) {

	accessToken := os.Getenv("BITBUCKET_TEST_ACCESS_TOKEN")
	workspace := os.Getenv("BITBUCKET_TEST_OWNER")
	repo := os.Getenv("BITBUCKET_TEST_REPOSLUG")

	if accessToken == "" {
		t.Errorf("BITBUCKET_TEST_ACCESS_TOKEN is unset")
	}
	if workspace == "" {
		t.Errorf("BITBUCKET_TEST_OWNER is unset")
	}
	if repo == "" {
		t.Errorf("BITBUCKET_TEST_REPOSLUG is unset")
	}

	c, err := bitbucket.NewOAuthbearerToken(accessToken)
	if err != nil {
		t.Error(err)
	}

	variable := &bitbucket.RepositoryPipelineVariableOptions{
		Owner:    workspace,
		RepoSlug: repo,
		Key:      "test_key_to_delete",
		Value:    "Some value to delete",
		Secured:  false,
	}

	res, err := c.Repositories.Repository.AddPipelineVariable(variable)
	if err != nil {
		t.Error(err)
	}

	opt := &bitbucket.RepositoryPipelineVariableOptions{
		Owner:    workspace,
		RepoSlug: repo,
		Uuid:     res.Uuid,
	}

	optd := &bitbucket.RepositoryPipelineVariableDeleteOptions{
		Owner:    workspace,
		RepoSlug: repo,
		Uuid:     res.Uuid,
	}

	res, err = c.Repositories.Repository.GetPipelineVariable(opt)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, opt.Uuid, res.Uuid)

	// On success the delete API doesn't return any content (HTTP status 204)
	_, err = c.Repositories.Repository.DeletePipelineVariable(optd)
	if err != nil {
		t.Error(err)
	}

}

func TestAddGetandDeletePipelineVariableAccessTokenCaCert(t *testing.T) {

	accessToken := os.Getenv("BITBUCKET_TEST_ACCESS_TOKEN")
	workspace := os.Getenv("BITBUCKET_TEST_OWNER")
	repo := os.Getenv("BITBUCKET_TEST_REPOSLUG")

	if accessToken == "" {
		t.Errorf("BITBUCKET_TEST_ACCESS_TOKEN is unset")
	}
	if workspace == "" {
		t.Errorf("BITBUCKET_TEST_OWNER is unset")
	}
	if repo == "" {
		t.Errorf("BITBUCKET_TEST_REPOSLUG is unset")
	}

	caCert, err := fetchCACerts("api.bitbucket.org", "443")
	if err != nil {
		t.Error(err)
	}
	c, err := bitbucket.NewOAuthbearerTokenWithCaCert(accessToken, caCert)
	if err != nil {
		t.Error(err)
	}

	variable := &bitbucket.RepositoryPipelineVariableOptions{
		Owner:    workspace,
		RepoSlug: repo,
		Key:      "test_key_to_delete",
		Value:    "Some value to delete",
		Secured:  false,
	}

	res, err := c.Repositories.Repository.AddPipelineVariable(variable)
	if err != nil {
		t.Error(err)
	}

	opt := &bitbucket.RepositoryPipelineVariableOptions{
		Owner:    workspace,
		RepoSlug: repo,
		Uuid:     res.Uuid,
	}

	optd := &bitbucket.RepositoryPipelineVariableDeleteOptions{
		Owner:    workspace,
		RepoSlug: repo,
		Uuid:     res.Uuid,
	}

	res, err = c.Repositories.Repository.GetPipelineVariable(opt)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, opt.Uuid, res.Uuid)

	// On success the delete API doesn't return any content (HTTP status 204)
	_, err = c.Repositories.Repository.DeletePipelineVariable(optd)
	if err != nil {
		t.Error(err)
	}

}
