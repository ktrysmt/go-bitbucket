package tests

import (
	"reflect"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
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

func TestClientNewBasicAuthWithCaCert(t *testing.T) {

	caCerts, err := FetchCACerts("api.bitbucket.org", "443")
	if err != nil {
		t.Error(err)
	}

	c, err := bitbucket.NewBasicAuthWithCaCert("example", "password", caCerts)
	if err != nil {
		t.Fatal(err)
	}

	r := reflect.ValueOf(c)

	if r.Type().String() != "*bitbucket.Client" {
		t.Error("Unknown error by `NewBasicAuthWithCaCert`.")
	}
}

func TestClientWithBearerToken(t *testing.T) {

	c, err := bitbucket.NewOAuthbearerToken("token")
	if err != nil {
		t.Fatal(err)
	}

	r := reflect.ValueOf(c)
	if r.Type().String() != "*bitbucket.Client" {
		t.Error("Unknown error by `NewOAuthbearerToken`.")
	}
}

func TestClientWithBearerTokenWithCaCert(t *testing.T) {

	caCerts, err := FetchCACerts("api.bitbucket.org", "443")
	if err != nil {
		t.Error(err)
	}

	c, err := bitbucket.NewOAuthbearerTokenWithCaCert("token", caCerts)
	if err != nil {
		t.Fatal(err)
	}

	r := reflect.ValueOf(c)
	if r.Type().String() != "*bitbucket.Client" {
		t.Error("Unknown error by `NewOAuthbearerTokenWithCaCert`.")
	}
}
