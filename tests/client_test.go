package tests

import (
	"reflect"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
)

const (
	EXPECTED_CLIENT_TYPE_STR = "*bitbucket.Client"
	EXPECTED_BASE_URL_STR    = "https://api.bitbucket.org/2.0"
	EXPECTED_BASE_URL_HOST   = "api.bitbucket.org"
	EXPECTED_BASE_URL_PORT   = "443"
	EXPECTED_TOKEN           = "token"
	EXPECTED_USERNAME        = "example"
	EXPECTED_PASSWORD        = "password"
)

/*
These are critical tests at the client generation stage that will cause upstream failures
if not addressed early, so a fatal errors are expected.
*/

func TestClientNewBasicAuth(t *testing.T) {

	c, err := bitbucket.NewBasicAuth(EXPECTED_USERNAME, EXPECTED_PASSWORD)
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

	caCerts, err := FetchCACerts(EXPECTED_BASE_URL_HOST, EXPECTED_BASE_URL_PORT)
	if err != nil {
		t.Fatal("Error returned from `FetchCACerts` got: ", err)
	}

	c, err := bitbucket.NewBasicAuthWithCaCert(EXPECTED_USERNAME, EXPECTED_PASSWORD, caCerts)
	if err != nil {
		t.Fatal("Error returned from `NewBasicAuthWithCaCert got: ", err)
	}

	r := reflect.ValueOf(c)
	actualClientTypeStr := r.Type().String()
	if actualClientTypeStr != EXPECTED_CLIENT_TYPE_STR {
		t.Fatalf("Incorrect client type generated, expected: %s, got: %s", EXPECTED_CLIENT_TYPE_STR, actualClientTypeStr)
	}
}

func TestClientNewBasicAuthBaseUrlStr(t *testing.T) {

	c, err := bitbucket.NewBasicAuthWithBaseUrlStr(EXPECTED_USERNAME, EXPECTED_PASSWORD, EXPECTED_BASE_URL_STR)
	if err != nil {
		t.Fatal("Error returned from `NewBasicAuthWithBaseUrlStr` got: ", err)
	}
	actualBaseUrlStr := c.GetApiBaseURL()
	if actualBaseUrlStr != EXPECTED_BASE_URL_STR {
		t.Fatalf("Incorrect base url generated, expected: %s, got: %s", EXPECTED_BASE_URL_STR, actualBaseUrlStr)
	}

	r := reflect.ValueOf(c)
	actualClientTypeStr := r.Type().String()
	if actualClientTypeStr != EXPECTED_CLIENT_TYPE_STR {
		t.Fatalf("Incorrect client type generated, expected: %s, got: %s", EXPECTED_CLIENT_TYPE_STR, actualClientTypeStr)
	}
}

func TestClientNewBasicAuthBaseUrlStrCaCert(t *testing.T) {

	caCerts, err := FetchCACerts(EXPECTED_BASE_URL_HOST, EXPECTED_BASE_URL_PORT)
	if err != nil {
		t.Fatal("Error returned from `FetchCACerts` got: ", err)
	}

	c, err := bitbucket.NewBasicAuthWithBaseUrlStrCaCert(EXPECTED_USERNAME, EXPECTED_PASSWORD, EXPECTED_BASE_URL_STR, caCerts)
	if err != nil {
		t.Fatal("Error returned from `NewBasicAuthWithBaseUrlStrCaCert` got: ", err)
	}
	actualBaseUrlStr := c.GetApiBaseURL()
	if actualBaseUrlStr != EXPECTED_BASE_URL_STR {
		t.Fatalf("Incorrect base url generated, expected: %s, got: %s", EXPECTED_BASE_URL_STR, actualBaseUrlStr)
	}

	r := reflect.ValueOf(c)
	actualClientTypeStr := r.Type().String()
	if actualClientTypeStr != EXPECTED_CLIENT_TYPE_STR {
		t.Fatalf("Incorrect client type generated, expected: %s, got: %s", EXPECTED_CLIENT_TYPE_STR, actualClientTypeStr)
	}
}

func TestClientWithBearerToken(t *testing.T) {

	c, err := bitbucket.NewOAuthbearerToken(EXPECTED_TOKEN)
	if err != nil {
		t.Fatal("Error returned from `NewOAuthbearerToken` got: ", err)
	}

	r := reflect.ValueOf(c)
	actualClientTypeStr := r.Type().String()
	if actualClientTypeStr != EXPECTED_CLIENT_TYPE_STR {
		t.Fatalf("Incorrect client type generated, expected: %s, got: %s", EXPECTED_CLIENT_TYPE_STR, actualClientTypeStr)
	}
}

func TestClientWithBearerTokenWithCaCert(t *testing.T) {

	caCerts, err := FetchCACerts(EXPECTED_BASE_URL_HOST, EXPECTED_BASE_URL_PORT)
	if err != nil {
		t.Fatal("Error returned from `FetchCACerts` got: ", err)
	}

	c, err := bitbucket.NewOAuthbearerTokenWithCaCert(EXPECTED_TOKEN, caCerts)
	if err != nil {
		t.Fatal("Error returned from `NewOAuthbearerTokenWithCaCert` got: ", err)
	}
	r := reflect.ValueOf(c)
	actualClientTypeStr := r.Type().String()
	if actualClientTypeStr != EXPECTED_CLIENT_TYPE_STR {
		t.Fatalf("Incorrect client type generated, expected: %s, got: %s", EXPECTED_CLIENT_TYPE_STR, actualClientTypeStr)
	}
}

func TestClientWithBearerTokenWithBaseUrlStrCaCert(t *testing.T) {

	caCerts, err := FetchCACerts(EXPECTED_BASE_URL_HOST, EXPECTED_BASE_URL_PORT)
	if err != nil {
		t.Fatal("Error returned from `FetchCACerts` got: ", err)
	}

	c, err := bitbucket.NewOAuthbearerTokenWithBaseUrlStrCaCert(EXPECTED_TOKEN, EXPECTED_BASE_URL_STR, caCerts)
	if err != nil {
		t.Fatal("Error returned from `NewOAuthbearerTokenWithBaseUrlStrCaCert` got: ", err)
	}

	actualBaseUrlStr := c.GetApiBaseURL()
	if actualBaseUrlStr != EXPECTED_BASE_URL_STR {
		t.Fatalf("Incorrect base url generated, expected: %s, got: %s", EXPECTED_BASE_URL_STR, actualBaseUrlStr)
	}

	r := reflect.ValueOf(c)
	actualClientTypeStr := r.Type().String()
	if actualClientTypeStr != EXPECTED_CLIENT_TYPE_STR {
		t.Fatalf("Incorrect client type generated, expected: %s, got: %s", EXPECTED_CLIENT_TYPE_STR, actualClientTypeStr)
	}
}
