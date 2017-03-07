package tests

import (
	"github.com/k0kubun/pp"
	"github.com/ktrysmt/go-bitbucket"
	"os"
	"testing"
)

func TestGetsBranchRestrictions(t *testing.T) {

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

	opt := &bitbucket.BranchRestrictionsOptions{
		Owner:     owner,
		Pattern:   "master",
		Repo_slug: repo,
		Kind:      "push",
		Users:     []string{"kotaro_yoshimatsu"},
	}

	res := c.Repositories.BranchRestrictions.Create(opt)
	jsonMap := res.(map[string]interface{})
	pp.Println(jsonMap)
	// if res.Full_name != owner+"/"+repo {
	// 	t.Error("Cannot catch repos full name.")
	// }
}
