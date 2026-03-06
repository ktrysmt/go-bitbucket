package bitbucket

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateSelfSignedCert(t *testing.T) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{Organization: []string{"Test"}},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		IsCA:         true,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
}

func TestRequestUrl_WithArgs(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	result := client.requestUrl("/repositories/%s/%s", "owner", "repo")

	assert.Contains(t, result, "/2.0/repositories/owner/repo")
}

func TestRequestUrl_WithEmptyArg(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	result := client.requestUrl("/users/%s/", "")

	assert.Contains(t, result, "/2.0/users/")
	assert.NotContains(t, result, "%!(EXTRA")
}

func TestRequestUrl_NoArgs(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	result := client.requestUrl("/workspaces")

	assert.Contains(t, result, "/2.0/workspaces")
}

func TestRequestUrl_WithSpecialCharacters(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	result := client.requestUrl("/repositories/%s/%s/refs/branches/%s", "owner", "my-repo", "feature/special~chars")
	assert.Contains(t, result, "/2.0/repositories/owner/my-repo/refs/branches/feature/special~chars")

	resultWithSpaces := client.requestUrl("/repositories/%s/%s", "my owner", "my repo")
	assert.Contains(t, resultWithSpaces, "my owner")
	assert.Contains(t, resultWithSpaces, "my repo")

	resultWithUnicode := client.requestUrl("/repositories/%s/%s", "owner", "repo-\u00e9")
	assert.Contains(t, resultWithUnicode, "repo-\u00e9")
}

func TestAuthenticateRequest_BasicAuth(t *testing.T) {
	t.Parallel()
	var receivedUser, receivedPass string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedUser, receivedPass, _ = r.BasicAuth()
		respondJSON(w, http.StatusOK, map[string]interface{}{"ok": true})
	})
	defer server.Close()

	urlStr := client.requestUrl("/user")
	_, err := client.execute("GET", urlStr, "")

	require.NoError(t, err)
	assert.Equal(t, "test-user", receivedUser)
	assert.Equal(t, "test-pass", receivedPass)
}

func TestAuthenticateRequest_BearerToken(t *testing.T) {
	t.Parallel()
	var receivedAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		respondJSON(w, http.StatusOK, map[string]interface{}{"ok": true})
	}))
	defer server.Close()

	client, _ := NewOAuthbearerToken("my-token")
	serverURL, _ := url.Parse(server.URL + "/2.0")
	client.SetApiBaseURL(*serverURL)
	client.HttpClient = server.Client()

	urlStr := client.requestUrl("/user")
	_, err := client.execute("GET", urlStr, "")

	require.NoError(t, err)
	assert.Equal(t, "Bearer my-token", receivedAuth)
}

func TestExecute_Success(t *testing.T) {
	t.Parallel()
	expected := map[string]interface{}{
		"username": "testuser",
		"type":     "user",
	}

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, expected)
	})
	defer server.Close()

	urlStr := client.requestUrl("/user")
	result, err := client.execute("GET", urlStr, "")

	require.NoError(t, err)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "result should be a map[string]interface{}")
	assert.Len(t, resultMap, 2)
	assert.Contains(t, resultMap, "username")
	assert.Contains(t, resultMap, "type")
	assert.Equal(t, "testuser", resultMap["username"])
	assert.Equal(t, "user", resultMap["type"])
}

func TestExecute_PostWithBody(t *testing.T) {
	t.Parallel()
	var receivedMethod string
	var receivedContentType string
	var receivedBody string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedContentType = r.Header.Get("Content-Type")
		bodyBytes, _ := io.ReadAll(r.Body)
		receivedBody = string(bodyBytes)
		respondJSON(w, http.StatusCreated, map[string]interface{}{"id": "123"})
	})
	defer server.Close()

	body := `{"title":"test"}`
	urlStr := client.requestUrl("/repositories/%s/%s/issues", "owner", "repo")
	_, err := client.execute("POST", urlStr, body)

	require.NoError(t, err)
	assert.Equal(t, "POST", receivedMethod)
	assert.Equal(t, "application/json", receivedContentType)
	assert.Equal(t, body, receivedBody)
}

func TestExecute_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": map[string]interface{}{"message": "not found"},
		})
	})
	defer server.Close()

	urlStr := client.requestUrl("/repositories/%s/%s", "owner", "nonexistent")
	_, err := client.execute("GET", urlStr, "")

	assert.Error(t, err)
	var unexpectedErr *UnexpectedResponseStatusError
	assert.ErrorAs(t, err, &unexpectedErr)
}

func TestExecute_NoContentResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	urlStr := client.requestUrl("/repositories/%s/%s", "owner", "repo")
	result, err := client.execute("DELETE", urlStr, "")

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestExecutePaginated_SinglePage(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"name": "item1"},
			map[string]interface{}{"name": "item2"},
		}))
	})
	defer server.Close()

	urlStr := client.requestUrl("/repositories/%s/%s/commits/main", "owner", "repo")
	result, err := client.executePaginated("GET", urlStr, "", nil)

	require.NoError(t, err)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 2)
}

func TestExecutePaginated_MultiplePages(t *testing.T) {
	t.Parallel()
	callCount := 0

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			resp := paginatedResponse([]interface{}{
				map[string]interface{}{"name": "item1"},
			})
			// Point "next" to same server for second page
			resp["next"] = "http://" + r.Host + r.URL.Path + "?page=2"
			respondJSON(w, http.StatusOK, resp)
		} else {
			respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
				map[string]interface{}{"name": "item2"},
			}))
		}
	})
	defer server.Close()

	urlStr := client.requestUrl("/repositories/%s/%s/commits/main", "owner", "repo")
	result, err := client.executePaginated("GET", urlStr, "", nil)

	require.NoError(t, err)
	resultMap := result.(map[string]interface{})
	values := resultMap["values"].([]interface{})
	assert.Len(t, values, 2)
	assert.Equal(t, 2, callCount)
	firstItem := values[0].(map[string]interface{})
	assert.Equal(t, "item1", firstItem["name"])
	secondItem := values[1].(map[string]interface{})
	assert.Equal(t, "item2", secondItem["name"])
}

func TestExecutePaginated_DisableAutoPaging(t *testing.T) {
	t.Parallel()
	callCount := 0

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := paginatedResponse([]interface{}{
			map[string]interface{}{"name": "item1"},
		})
		resp["next"] = "http://" + r.Host + r.URL.Path + "?page=2"
		respondJSON(w, http.StatusOK, resp)
	})
	defer server.Close()

	client.DisableAutoPaging = true

	urlStr := client.requestUrl("/repositories/%s/%s/commits/main", "owner", "repo")
	_, err := client.executePaginated("GET", urlStr, "", nil)

	require.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func TestExecutePaginated_LimitPages(t *testing.T) {
	t.Parallel()
	callCount := 0

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		resp := paginatedResponse([]interface{}{
			map[string]interface{}{"name": "item"},
		})
		resp["next"] = "http://" + r.Host + r.URL.Path + "?page=" + fmt.Sprintf("%d", callCount+1)
		respondJSON(w, http.StatusOK, resp)
	})
	defer server.Close()

	client.LimitPages = 2

	urlStr := client.requestUrl("/repositories/%s/%s/commits/main", "owner", "repo")
	_, err := client.executePaginated("GET", urlStr, "", nil)

	require.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

func TestExecutePaginated_WithSpecificPage(t *testing.T) {
	t.Parallel()
	var receivedPage string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPage = r.URL.Query().Get("page")
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{
			map[string]interface{}{"name": "item"},
		}))
	})
	defer server.Close()

	page := 3
	urlStr := client.requestUrl("/repositories/%s/%s/commits/main", "owner", "repo")
	_, err := client.executePaginated("GET", urlStr, "", &page)

	require.NoError(t, err)
	assert.Equal(t, "3", receivedPage)
}

func TestExecutePaginated_CustomPagelen(t *testing.T) {
	t.Parallel()
	var receivedPagelen string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPagelen = r.URL.Query().Get("pagelen")
		respondJSON(w, http.StatusOK, paginatedResponse([]interface{}{}))
	})
	defer server.Close()

	client.Pagelen = 50

	urlStr := client.requestUrl("/repositories/%s/%s/commits/main", "owner", "repo")
	_, err := client.executePaginated("GET", urlStr, "", nil)

	require.NoError(t, err)
	assert.Equal(t, "50", receivedPagelen)
}

func TestExecuteRaw_Success(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("raw diff content"))
	})
	defer server.Close()

	urlStr := client.requestUrl("/repositories/%s/%s/diff/abc123", "owner", "repo")
	body, err := client.executeRaw("GET", urlStr, "")

	require.NoError(t, err)
	defer func() { _ = body.Close() }()
	data, _ := io.ReadAll(body)
	assert.Equal(t, "raw diff content", string(data))
}

func TestExecuteRaw_ErrorResponse(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("server error"))
	})
	defer server.Close()

	urlStr := client.requestUrl("/repositories/%s/%s/diff/abc123", "owner", "repo")
	_, err := client.executeRaw("GET", urlStr, "")

	assert.Error(t, err)
}

func TestUnexpectedHttpStatusCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"200 OK", http.StatusOK, false},
		{"201 Created", http.StatusCreated, false},
		{"202 Accepted", http.StatusAccepted, false},
		{"204 No Content", http.StatusNoContent, false},
		{"400 Bad Request", http.StatusBadRequest, true},
		{"401 Unauthorized", http.StatusUnauthorized, true},
		{"403 Forbidden", http.StatusForbidden, true},
		{"404 Not Found", http.StatusNotFound, true},
		{"500 Internal Server Error", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := unexpectedHttpStatusCode(tt.statusCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAddMaxDepthParam_Default(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	params := &url.Values{}
	client.addMaxDepthParam(params, nil)

	assert.Empty(t, params.Get("max_depth"))
}

func TestAddMaxDepthParam_Custom(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	params := &url.Values{}
	customDepth := 5
	client.addMaxDepthParam(params, &customDepth)

	assert.Equal(t, "5", params.Get("max_depth"))
}

func TestAddMaxDepthParam_ClientOverride(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	client.MaxDepth = 3
	params := &url.Values{}
	client.addMaxDepthParam(params, nil)

	assert.Equal(t, "3", params.Get("max_depth"))
}

func TestGetApiBaseURL(t *testing.T) {
	t.Parallel()
	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {})
	defer server.Close()

	apiBaseURL := client.GetApiBaseURL()

	assert.Contains(t, apiBaseURL, "/2.0")
}

func TestDoRequest_JsonUnmarshal(t *testing.T) {
	t.Parallel()
	var receivedContentType string

	client, server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedContentType = r.Header.Get("Content-Type")
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"owner": map[string]interface{}{
				"username": "nested-user",
				"links": map[string]interface{}{
					"self": "https://api.bitbucket.org/2.0/users/nested-user",
				},
			},
			"full_name": "owner/repo",
		})
	})
	defer server.Close()

	urlStr := client.requestUrl("/test")
	result, err := client.execute("POST", urlStr, `{"name":"test"}`)

	require.NoError(t, err)
	assert.Equal(t, "application/json", receivedContentType)
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "owner/repo", resultMap["full_name"])
	owner, ok := resultMap["owner"].(map[string]interface{})
	require.True(t, ok, "owner should be a nested map")
	assert.Equal(t, "nested-user", owner["username"])
	links, ok := owner["links"].(map[string]interface{})
	require.True(t, ok, "links should be a nested map")
	assert.Equal(t, "https://api.bitbucket.org/2.0/users/nested-user", links["self"])
}

func TestNewBasicAuth(t *testing.T) {
	t.Parallel()
	client, err := NewBasicAuth("user", "pass")

	require.NoError(t, err)
	assert.Equal(t, "user", client.Auth.user)
	assert.Equal(t, "pass", client.Auth.password)
	assert.Equal(t, DEFAULT_PAGE_LENGTH, client.Pagelen)
	assert.Equal(t, DEFAULT_MAX_DEPTH, client.MaxDepth)
	assert.Equal(t, DEFAULT_LIMIT_PAGES, client.LimitPages)
	assert.False(t, client.DisableAutoPaging)
	assert.NotNil(t, client.Repositories)
	assert.NotNil(t, client.Repositories.PullRequests)
	assert.NotNil(t, client.Repositories.Pipelines)
	assert.NotNil(t, client.Repositories.Repository)
	assert.NotNil(t, client.Repositories.Issues)
	assert.NotNil(t, client.Repositories.Commits)
	assert.NotNil(t, client.Repositories.Diff)
	assert.NotNil(t, client.Repositories.BranchRestrictions)
	assert.NotNil(t, client.Repositories.Webhooks)
	assert.NotNil(t, client.Repositories.Downloads)
	assert.NotNil(t, client.Repositories.DeployKeys)
	assert.NotNil(t, client.Users)
	assert.NotNil(t, client.Users.SSHKeys)
	assert.NotNil(t, client.User)
	assert.NotNil(t, client.Workspaces)
	assert.NotNil(t, client.Workspaces.Permissions)
	assert.NotNil(t, client.HttpClient)
}

func TestNewBasicAuthWithBaseUrlStr(t *testing.T) {
	t.Parallel()
	client, err := NewBasicAuthWithBaseUrlStr("user", "pass", "https://custom.bitbucket.org/2.0")

	require.NoError(t, err)
	assert.Contains(t, client.GetApiBaseURL(), "custom.bitbucket.org")
}

func TestNewBasicAuthWithBaseUrlStr_InvalidURL(t *testing.T) {
	t.Parallel()
	_, err := NewBasicAuthWithBaseUrlStr("user", "pass", "://invalid")

	assert.Error(t, err)
}

func TestNewOAuthbearerToken(t *testing.T) {
	t.Parallel()
	client, err := NewOAuthbearerToken("my-token")

	require.NoError(t, err)
	assert.Equal(t, "my-token", client.Auth.bearerToken)
	assert.Equal(t, DEFAULT_PAGE_LENGTH, client.Pagelen)
	assert.Equal(t, DEFAULT_LIMIT_PAGES, client.LimitPages)
	assert.False(t, client.DisableAutoPaging)
	assert.NotNil(t, client.HttpClient)
	assert.NotNil(t, client.Repositories)
	assert.NotNil(t, client.Users)
	assert.NotNil(t, client.Workspaces)
}

func TestNewOAuthbearerTokenWithBaseUrlStr(t *testing.T) {
	t.Parallel()
	client, err := NewOAuthbearerTokenWithBaseUrlStr("token", "https://custom.example.com/2.0")

	require.NoError(t, err)
	assert.Contains(t, client.GetApiBaseURL(), "custom.example.com")
	assert.Equal(t, "token", client.Auth.bearerToken)
	assert.Equal(t, DEFAULT_PAGE_LENGTH, client.Pagelen)
	assert.Equal(t, DEFAULT_LIMIT_PAGES, client.LimitPages)
	assert.False(t, client.DisableAutoPaging)
	assert.NotNil(t, client.HttpClient)
}

func TestSetApiBaseURL(t *testing.T) {
	t.Parallel()
	client, err := NewBasicAuth("user", "pass")
	require.NoError(t, err)

	newURL, _ := url.Parse("https://custom.example.com/2.0")
	client.SetApiBaseURL(*newURL)

	assert.Contains(t, client.GetApiBaseURL(), "custom.example.com")
}

func TestNewBasicAuthWithCaCert(t *testing.T) {
	t.Parallel()
	cert := generateSelfSignedCert(t)
	client, err := NewBasicAuthWithCaCert("user", "pass", cert)

	require.NoError(t, err)
	assert.Equal(t, "user", client.Auth.user)
	require.NotNil(t, client.HttpClient)
	transport, ok := client.HttpClient.Transport.(*http.Transport)
	require.True(t, ok, "transport should be *http.Transport")
	require.NotNil(t, transport.TLSClientConfig)
	assert.NotNil(t, transport.TLSClientConfig.RootCAs)
	assert.Equal(t, uint16(tls.VersionTLS12), transport.TLSClientConfig.MinVersion)
}

func TestNewBasicAuthWithBaseUrlStrCaCert(t *testing.T) {
	t.Parallel()
	cert := generateSelfSignedCert(t)
	client, err := NewBasicAuthWithBaseUrlStrCaCert("user", "pass", "https://custom.example.com/2.0", cert)

	require.NoError(t, err)
	assert.Contains(t, client.GetApiBaseURL(), "custom.example.com")
}

func TestNewBasicAuthWithBaseUrlStrCaCert_InvalidURL(t *testing.T) {
	t.Parallel()
	cert := generateSelfSignedCert(t)
	_, err := NewBasicAuthWithBaseUrlStrCaCert("user", "pass", "://invalid", cert)

	assert.Error(t, err)
}

func TestNewOAuthbearerTokenWithCaCert(t *testing.T) {
	t.Parallel()
	cert := generateSelfSignedCert(t)
	client, err := NewOAuthbearerTokenWithCaCert("my-token", cert)

	require.NoError(t, err)
	assert.Equal(t, "my-token", client.Auth.bearerToken)
	require.NotNil(t, client.HttpClient)
	transport, ok := client.HttpClient.Transport.(*http.Transport)
	require.True(t, ok, "transport should be *http.Transport")
	require.NotNil(t, transport.TLSClientConfig)
	assert.NotNil(t, transport.TLSClientConfig.RootCAs)
	assert.Equal(t, uint16(tls.VersionTLS12), transport.TLSClientConfig.MinVersion)
}

func TestNewOAuthbearerTokenWithBaseUrlStrCaCert(t *testing.T) {
	t.Parallel()
	cert := generateSelfSignedCert(t)
	client, err := NewOAuthbearerTokenWithBaseUrlStrCaCert("token", "https://custom.example.com/2.0", cert)

	require.NoError(t, err)
	assert.Contains(t, client.GetApiBaseURL(), "custom.example.com")
}

func TestNewOAuthbearerTokenWithBaseUrlStrCaCert_InvalidURL(t *testing.T) {
	t.Parallel()
	cert := generateSelfSignedCert(t)
	_, err := NewOAuthbearerTokenWithBaseUrlStrCaCert("token", "://invalid", cert)

	assert.Error(t, err)
}

func TestGetOAuthToken(t *testing.T) {
	t.Parallel()
	client, err := NewBasicAuth("user", "pass")
	require.NoError(t, err)

	token := client.GetOAuthToken()
	assert.False(t, token.Valid())
}

func TestNewOAuthbearerTokenWithBaseUrlStr_InvalidURL(t *testing.T) {
	t.Parallel()
	_, err := NewOAuthbearerTokenWithBaseUrlStr("token", "://invalid")

	assert.Error(t, err)
}

func TestAppendCaCerts_InvalidPEM(t *testing.T) {
	t.Parallel()
	_, err := appendCaCerts([]byte("not a valid PEM"))

	assert.Error(t, err)
}
