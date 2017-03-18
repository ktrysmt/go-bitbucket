package tests

import (
	"github.com/ktrysmt/go-bitbucket"
	"os"
	"testing"
)

var (
	user  = os.Getenv("BITBUCKET_TEST_USERNAME")
	pass  = os.Getenv("BITBUCKET_TEST_PASSWORD")
	owner = os.Getenv("BITBUCKET_TEST_OWNER")
	repo  = os.Getenv("BITBUCKET_TEST_REPOSLUG")
)

func setup(t *testing.T) *bitbucket.Client {

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
	return c
}

func TestCreateBranchRestrictionsKindPush(t *testing.T) {

	c := setup(t)

	opt := &bitbucket.BranchRestrictionsOptions{
		Owner:     owner,
		Pattern:   "develop",
		Repo_slug: repo,
		Kind:      "push",
		Users:     []string{user},
	}
	res := c.Repositories.BranchRestrictions.Create(opt)
	jsonMap := res.(map[string]interface{})
	if jsonMap["type"] != "branchrestriction" {
		t.Error("is not branchrestriction type")
	}
	if jsonMap["kind"] != "push" {
		t.Error("did not match branchrestriction kind")
	}
}

func TestCreateBranchRestrictionsKindRequirePassingBuilds(t *testing.T) {

	c := setup(t)

	opt := &bitbucket.BranchRestrictionsOptions{
		Owner:     owner,
		Pattern:   "master",
		Repo_slug: repo,
		Kind:      "require_passing_builds_to_merge",
		Value:     2,
	}
	res := c.Repositories.BranchRestrictions.Create(opt)
	jsonMap := res.(map[string]interface{})
	if jsonMap["type"] != "branchrestriction" {
		t.Error("is not branchrestriction type")
	}
	if jsonMap["kind"] != "require_passing_builds_to_merge" {
		t.Error("did not match branchrestriction kind")
	}
}
