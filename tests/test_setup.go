package tests

import (
	"os"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
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

	c, err := bitbucket.NewBasicAuth(user, pass)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
