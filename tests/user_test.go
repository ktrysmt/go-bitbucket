package tests

import (
	"os"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
)

func TestProfile(t *testing.T) {

	user := getUsername()
	pass := getPassword()

	if user == "" {
		t.Error("BITBUCKET_TEST_USERNAME is empty.")
	}

	if pass == "" {
		t.Error("BITBUCKET_TEST_PASSWORD is empty.")
	}

	c := bitbucket.NewBasicAuth(user, pass)

	res, err := c.User.Profile()
	if err != nil {
		t.Error(err)
	}

	if res.Username != user {
		t.Error("Cannot catch the Profile.username.")
	}
}

func getUsername() string {
	ev := os.Getenv("BITBUCKET_TEST_USERNAME")
	if ev != "" {
		return ev
	}

	return "example-username"
}

func getPassword() string {
	ev := os.Getenv("BITBUCKET_TEST_PASSWORD")
	if ev != "" {
		return ev
	}

	return "password"
}
