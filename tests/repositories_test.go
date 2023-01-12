package tests

import (
	"os"
	"testing"

	"github.com/elvenworks/go-bitbucket"
)

func TestListForAccount(t *testing.T) {
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

	repositories, err := c.Repositories.ListForAccount(&bitbucket.RepositoriesOptions{
		Owner: owner,
	})
	if err != nil {
		t.Error("Unable to fetch repositories")
	}

	found := false
	for _, r := range repositories.Items {
		if r.Slug == repo {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find repository in list")
	}
}

func TestListForTeam(t *testing.T) {
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

	//goland:noinspection GoDeprecation
	repositories, err := c.Repositories.ListForTeam(&bitbucket.RepositoriesOptions{

		Owner: owner,
	})
	if err != nil {
		t.Error("Unable to fetch repositories")
	}

	found := false
	for _, r := range repositories.Items {
		if r.Slug == repo {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find repository in list")
	}
}
