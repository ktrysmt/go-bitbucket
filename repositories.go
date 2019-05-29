package bitbucket

import (
	"net/url"
	"strconv"
)

//"github.com/k0kubun/pp"

type Repositories struct {
	c                  *Client
	Issues             *IssuesService
	PullRequests       *PullRequestsService
	Repository         *Repository
	Commits            *Commits
	Diff               *Diff
	BranchRestrictions *BranchRestrictions
	Webhooks           *Webhooks
	Downloads          *Downloads
	repositories
}

func safelyConvertToURIParameters(ro *RepositoriesOptions) string {
	params := url.Values{}

	if ro.Role != "" {
		params.Set("role", ro.Role)
	}

	if ro.ListOptions != nil {
		if ro.ListOptions.Page > 0 {
			pageString := strconv.Itoa(int(ro.ListOptions.Page))
			params.Set("page", pageString)
		}

		if ro.ListOptions.PageLen > 0 {
			pageLenString := strconv.Itoa(int(ro.ListOptions.PageLen))
			params.Set("pagelen", pageLenString)
		}
	}

	return params.Encode()
}

func (r *Repositories) ListForAccount(ro *RepositoriesOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s", ro.Owner)

	if params := safelyConvertToURIParameters(ro); params != "" {
		urlStr += "?" + params
	}
	return r.c.execute("GET", urlStr, "", "")
}

func (r *Repositories) ListForTeam(ro *RepositoriesOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s", ro.Owner)

	if params := safelyConvertToURIParameters(ro); params != "" {
		urlStr += "?" + params
	}
	return r.c.execute("GET", urlStr, "", "")
}

func (r *Repositories) ListPublic() (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/", "")
	return r.c.execute("GET", urlStr, "", "")
}
