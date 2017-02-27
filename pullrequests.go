package bitbucket

import (
	"encoding/json"
	"os"

	"github.com/k0kubun/pp"
)

type PullRequests struct {
	c *Client
}

func (p *PullRequests) Create(po *PullRequestsOptions) interface{} {
	data := p.buildPullRequestBody(po)
	url := p.c.requestUrl("/repositories/%s/%s/pullrequests/", po.Owner, po.RepoSlug)
	return p.c.execute("POST", url, data)
}

func (p *PullRequests) Update(po *PullRequestsOptions) interface{} {
	data := p.buildPullRequestBody(po)
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID
	return p.c.execute("PUT", url, data)
}

func (p *PullRequests) Gets(po *PullRequestsOptions) interface{} {
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/"
	return p.c.execute("GET", url, "")
}

func (p *PullRequests) Get(po *PullRequestsOptions) interface{} {
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID
	return p.c.execute("GET", url, "")
}

func (p *PullRequests) Activities(po *PullRequestsOptions) interface{} {
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/activity"
	return p.c.execute("GET", url, "")
}

func (p *PullRequests) Activity(po *PullRequestsOptions) interface{} {
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/activity"
	return p.c.execute("GET", url, "")
}

func (p *PullRequests) Commits(po *PullRequestsOptions) interface{} {
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/commits"
	return p.c.execute("GET", url, "")
}

func (p *PullRequests) Patch(po *PullRequestsOptions) interface{} {
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/patch"
	return p.c.execute("GET", url, "")
}

func (p *PullRequests) Diff(po *PullRequestsOptions) interface{} {
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/diff"
	return p.c.execute("GET", url, "")
}

func (p *PullRequests) Merge(po *PullRequestsOptions) interface{} {
	data := p.buildPullRequestBody(po)
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/merge"
	return p.c.execute("POST", url, data)
}

func (p *PullRequests) Decline(po *PullRequestsOptions) interface{} {
	data := p.buildPullRequestBody(po)
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/decline"
	return p.c.execute("POST", url, data)
}

func (p *PullRequests) GetComments(po *PullRequestsOptions) interface{} {
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/comments/"
	return p.c.execute("GET", url, "")
}

func (p *PullRequests) GetComment(po *PullRequestsOptions) interface{} {
	url := GetAPIBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/comments/" + po.CommentID
	return p.c.execute("GET", url, "")
}

func (p *PullRequests) buildPullRequestBody(po *PullRequestsOptions) string {

	body := map[string]interface{}{}
	body["source"] = map[string]interface{}{}
	body["destination"] = map[string]interface{}{}
	body["reviewers"] = []map[string]string{}
	body["title"] = ""
	body["description"] = ""
	body["message"] = ""
	body["close_source_branch"] = false

	if n := len(po.Reviewers); n > 0 {
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
		pp.Println(err)
		os.Exit(9)
	}

	return string(data)
}
