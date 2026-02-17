package tests

import (
	"testing"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/stretchr/testify/assert"
)

/*
Note: Run test with test_utils.go for util function dependencies
go test tests/repository_access_token_test.go tests/test_utils.go
*/

func TestAddGetandDeletePipelineVariableAccess(t *testing.T) {

	client, err := SetupBearerToken(t)
	if err != nil {
		t.Error(err)
		return
	}

	variable := &bitbucket.RepositoryPipelineVariableOptions{
		Owner:    ownerEnv,
		RepoSlug: repoEnv,
		Key:      "test_key_to_delete",
		Value:    "Some value to delete",
		Secured:  false,
	}

	testApiCalls(t, client, variable)
}

func TestAddGetandDeletePipelineVariableAccessWithTokenBaseUrlCaCert(t *testing.T) {

	client0, err := SetupBearerTokenWithBaseUrlStrCaCert(t, "", nil)
	if err != nil {
		t.Error(err)
		assert.Nil(t, err)
	}

	expectedBaseUrlStr := "https://api.bitbucket.org/2.0"
	client1, err := SetupBearerTokenWithBaseUrlStrCaCert(t, expectedBaseUrlStr, nil)
	if err != nil {
		t.Error(err)
		assert.Nil(t, err)
	}
	assert.Equal(t, client0.Auth, client1.Auth)

	variable := &bitbucket.RepositoryPipelineVariableOptions{
		Owner:    ownerEnv,
		RepoSlug: repoEnv,
		Key:      "test_key_to_delete",
		Value:    "Some value to delete",
		Secured:  false,
	}

	testApiCalls(t, client0, variable)
}

func testApiCalls(t *testing.T, c *bitbucket.Client, v *bitbucket.RepositoryPipelineVariableOptions) {
	res, err := c.Repositories.Repository.AddPipelineVariable(v)
	if err != nil {
		t.Error(err)
	}

	opt := &bitbucket.RepositoryPipelineVariableOptions{
		Owner:    v.Owner,
		RepoSlug: v.RepoSlug,
		Uuid:     res.Uuid,
	}

	optd := &bitbucket.RepositoryPipelineVariableDeleteOptions{
		Owner:    v.Owner,
		RepoSlug: v.RepoSlug,
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
		assert.Nilf(t, err, "expected no error returned, but got %v", err)
	}
}
