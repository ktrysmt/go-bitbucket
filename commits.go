package bitbucket

import "net/url"

type Commits struct {
	c *Client
}

func (cm *Commits) GetCommits(cmo *CommitsOptions) interface{} {
	url := cm.c.requestUrl("/repositories/%s/%s/commits/%s", cmo.Owner, cmo.Repo_slug, cmo.Branchortag)
	url += cm.buildCommitsQuery(cmo.Include, cmo.Exclude)
	return cm.c.execute("GET", url, "")
}

func (cm *Commits) GetCommit(cmo *CommitsOptions) interface{} {
	url := cm.c.requestUrl("/repositories/%s/%s/commit/%s", cmo.Owner, cmo.Repo_slug, cmo.Revision)
	return cm.c.execute("GET", url, "")
}

func (cm *Commits) GetCommitComments(cmo *CommitsOptions) interface{} {
	url := cm.c.requestUrl("/repositories/%s/%s/commit/%s/comments", cmo.Owner, cmo.Repo_slug, cmo.Revision)
	return cm.c.execute("DELETE", url, "")
}

func (cm *Commits) GetCommitComment(cmo *CommitsOptions) interface{} {
	url := cm.c.requestUrl("/repositories/%s/%s/commit/%s/comments/%s", cmo.Owner, cmo.Repo_slug, cmo.Revision, cmo.Comment_id)
	return cm.c.execute("GET", url, "")
}

func (cm *Commits) GetCommitStatus(cmo *CommitsOptions, commitStatusKey string) interface{} {
	url := cm.c.requestUrl("/repositories/%s/%s/commit/%s/statuses/build/%s", cmo.Owner, cmo.Repo_slug, cmo.Revision, commitStatusKey)
	return cm.c.execute("GET", url, "")
}

func (cm *Commits) GiveApprove(cmo *CommitsOptions) interface{} {
	url := cm.c.requestUrl("/repositories/%s/%s/commit/%s/approve", cmo.Owner, cmo.Repo_slug, cmo.Revision)
	return cm.c.execute("POST", url, "")
}

func (cm *Commits) RemoveApprove(cmo *CommitsOptions) interface{} {
	url := cm.c.requestUrl("/repositories/%s/%s/commit/%s/approve", cmo.Owner, cmo.Repo_slug, cmo.Revision)
	return cm.c.execute("DELETE", url, "")
}

func (cm *Commits) buildCommitsQuery(include, exclude string) string {

	p := url.Values{}

	if include != "" {
		p.Add("include", include)
	}
	if exclude != "" {
		p.Add("exclude", exclude)
	}

	return p.Encode()
}
