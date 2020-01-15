package bitbucket

import (
	"net/url"
	"encoding/json"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

type Commits struct {
	c *Client
}

type CommitStatus struct {
	Key                    string
	Refname                string
	Url                    string
	State                  string
	Name                   string
	Description            string
	Type                   string
	Links                  map[string]interface{}
}

type CommitStatusesResponse struct {
	Page	        int
	Pagelen	        int
	MaxDepth        int
	Size            int
	Next            string
	Previous        string
	CommitStatuses  []CommitStatus
}

func (cm *Commits) GetCommits(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commits/%s", cmo.Owner, cmo.RepoSlug, cmo.Branchortag)
	urlStr += cm.buildCommitsQuery(cmo.Include, cmo.Exclude)
	return cm.c.execute("GET", urlStr, "")
}

func (cm *Commits) GetCommit(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	return cm.c.execute("GET", urlStr, "")
}

func (cm *Commits) GetCommitComments(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/comments", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	return cm.c.execute("DELETE", urlStr, "")
}

func (cm *Commits) GetCommitComment(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/comments/%s", cmo.Owner, cmo.RepoSlug, cmo.Revision, cmo.CommentID)
	return cm.c.execute("GET", urlStr, "")
}

func (cm *Commits) GetCommitStatuses(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/statuses", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	return cm.c.execute("GET", urlStr, "")
}

func (cm *Commits) GetCommitStatus(cmo *CommitsOptions, commitStatusKey string) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/statuses/build/%s", cmo.Owner, cmo.RepoSlug, cmo.Revision, commitStatusKey)
	return cm.c.execute("GET", urlStr, "")
}

func (cm *Commits) GiveApprove(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/approve", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	return cm.c.execute("POST", urlStr, "")
}

func (cm *Commits) RemoveApprove(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/approve", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	return cm.c.execute("DELETE", urlStr, "")
}


func (cm *Commits) CreateCommitStatus(cmo *CommitsOptions, cso *CommitStatusOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/statuses/build", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	data, err := json.Marshal(cso)
	if err != nil {
		return nil, err
	}
	return cm.c.execute("POST", urlStr, string(data))
}


func (cm *Commits) buildCommitsQuery(include, exclude string) string {

	p := url.Values{}

	if include != "" {
		p.Add("include", include)
	}
	if exclude != "" {
		p.Add("exclude", exclude)
	}

	if res := p.Encode(); len(res) > 0 {
		return "?" + res
	}
	return ""
}

func (cm *Commits) GetCommitStatusesV2(cmo *CommitsOptions) (*CommitStatusesResponse, error) {
	params := url.Values{}
	if cmo.PageNum > 0 {
		params.Add("page", strconv.Itoa(cmo.PageNum))
	}

	if cmo.Pagelen > 0 {
		params.Add("pagelen", strconv.Itoa(cmo.Pagelen))
	}

	if cmo.MaxDepth > 0 {
		params.Add("max_depth", strconv.Itoa(cmo.MaxDepth))
	}

	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/statuses?%s", cmo.Owner, cmo.RepoSlug, cmo.Revision, params.Encode())
	response, err :=  cm.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodeCommitStatuses(response)
}

func decodeCommitStatuses(response interface{}) (*CommitStatusesResponse, error) {
	var commitStatusesResponseMap map[string]interface{}
	err := json.Unmarshal(response.([]byte), &commitStatusesResponseMap)
	if err != nil {
		return nil, err
	}

	commitStatusesArray := commitStatusesResponseMap["values"].([]interface{})
	var css []CommitStatus
	for _, commitStatusEntry := range commitStatusesArray {
		var commitStatus CommitStatus
		err = mapstructure.Decode(commitStatusEntry, &commitStatus)
		if err == nil {
			css = append(css, commitStatus)
		}
	}

	page, ok := commitStatusesResponseMap["page"].(float64)
	if !ok {
		page = 0
	}

	pagelen, ok := commitStatusesResponseMap["pagelen"].(float64)
	if !ok {
		pagelen = 0
	}
	max_depth, ok := commitStatusesResponseMap["max_depth"].(float64)
	if !ok {
		max_depth = 0
	}
	size, ok := commitStatusesResponseMap["size"].(float64)
	if !ok {
		size = 0
	}

	next, ok := commitStatusesResponseMap["next"].(string)
	if !ok {
		next = ""
	}

	commitStatuses :=  CommitStatusesResponse {
		Page:     int(page),
		Pagelen:  int(pagelen),
		MaxDepth: int(max_depth),
		Size:     int(size),
		Next:     next,
		CommitStatuses: css,
	}
	return &commitStatuses, nil
}
