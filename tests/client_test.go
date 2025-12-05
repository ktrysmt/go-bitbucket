package tests

import (
	"github.com/ktrysmt/go-bitbucket"
	"reflect"
	"testing"
)

func TestClientNewBasicAuth(t *testing.T) {

	c, err := bitbucket.NewBasicAuth("example", "password")
	if err != nil {
		t.Fatal(err)
	}

	r := reflect.ValueOf(c)

	if r.Type().String() != "*bitbucket.Client" {
		t.Error("Unknown error by `NewBasicAuth`.")
	}
}
