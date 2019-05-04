package bitbucket

import (
	"fmt"
	"net/url"
	"strconv"
)

//"github.com/k0kubun/pp"

type Repositories struct {
	c                  *Client
	PullRequests       *PullRequests
	Repository         *Repository
	Commits            *Commits
	Diff               *Diff
	BranchRestrictions *BranchRestrictions
	Webhooks           *Webhooks
	Downloads		   *Downloads
	repositories
}

func safelyConvertToURIParameters(ro *RepositoriesOptions) string {
	params := ""

	if ro.Role != "" {
		params += fmt.Sprintf("role=%s&", url.QueryEscape(ro.Role))
	}

	if ro.ListOptions != nil {
		if ro.ListOptions.Page > 0 {
			pageString := strconv.Itoa(int(ro.ListOptions.Page))
			params += fmt.Sprintf("page=%s&", url.QueryEscape(pageString))
		}

		if ro.ListOptions.PageLen > 0 {
			pageLenString := strconv.Itoa(int(ro.ListOptions.PageLen))
			params += fmt.Sprintf("pagelen=%s&", url.QueryEscape(pageLenString))
		}
	}

	return params
}

func (r *Repositories) ListForAccount(ro *RepositoriesOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s", ro.Owner)

	if params := safelyConvertToURIParameters(ro); params != "" {
		urlStr += "?" + params
	}

	return r.c.execute("GET", urlStr, "")
}

func (r *Repositories) ListForTeam(ro *RepositoriesOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s", ro.Owner)

	if params := safelyConvertToURIParameters(ro); params != "" {
		urlStr += "?" + params
	}

	return r.c.execute("GET", urlStr, "")
}

func (r *Repositories) ListPublic() (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/", "")
	return r.c.execute("GET", urlStr, "")
}
