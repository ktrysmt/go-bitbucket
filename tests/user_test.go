package tests

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
)

func TestProfile(t *testing.T) {

	user := getUsername()
	pass := getPassword()

	c, err := bitbucket.NewBasicAuth(user, pass)
	if err != nil {
		t.Fatal(err)
	}

	res, err := c.User.Profile()

	assert.NoError(t, err)
	assert.NotNil(t, res)
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
