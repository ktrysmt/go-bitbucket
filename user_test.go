package bitbucket

import (
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
	"os"
	"testing"
)

func TestProfile(t *testing.T) {

	user := os.Getenv("BITBUCKET_TEST_USERNAME")
	pass := os.Getenv("BITBUCKET_TEST_PASSWORD")

	if user == "" {
		t.Error("username is empty.")
	}

	if pass == "" {
		t.Error("password is empty.")
	}

	c := bitbucket.NewBasicAuth(user, pass)

	res := c.User.Profile()

	fmt.Println(res) // receive the data as json format
}
