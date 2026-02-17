package tests

import (
	"errors"
	"net/url"
	"os"
	"testing"

	"github.com/ktrysmt/go-bitbucket"
)

var (
	userEnv        = os.Getenv("BITBUCKET_TEST_USERNAME")
	passEnv        = os.Getenv("BITBUCKET_TEST_PASSWORD")
	ownerEnv       = os.Getenv("BITBUCKET_TEST_OWNER")
	repoEnv        = os.Getenv("BITBUCKET_TEST_REPOSLUG")
	accessTokenEnv = os.Getenv("BITBUCKET_TEST_ACCESS_TOKEN")
	baseUrlStrEnv  = os.Getenv("BITBUCKET_API_BASE_URL")
)

func checkOwnerRepoSet(t *testing.T) error {
	if ownerEnv == "" {
		err := errors.New("BITBUCKET_TEST_OWNER is empty.")
		t.Error(err)
		return err
	}
	if repoEnv == "" {
		err := errors.New("BITBUCKET_TEST_REPOSLUG is empty.")
		t.Error(err)
		return err
	}
	return nil
}

func checkAccessTokenSet(t *testing.T) error {
	if accessTokenEnv == "" {
		err := errors.New("BITBUCKET_TEST_ACCESS_TOKEN is empty.")
		t.Error(err)
		return err
	}
	return nil
}

func checkBaseUrlStrSet(t *testing.T, baseUrlStr string) (string, error) {
	if baseUrlStr == "" {
		if baseUrlStrEnv == "" {
			err := errors.New("BITBUCKET_TEST_BASE_URL is empty.")
			t.Error(err)
			return "", err
		}
		baseUrlStr = baseUrlStrEnv
	}
	return baseUrlStr, nil
}

func setupBasicAuthTest(t *testing.T) (*bitbucket.Client, error) {

	if userEnv == "" {
		err := errors.New("BITBUCKET_TEST_USERNAME is empty.")
		t.Error(err)
		return nil, err
	}
	if passEnv == "" {
		err := errors.New("BITBUCKET_TEST_PASSWORD is empty.")
		t.Error(err)
		return nil, err
	}
	err := checkOwnerRepoSet(t)
	if err != nil {
		t.Error(err)
		return nil, err
	}

	c, err := bitbucket.NewBasicAuth(userEnv, passEnv)
	if err != nil {
		t.Error(err)
		return nil, err
	}
	return c, nil
}

func SetupBearerToken(t *testing.T) (*bitbucket.Client, error) {

	err := checkOwnerRepoSet(t)
	if err != nil {
		t.Error(err)
		return nil, err
	}

	c, err := bitbucket.NewOAuthbearerToken(accessTokenEnv)
	if err != nil {
		t.Error(err)
		return nil, err
	}
	return c, nil
}

func SetupBearerTokenWithBaseUrlStr(t *testing.T, baseUrlStr string) (*bitbucket.Client, error) {
	err := checkAccessTokenSet(t)
	if err != nil {
		t.Error(err)
		return nil, err
	}

	baseUrlStr, err = checkBaseUrlStrSet(t, baseUrlStr)
	if err != nil {
		t.Error(err)
		return nil, err
	}

	c, err := bitbucket.NewOAuthbearerTokenWithBaseUrlStr(accessTokenEnv, baseUrlStr)
	if err != nil {
		t.Error(err)
		return nil, err
	}
	return c, nil
}

func SetupBearerTokenWithBaseUrlStrCaCert(t *testing.T, baseUrlStr string, caCerts []byte) (*bitbucket.Client, error) {
	err := checkAccessTokenSet(t)
	if err != nil {
		t.Error(err)
		return nil, err
	}

	baseUrlStr, err = checkBaseUrlStrSet(t, baseUrlStr)
	if err != nil {
		t.Error(err)
		return nil, err
	}

	if caCerts == nil {
		parsedUrl, err := url.Parse(baseUrlStr)
		if err != nil {
			t.Error(err)
			return nil, err
		}
		parsedPort := parsedUrl.Port()
		// url.Port does not consider default ports for http (80) and https (443)
		if parsedPort == "" {
			parsedPort = "443"
		}
		extractedCaCerts, err := FetchCACerts(parsedUrl.Hostname(), parsedPort)
		if err != nil {
			t.Error(err)
			return nil, err
		}
		caCerts = extractedCaCerts
	}

	c, err := bitbucket.NewOAuthbearerTokenWithBaseUrlStrCaCert(accessTokenEnv, baseUrlStr, caCerts)
	if err != nil {
		t.Error(err)
		return nil, err
	}
	return c, nil
}
