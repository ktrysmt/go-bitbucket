package tests

import (
	"reflect"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
)

const (
	EXPECTED_CLIENT_TYPE_STR = "*bitbucket.Client"
)

/*
These are critical tests at the client generation stage that will cause upstream failures
if not addressed early, so a fatal errors are expected.
*/

func TestClientNewBasicAuth(t *testing.T) {

	c, err := bitbucket.NewBasicAuth("example", "password")
	if err != nil {
		t.Fatal("Error returned from `NewBasicAuth` got: ", err)
	}

	r := reflect.ValueOf(c)
	actualClientTypeStr := r.Type().String()
	if actualClientTypeStr != EXPECTED_CLIENT_TYPE_STR {
		t.Fatalf("Incorrect client type generated, expected: %s, got: %s", EXPECTED_CLIENT_TYPE_STR, actualClientTypeStr)
	}
}

func TestClientNewBasicAuthWithCaCert(t *testing.T) {

	caCerts, err := FetchCACerts("api.bitbucket.org", "443")
	if err != nil {
		t.Fatal("Error returned from `FetchCACerts` got: ", err)
	}

	c, err := bitbucket.NewBasicAuthWithCaCert("example", "password", caCerts)
	if err != nil {
		t.Fatal("Error returned from `NewBasicAuthWithCaCert got: ", err)
	}

	r := reflect.ValueOf(c)
	actualClientTypeStr := r.Type().String()
	if actualClientTypeStr != "*bitbucket.Client" {
		t.Fatalf("Incorrect client type generated, expected: %s, got: %s", EXPECTED_CLIENT_TYPE_STR, actualClientTypeStr)
	}
}

func TestClientWithBearerToken(t *testing.T) {

	c, err := bitbucket.NewOAuthbearerToken("token")
	if err != nil {
		t.Fatal("Error returned from `NewOAuthbearerToken` got: ", err)
	}

	r := reflect.ValueOf(c)
	actualClientTypeStr := r.Type().String()
	if actualClientTypeStr != "*bitbucket.Client" {
		t.Fatalf("Incorrect client type generated, expected: %s, got: %s", EXPECTED_CLIENT_TYPE_STR, actualClientTypeStr)
	}
}

func TestClientWithBearerTokenWithCaCert(t *testing.T) {

	caCerts, err := FetchCACerts("api.bitbucket.org", "443")
	if err != nil {
		t.Fatal("Error returned from `FetchCACerts` got: ", err)
	}

	c, err := bitbucket.NewOAuthbearerTokenWithCaCert("token", caCerts)
	if err != nil {
		t.Fatal("Error returned from `NewOAuthbearerTokenWithCaCert` got: ", err)
	}
	r := reflect.ValueOf(c)
	actualClientTypeStr := r.Type().String()
	if actualClientTypeStr != EXPECTED_CLIENT_TYPE_STR {
		t.Fatalf("Incorrect client type generated, expected: %s, got: %s", EXPECTED_CLIENT_TYPE_STR, actualClientTypeStr)
	}
}
