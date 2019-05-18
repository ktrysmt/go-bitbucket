package tests

import (
	"github.com/davidji99/go-bitbucket"
	"reflect"
	"testing"
)

func TestClientNewBasicAuth(t *testing.T) {

	c := bitbucket.NewBasicAuth("example", "password")

	r := reflect.ValueOf(c)

	if r.Type().String() != "*bitbucket.Client" {
		t.Error("Unknown error by `NewBasicAuth`.")
	}
}
