package tests

import (
	"os"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/ktrysmt/go-bitbucket"
)

func TestDiff(t *testing.T) {

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

	c := bitbucket.NewBasicAuth(user, pass)

	spec := "master..develop"

	opt := &bitbucket.DiffOptions{
		Owner:    owner,
		RepoSlug: repo,
		Spec:     spec,
	}
	res, _ := c.Repositories.Diff.GetDiff(opt)

	pp.Println(res)

	if res == nil {
		t.Error("It could not get the raw response.")
	}
}
