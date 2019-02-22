package bitbucket

import (
	"encoding/json"
	"fmt"
	"os"
)

// PullRequestService handles communication with the pull requests related
// methods of the Bitbucket API.
//
// Bitbucket API docs: https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories/%7Busername%7D/%7Brepo_slug%7D/pullrequests
type PullRequestsService struct {
	client *Client
}

// PullRequest represents a list of Bitbucket pull request.
//
// Bitbucket API docs: https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories/%7Busername%7D/%7Brepo_slug%7D/pullrequests
type PullRequests struct {
	Page    int            `json:"page,omitempty"`
	Size    int            `json:"size,omitempty"`
	PageLen int            `json:"pagelen,omitempty"`
	Values  *[]PullRequest `json:"values,omitempty"`
}

// PullRequest represents a Bitbucket pull request.
//
// Bitbucket API docs: https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories/%7Busername%7D/%7Brepo_slug%7D/pullrequests
type PullRequest struct {
	Rendered struct {
		Description struct {
			Raw    string `json:"raw,omitempty"`
			Markup string `json:"markup,omitempty"`
			HTML   string `json:"html,omitempty"`
			Type   string `json:"type,omitempty"`
		} `json:"description,omitempty"`
		Title struct {
			Raw    string `json:"raw,omitempty"`
			Markup string `json:"markup,omitempty"`
			HTML   string `json:"html,omitempty"`
			Type   string `json:"type,omitempty"`
		} `json:"title,omitempty"`
	} `json:"rendered,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Links       struct {
		Html struct {
			Href string `json:"href,omitempty"`
		} `json:"html,omitempty"`
	} `json:"html,omitempty"`
	Title             string `json:"title,omitempty"`
	CloseSourceBranch bool   `json:"close_source_branch,omitempty"`
	ID                int64  `json:"id,omitempty"`
	Destination       struct {
		Commit struct {
			Hash string `json:"hash,omitempty"`
		} `json:"commit,omitempty"`
		Repository struct {
			Name     string `json:"name,omitempty"`
			FullName string `json:"full_name,omitempty"`
			Uuid     string `json:"uuid,omitempty"`
		} `json:"repository,omitempty"`
		Branch struct {
			Name string `json:"name,omitempty"`
		} `json:"branch,omitempty"`
	} `json:"destination,omitempty"`
	Summary struct {
		Raw    string `json:"raw,omitempty"`
		Markup string `json:"markup,omitempty"`
		Html   string `json:"html,omitempty"`
		Type   string `json:"type,omitempty"`
	} `json:"summary,omitempty"`
	Source struct {
		Commit struct {
			Hash string `json:"hash,omitempty"`
		}
		Repository struct {
			Name     string `json:"name,omitempty"`
			FullName string `json:"full_name,omitempty"`
			Uuid     string `json:"uuid,omitempty"`
		} `json:"repository,omitempty"`
		Branch struct {
			Name string `json:"name,omitempty"`
		} `json:"destination,omitempty"`
	} `json:"source,omitempty"`
	CommentCount int    `json:"comment_count,omitempty"`
	State        string `json:"state,omitempty"`
	TaskCount    int    `json:"task_count,omitempty"`
	Reason       string `json:"reason,omitempty"`
}

type PullRequestsOpts struct {
	ID                string   `json:"id"`
	State             string   `json:"state"`
	CommentID         string   `json:"comment_id"`
	Owner             string   `json:"owner"`
	RepoSlug          string   `json:"repo_slug"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	CloseSourceBranch bool     `json:"close_source_branch"`
	SourceBranch      string   `json:"source_branch"`
	SourceRepository  string   `json:"source_repository"`
	DestinationBranch string   `json:"destination_branch"`
	DestinationCommit string   `json:"destination_repository"`
	Message           string   `json:"message"`
	Reviewers         []string `json:"reviewers"`
}

func (p *PullRequestsService) List(owner, repo, opts string) (*PullRequests, *Response, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + owner + "/" + repo + "/pullrequests"

	result := new(PullRequests)
	response, err := p.client.executeNew("GET", urlStr, result, nil, opts)

	return result, response, err
}

func (p *PullRequestsService) Get(owner, repo, id string) (*PullRequest, *Response, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + owner + "/" + repo + "/pullrequests/" + id

	result := new(PullRequest)
	response, err := p.client.executeNew("GET", urlStr, result, nil, "")

	return result, response, err
}

func (p *PullRequestsService) Create(owner, repo string, po *PullRequestsOpts) (*PullRequest, *Response, error) {
	urlStr := p.client.requestUrl("/repositories/%s/%s/pullrequests/", owner, repo)

	result := new(PullRequest)
	response, err := p.client.executeNew("POST", urlStr, result, po, "")

	return result, response, err
}

func (p *PullRequestsService) Patch(po *PullRequestsOpts) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/patch"
	return p.client.execute("GET", urlStr, "", "")
}

func (p *PullRequestsService) Diff(po *PullRequestsOpts) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/diff"
	return p.client.execute("GET", urlStr, "", "")
}

func (p *PullRequestsService) Merge(po *PullRequestsOpts) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/merge"
	return p.client.execute("POST", urlStr, data, "")
}

func (p *PullRequestsService) Decline(po *PullRequestsOpts) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/decline"
	return p.client.execute("POST", urlStr, data, "")
}

func (p *PullRequestsService) GetComments(po *PullRequestsOpts) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/comments/"
	return p.client.execute("GET", urlStr, "", "")
}

func (p *PullRequestsService) GetComment(po *PullRequestsOpts) (interface{}, error) {
	urlStr := GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/comments/" + po.CommentID
	return p.client.execute("GET", urlStr, "", "")
}

func (p *PullRequestsService) buildPullRequestBody(po *PullRequestsOpts) string {

	body := map[string]interface{}{}
	body["source"] = map[string]interface{}{}
	body["destination"] = map[string]interface{}{}
	body["reviewers"] = []map[string]string{}
	body["title"] = ""
	body["description"] = ""
	body["message"] = ""
	body["close_source_branch"] = false

	if n := len(po.Reviewers); n > 0 {
		body["reviewers"] = make([]map[string]string, n)
		for i, user := range po.Reviewers {
			body["reviewers"].([]map[string]string)[i] = map[string]string{"username": user}
		}
	}

	if po.SourceBranch != "" {
		body["source"].(map[string]interface{})["branch"] = map[string]string{"name": po.SourceBranch}
	}

	if po.SourceRepository != "" {
		body["source"].(map[string]interface{})["repository"] = map[string]interface{}{"full_name": po.SourceRepository}
	}

	if po.DestinationBranch != "" {
		body["destination"].(map[string]interface{})["branch"] = map[string]interface{}{"name": po.DestinationBranch}
	}

	if po.DestinationCommit != "" {
		body["destination"].(map[string]interface{})["commit"] = map[string]interface{}{"hash": po.DestinationCommit}
	}

	if po.Title != "" {
		body["title"] = po.Title
	}

	if po.Description != "" {
		body["description"] = po.Description
	}

	if po.Message != "" {
		body["message"] = po.Message
	}

	if po.CloseSourceBranch == true || po.CloseSourceBranch == false {
		body["close_source_branch"] = po.CloseSourceBranch
	}

	data, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		os.Exit(9)
	}

	return string(data)
}
