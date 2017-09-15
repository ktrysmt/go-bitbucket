package bitbucket

//"github.com/k0kubun/pp"

type Repositories struct {
	c                  *Client
	PullRequests       *PullRequests
	Repository         *Repository
	Commits            *Commits
	Diff               *Diff
	BranchRestrictions *BranchRestrictions
	Webhooks           *Webhooks
	repositories
}

func (r *Repositories) ListForAccount(ro *RepositoriesOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s", ro.Owner)
	if ro.Role != "" {
		urlStr += "?role=" + ro.Role
	}
	return r.c.execute("GET", urlStr, "")
}

func (r *Repositories) ListForTeam(ro *RepositoriesOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s", ro.Owner)
	if ro.Role != "" {
		urlStr += "?role=" + ro.Role
	}
	return r.c.execute("GET", urlStr, "")
}

func (r *Repositories) ListPublic() (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/", "")
	return r.c.execute("GET", urlStr, "")
}
