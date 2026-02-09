package mock_tests

import (
	"errors"
	"testing"

	go_bitbucket "github.com/ktrysmt/go-bitbucket"
	"github.com/ktrysmt/go-bitbucket/mockgen"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMockPullRequests_Gets_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPullRequestInst := mockgen.NewMockpullrequests(ctrl)

	inPullRequestOpts := go_bitbucket.PullRequestsOptions{
		Owner:    "test-workspace",
		RepoSlug: "test-repo",
	}

	outPullRequest := go_bitbucket.PullRequestsOptions{
		ID:       "test-pull-request-id",
		Owner:    "test-workspace",
		RepoSlug: "test-repo",
		Title:    "test-pull-request",
		Commit:   "test-commit",
	}

	var expectedPullRequestList []go_bitbucket.PullRequestsOptions
	expectedPullRequestList = append(expectedPullRequestList, outPullRequest)

	mockPullRequestInst.EXPECT().
		List(inPullRequestOpts).
		Times(1).
		Return(expectedPullRequestList, nil)

	actualPullRequestList, actualErr := mockPullRequestInst.List(inPullRequestOpts)

	assert.Nil(t, actualErr, "No error should have been returned, but got: %v", actualErr)
	assert.NotNil(t, actualPullRequestList, "Pull requests should be returned.")
	assert.Equal(t, expectedPullRequestList, actualPullRequestList, "Actual and Expected pull request lists should be equal.")
}

func TestMockPullRequests_Gets_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPullRequestInst := mockgen.NewMockpullrequests(ctrl)
	expectedMockError := errors.New("Not Found")

	inPullRequestOps := go_bitbucket.PullRequestsOptions{
		Owner:    "test-workspace",
		RepoSlug: "test-repo",
	}

	mockPullRequestInst.EXPECT().
		List(inPullRequestOps).
		Times(1).
		Return(nil, expectedMockError)

	actualPullRequestList, actualErr := mockPullRequestInst.List(inPullRequestOps)

	assert.NotNil(t, actualErr)
	assert.Nil(t, actualPullRequestList, "The returned list of pull requests should be nil, but got: %v", actualPullRequestList)
}
