package bitbucket

import (
	"testing"
)

func TestUsername(t *testing.T) {

	c := New("example", "password")

	if c.app_id != "example" {
		t.Error("username not equal")
	}

	if c.secret != "password" {
		t.Error("password not equal")
	}
}
