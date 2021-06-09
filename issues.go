package bitbucket

import (
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"github.com/k0kubun/pp"
)

type Issues struct {
	c *Client
}

func (p *Issues) Gets(io *IssuesOptions) (interface{}, error) {
	url, err := url.Parse(p.c.GetApiBaseURL() + "/repositories/" + io.Owner + "/" + io.RepoSlug + "/issues/")
	if err != nil {
		return nil, err
	}

	if io.States != nil && len(io.States) != 0 {
		query := url.Query()
		for _, state := range io.States {
			query.Set("state", state)
		}
		url.RawQuery = query.Encode()
	}

	if io.Query != "" {
		query := url.Query()
		query.Set("q", io.Query)
		url.RawQuery = query.Encode()
	}

	if io.Sort != "" {
		query := url.Query()
		query.Set("sort", io.Sort)
		url.RawQuery = query.Encode()
	}

	return p.c.execute("GET", url.String(), "")
}

func (p *Issues) Get(io *IssuesOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + io.Owner + "/" + io.RepoSlug + "/issues/" + io.ID
	return p.c.execute("GET", urlStr, "")
}

func (p *Issues) Delete(io *IssuesOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + io.Owner + "/" + io.RepoSlug + "/issues/" + io.ID
	return p.c.execute("DELETE", urlStr, "")
}

func (p *Issues) Update(io *IssuesOptions) (interface{}, error) {
	data := p.buildIssueBody(io)
	urlStr := p.c.requestUrl("/repositories/%s/%s/issues%s", io.Owner, io.RepoSlug, io.ID)
	return p.c.execute("POST", urlStr, data)
}

func (p *Issues) Create(io *IssuesOptions) (interface{}, error) {
	data := p.buildIssueBody(io)
	urlStr := p.c.requestUrl("/repositories/%s/%s/issues", io.Owner, io.RepoSlug)
	return p.c.execute("POST", urlStr, data)
}

func (p *Issues) GetVote(io *IssuesOptions) (bool, error) {
	// A 404 indicates that the user hasn't voted
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + io.Owner + "/" + io.RepoSlug + "/issues/" + io.ID + "/vote"
	_, err := p.c.execute("GET", urlStr, "")
	if strings.HasPrefix(err.Error(), "404") {
		return false, nil
	}
	return true, err
}

func (p *Issues) PutVote(io *IssuesOptions) error {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + io.Owner + "/" + io.RepoSlug + "/issues/" + io.ID + "/vote"
	_, err := p.c.execute("PUT", urlStr, "")
	return err
}

func (p *Issues) DeleteVote(io *IssuesOptions) error {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + io.Owner + "/" + io.RepoSlug + "/issues/" + io.ID + "/vote"
	_, err := p.c.execute("DELETE", urlStr, "")
	return err
}

func (p *Issues) buildIssueBody(io *IssuesOptions) string {
	body := map[string]interface{}{}

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	// This feld is required
	body["title"] = io.Title

	if io.Content != "" {
		body["content"] = map[string]interface{}{
			"raw": io.Content,
		}
	}

	if io.State != "" {
		body["state"] = io.State
	}

	if io.Kind != "" {
		body["kind"] = io.Kind
	}

	if io.Priority != "" {
		body["priority"] = io.Priority
	}

	if io.Milestone != "" {
		body["milestone"] = map[string]interface{}{
			"name": io.Milestone,
		}
	}

	if io.Component != "" {
		body["component"] = map[string]interface{}{
			"name": io.Component,
		}
	}

	if io.Version != "" {
		body["version"] = map[string]interface{}{
			"name": io.Component,
		}
	}
	if io.Assignee != "" {
		body["assignee"] = map[string]interface{}{
			"uuid": io.Assignee,
		}
	}

	return string(data)
}
