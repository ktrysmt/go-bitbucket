package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raphaeldevs/go-bitbucket"
)

func TestProfile(t *testing.T) {

	user := getUsername()
	pass := getPassword()

	c := bitbucket.NewBasicAuth(user, pass)

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
