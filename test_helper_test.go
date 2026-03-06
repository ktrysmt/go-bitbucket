package bitbucket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// setupMockServer creates a httptest server and a Client pointing at it.
// The handler receives the actual HTTP requests made by the Client.
// Caller must call server.Close() when done.
func setupMockServer(handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client, _ := NewBasicAuth("test-user", "test-pass")
	serverURL, _ := url.Parse(server.URL + "/2.0")
	client.SetApiBaseURL(*serverURL)
	client.HttpClient = server.Client()
	return client, server
}

// respondJSON writes a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}

// paginatedResponse creates a standard Bitbucket paginated response body.
func paginatedResponse(values []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"page":    1,
		"pagelen": 10,
		"size":    len(values),
		"values":  values,
	}
}
