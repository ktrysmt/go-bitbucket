package bitbucket

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"

	"github.com/k0kubun/pp"
	"github.com/mitchellh/mapstructure"
)

type PullRequest struct {
	Author struct {
		AccountID   string `mapstructure:"account_id"`
		DisplayName string `mapstructure:"display_name"`
		Links       struct {
			Avatar struct {
				Href string `mapstructure:"href"`
			} `mapstructure:"avatar"`
			HTML struct {
				Href string `mapstructure:"href"`
			} `mapstructure:"html"`
			Self struct {
				Href string `mapstructure:"href"`
			} `mapstructure:"self"`
		} `mapstructure:"links"`
		Nickname string `mapstructure:"nickname"`
		Type     string `mapstructure:"type"`
		UUID     string `mapstructure:"uuid"`
	} `mapstructure:"author"`
	CloseSourceBranch bool        `mapstructure:"close_source_branch"`
	ClosedBy          interface{} `mapstructure:"closed_by"`
	CommentCount      int32       `mapstructure:"comment_count"`
	CreatedOn         string      `mapstructure:"created_on"`
	Description       string      `mapstructure:"description"`
	Destination       struct {
		Branch struct {
			Name string `mapstructure:"name"`
		} `mapstructure:"branch"`
		Commit struct {
			Hash  string `mapstructure:"hash"`
			Links struct {
				HTML struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"html"`
				Self struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"self"`
			} `mapstructure:"links"`
			Type string `mapstructure:"type"`
		} `mapstructure:"commit"`
		Repository struct {
			FullName string `mapstructure:"full_name"`
			Links    struct {
				Avatar struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"avatar"`
				HTML struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"html"`
				Self struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"self"`
			} `mapstructure:"links"`
			Name string `mapstructure:"name"`
			Type string `mapstructure:"type"`
			UUID string `mapstructure:"uuid"`
		} `mapstructure:"repository"`
	} `mapstructure:"destination"`
	ID    int32 `mapstructure:"id"`
	Links struct {
		Activity struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"activity"`
		Approve struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"approve"`
		Comments struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"comments"`
		Commits struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"commits"`
		Decline struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"decline"`
		Diff struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"diff"`
		Diffstat struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"diffstat"`
		HTML struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"html"`
		Merge struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"merge"`
		Self struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"self"`
		Statuses struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"statuses"`
	} `mapstructure:"links"`
	MergeCommit  interface{} `mapstructure:"merge_commit"`
	Participants []struct {
		Approved       bool   `mapstructure:"approved"`
		ParticipatedOn string `mapstructure:"participated_on"`
		Role           string `mapstructure:"role"`
		Type           string `mapstructure:"type"`
		User           struct {
			AccountID   string `mapstructure:"account_id"`
			DisplayName string `mapstructure:"display_name"`
			Links       struct {
				Avatar struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"avatar"`
				HTML struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"html"`
				Self struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"self"`
			} `mapstructure:"links"`
			Nickname string `mapstructure:"nickname"`
			Type     string `mapstructure:"type"`
			UUID     string `mapstructure:"uuid"`
		} `mapstructure:"user"`
	} `mapstructure:"participants"`
	Reason   string `mapstructure:"reason"`
	Rendered struct {
		Description struct {
			HTML   string `mapstructure:"html"`
			Markup string `mapstructure:"markup"`
			Raw    string `mapstructure:"raw"`
			Type   string `mapstructure:"type"`
		} `mapstructure:"description"`
		Title struct {
			HTML   string `mapstructure:"html"`
			Markup string `mapstructure:"markup"`
			Raw    string `mapstructure:"raw"`
			Type   string `mapstructure:"type"`
		} `mapstructure:"title"`
	} `mapstructure:"rendered"`
	Reviewers []struct {
		AccountID   string `mapstructure:"account_id"`
		DisplayName string `mapstructure:"display_name"`
		Links       struct {
			Avatar struct {
				Href string `mapstructure:"href"`
			} `mapstructure:"avatar"`
			HTML struct {
				Href string `mapstructure:"href"`
			} `mapstructure:"html"`
			Self struct {
				Href string `mapstructure:"href"`
			} `mapstructure:"self"`
		} `mapstructure:"links"`
		Nickname string `mapstructure:"nickname"`
		Type     string `mapstructure:"type"`
		UUID     string `mapstructure:"uuid"`
	} `mapstructure:"reviewers"`
	Source struct {
		Branch struct {
			Name string `mapstructure:"name"`
		} `mapstructure:"branch"`
		Commit struct {
			Hash  string `mapstructure:"hash"`
			Links struct {
				HTML struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"html"`
				Self struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"self"`
			} `mapstructure:"links"`
			Type string `mapstructure:"type"`
		} `mapstructure:"commit"`
		Repository struct {
			FullName string `mapstructure:"full_name"`
			Links    struct {
				Avatar struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"avatar"`
				HTML struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"html"`
				Self struct {
					Href string `mapstructure:"href"`
				} `mapstructure:"self"`
			} `mapstructure:"links"`
			Name string `mapstructure:"name"`
			Type string `mapstructure:"type"`
			UUID string `mapstructure:"uuid"`
		} `mapstructure:"repository"`
	} `mapstructure:"source"`
	State   string `mapstructure:"state"`
	Summary struct {
		HTML   string `mapstructure:"html"`
		Markup string `mapstructure:"markup"`
		Raw    string `mapstructure:"raw"`
		Type   string `mapstructure:"type"`
	} `mapstructure:"summary"`
	TaskCount int32  `mapstructure:"task_count"`
	Title     string `mapstructure:"title"`
	Type      string `mapstructure:"type"`
	UpdatedOn string `mapstructure:"updated_on"`
}

type PullRequestsResp struct {
	Page    int32         `mapstructure:"page"`
	Pagelen int32         `mapstructure:"pagelen"`
	Size    int32         `mapstructure:"size"`
	Values  []PullRequest `mapstructure:"values"`
}

type PullRequests struct {
	c *Client
}

func (p *PullRequests) Create(po *PullRequestsOptions) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := p.c.requestUrl("/repositories/%s/%s/pullrequests/", po.Owner, po.RepoSlug)
	return p.c.execute("POST", urlStr, data)
}

func (p *PullRequests) Update(po *PullRequestsOptions) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID
	return p.c.execute("PUT", urlStr, data)
}

func (p *PullRequests) Gets(po *PullRequestsOptions) (*PullRequestsResp, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/"

	if po.States != nil && len(po.States) != 0 {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		for _, state := range po.States {
			query.Set("state", state)
		}
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Query != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("q", po.Query)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Sort != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("sort", po.Sort)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	resp, err := p.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	pullRequestsResponseMap, ok := resp.(map[string]interface{})
	if !ok {
		return nil, errors.New("Not a valid format")
	}

	repoArray := pullRequestsResponseMap["values"].([]interface{})
	var prs []PullRequest
	for _, repoEntry := range repoArray {
		var pr PullRequest
		err := mapstructure.Decode(repoEntry, &pr)
		if err == nil {
			prs = append(prs, pr)
		}
	}

	page, ok := pullRequestsResponseMap["page"].(float64)
	if !ok {
		page = 0
	}

	pagelen, ok := pullRequestsResponseMap["pagelen"].(float64)
	if !ok {
		pagelen = 0
	}

	size, ok := pullRequestsResponseMap["size"].(float64)
	if !ok {
		size = 0
	}

	repositories := PullRequestsResp{
		Page:    int32(page),
		Pagelen: int32(pagelen),
		Size:    int32(size),
		Values:  prs,
	}
	return &repositories, nil
}

func (p *PullRequests) Get(po *PullRequestsOptions) (*PullRequest, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID
	resp, err := p.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodePullrequest(resp)
}

func (p *PullRequests) Activities(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/activity"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) Activity(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/activity"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) Commits(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/commits"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) Patch(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/patch"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) Diff(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/diff"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) Merge(po *PullRequestsOptions) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/merge"
	return p.c.execute("POST", urlStr, data)
}

func (p *PullRequests) Decline(po *PullRequestsOptions) (interface{}, error) {
	data := p.buildPullRequestBody(po)
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/decline"
	return p.c.execute("POST", urlStr, data)
}

func (p *PullRequests) GetComments(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/comments/"
	return p.c.execute("GET", urlStr, "")
}

func (p *PullRequests) GetComment(po *PullRequestsOptions) (interface{}, error) {
	urlStr := p.c.GetApiBaseURL() + "/repositories/" + po.Owner + "/" + po.RepoSlug + "/pullrequests/" + po.ID + "/comments/" + po.CommentID
	return p.c.execute("GET", urlStr, "")
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
		body["reviewers"] = make([]map[string]string, n)
		for i, uuid := range po.Reviewers {
			body["reviewers"].([]map[string]string)[i] = map[string]string{"uuid": uuid}
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

func decodePullrequest(pullRequest interface{}) (*PullRequest, error) {
	userMap := pullRequest.(map[string]interface{})

	if userMap["type"] == "error" {
		return nil, DecodeError(userMap)
	}

	var pr = new(PullRequest)
	err := mapstructure.Decode(userMap, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
