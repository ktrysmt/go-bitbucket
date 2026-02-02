package tests

import (
	"reflect"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
)

const DUMMY_CA_CERT = "-----BEGIN CERTIFICATE-----IxMDM1MV0ZDJkZjM...-----END CERTIFICATE-----"

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

	c, err := bitbucket.NewBasicAuthWithCaCert("example", "password", []byte(DUMMY_CA_CERT))
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

	c, err := bitbucket.NewOAuthbearerTokenWithCaCert("token", []byte(DUMMY_CA_CERT))
	if err != nil {
		t.Fatal(err)
	}

	r := reflect.ValueOf(c)
	if r.Type().String() != "*bitbucket.Client" {
		t.Error("Unknown error by `NewOAuthbearerTokenWithCaCert`.")
	}
}
