package tests

import (
	"os"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
)

func TestProfile(t *testing.T) {

	user := os.Getenv("BITBUCKET_TEST_USERNAME")
	pass := os.Getenv("BITBUCKET_TEST_PASSWORD")

	if user == "" {
		t.Error("BITBUCKET_TEST_USERNAME is empty.")
	}

	if pass == "" {
		t.Error("BITBUCKET_TEST_PASSWORD is empty.")
	}

	c := bitbucket.NewBasicAuth(user, pass)

	res, _ := c.User.Profile()

	jsonMap := res.(map[string]interface{})

	if jsonMap["username"] != user {
		t.Error("Cannot catch the Profile.username.")
	}
}
