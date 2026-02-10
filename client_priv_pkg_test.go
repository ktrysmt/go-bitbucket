package bitbucket

import (
	"testing"

	"github.com/ktrysmt/go-bitbucket/tests"
	"github.com/stretchr/testify/assert"
)

func TestAppendCaCerts_util_test(t *testing.T) {
	caCerts, err := tests.FetchCACerts("bitbucket.org", "443")
	if err != nil {
		t.Fatalf("Error fetching CA certs using `FetchCACerts`: %v", err)
	}
	httpClient, err := appendCaCerts(caCerts)
	if err != nil {
		t.Fatalf("Error returned from `appendCaCerts` failed to create the http client: %v", err)
	}
	assert.NotNil(t, httpClient)
}
