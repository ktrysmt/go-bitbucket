package bitbucket

import (
	"encoding/json"
	"net/url"
)

type Commits struct {
	c *Client
}

func (cm *Addons) GetCommits(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commits/%s", cmo.Owner, cmo.RepoSlug, cmo.Branchortag)
	urlStr += cm.buildCommitsQuery(cmo.Include, cmo.Exclude)
	return cm.c.executePaginated("GET", urlStr, "")
}

func (cm *Addons) GetCommit(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	return cm.c.execute("GET", urlStr, "")
}

func (cm *Addons) GetCommitComments(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/comments", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	return cm.c.executePaginated("GET", urlStr, "")
}

func (cm *Addons) GetCommitComment(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/comments/%s", cmo.Owner, cmo.RepoSlug, cmo.Revision, cmo.CommentID)
	return cm.c.execute("GET", urlStr, "")
}

func (cm *Addons) GetCommitStatuses(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/statuses", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	return cm.c.executePaginated("GET", urlStr, "")
}

func (cm *Addons) GetCommitStatus(cmo *CommitsOptions, commitStatusKey string) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/statuses/build/%s", cmo.Owner, cmo.RepoSlug, cmo.Revision, commitStatusKey)
	return cm.c.execute("GET", urlStr, "")
}

func (cm *Addons) GiveApprove(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/approve", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	return cm.c.execute("POST", urlStr, "")
}

func (cm *Addons) RemoveApprove(cmo *CommitsOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/approve", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	return cm.c.execute("DELETE", urlStr, "")
}

func (cm *Addons) CreateCommitStatus(cmo *CommitsOptions, cso *CommitStatusOptions) (interface{}, error) {
	urlStr := cm.c.requestUrl("/repositories/%s/%s/commit/%s/statuses/build", cmo.Owner, cmo.RepoSlug, cmo.Revision)
	data, err := json.Marshal(cso)
	if err != nil {
		return nil, err
	}
	return cm.c.execute("POST", urlStr, string(data))
}

func (cm *Addons) buildCommitsQuery(include, exclude string) string {

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
