package bitbucket

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"os"
)

type Issues struct {
	c *Client

	PageLen int     `json:"pagelen,omitempty"`
	Values  []Issue `json:"values,omitempty"`
}

type Issue struct {
	Priority  string  `json:"priority,omitempty"`
	Kind      string  `json:"kind,omitempty"`
	Title     string  `json:"title,omitempty"`
	Votes     int     `json:"votes,omitempty"`
	Watches   int     `json:"watches,omitempty"`
	Content   Content `json:"content,omitempty"`
	State     string  `json:"state,omitempty"`
	IssueType string  `json:"state,omitempty"`
	ID        int64   `json:"id,omitempty"`
}

type Content struct {
	Raw    string `json:"raw,omitempty"`
	Markup string `json:"markup,omitempty"`
	Html   string `json:"html,omitempty"`
	Type   string `json:"type,omitempty"`
}

func (i *Issues) List(owner, repoSlug string) (*Issues, error) {
	urlStr := i.c.requestUrl("/repositories/%s/%s/issues", owner, repoSlug)
	response, err := i.c.execute("GET", urlStr, "", "")
	if err != nil {
		return nil, err
	}

	return decodeIssues(response)
}

func (i *Issues) Get(owner, repoSlug, issueId string) (*Issue, error) {
	urlStr := i.c.requestUrl("/repositories/%s/%s/issues/%s", owner, repoSlug, issueId)
	response, err := i.c.execute("GET", urlStr, "", "")
	if err != nil {
		return nil, err
	}

	return decodeIssue(response)
}

func (i *Issues) Create(owner, repoSlug string, io *IssueOptions) (*Issue, error) {
	data := i.buildIssueBody(io)

	urlStr := i.c.requestUrl("/repositories/%s/%s/issues", owner, repoSlug)
	response, err := i.c.execute("POST", urlStr, data, "")
	if err != nil {
		return nil, err
	}

	return decodeIssue(response)
}

func (i *Issues) buildIssueBody(io *IssueOptions) string {

	body := map[string]interface{}{}

	if io.Title != "" {
		body["title"] = io.Title
	}

	if io.Kind != "" {
		body["kind"] = io.Kind
	}
	if io.Priority != "" {
		body["priority"] = io.Priority
	}
	if io.Content.Raw != "" {
		body["content"] = map[string]string{
			"raw": io.Content.Raw,
		}
	}

	return i.buildJsonBody(body)
}

func (i *Issues) buildJsonBody(body map[string]interface{}) string {

	data, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		os.Exit(9)
	}

	return string(data)
}

func decodeIssue(repoResponse interface{}) (*Issue, error) {
	repoMap := repoResponse.(map[string]interface{})

	if repoMap["type"] == "error" {
		return nil, DecodeError(repoMap)
	}

	var issue = new(Issue)
	err := mapstructure.Decode(repoMap, issue)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func decodeIssues(repoResponse interface{}) (*Issues, error) {
	repoMap := repoResponse.(map[string]interface{})

	if repoMap["type"] == "error" {
		return nil, DecodeError(repoMap)
	}

	var issues = new(Issues)
	err := mapstructure.Decode(repoMap, issues)
	if err != nil {
		return nil, err
	}

	return issues, nil
}
