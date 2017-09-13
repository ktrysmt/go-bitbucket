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
	url := r.c.requestUrl("/repositories/%s", ro.Owner)
	if ro.Role != "" {
		url += "?role=" + ro.Role
	}
	return r.c.execute("GET", url, "")
}

func (r *Repositories) ListForTeam(ro *RepositoriesOptions) (interface{}, error) {
	url := r.c.requestUrl("/repositories/%s", ro.Owner)
	if ro.Role != "" {
		url += "?role=" + ro.Role
	}
	return r.c.execute("GET", url, "")
}

func (r *Repositories) ListPublic() (interface{}, error) {
	url := r.c.requestUrl("/repositories/", "")
	return r.c.execute("GET", url, "")
}
