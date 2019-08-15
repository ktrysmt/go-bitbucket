package bitbucket

import "github.com/mitchellh/mapstructure"

//"github.com/k0kubun/pp"

type Repositories struct {
	c                  *Client
	PullRequests       *PullRequests
	Repository         *Repository
	Commits            *Commits
	Diff               *Diff
	BranchRestrictions *BranchRestrictions
	Webhooks           *Webhooks
	Downloads          *Downloads
	repositories
}

type RepositoriesRes struct {
	*PageRes
	Values []*Repositories `json:"values"`
}

func (r *Repositories) ListForAccount(ro *RepositoriesOptions) (*RepositoriesRes, error) {
	urlStr := r.c.requestUrl("/repositories/%s", ro.Owner)
	if ro.Role != "" {
		urlStr += "?role=" + ro.Role
	}
	repos, err := r.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}
	return decodeRepositorys(repos)
}

func (r *Repositories) ListForTeam(ro *RepositoriesOptions) (*RepositoriesRes, error) {
	urlStr := r.c.requestUrl("/repositories/%s", ro.Owner)
	if ro.Role != "" {
		urlStr += "?role=" + ro.Role
	}
	repos, err := r.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}
	return decodeRepositorys(repos)
}

func (r *Repositories) ListPublic() (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/", "")
	repos, err := r.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}
	return decodeRepositorys(repos)
}

func decodeRepositorys(reposResponse interface{}) (*RepositoriesRes, error) {
	repoMap := reposResponse.(map[string]interface{})

	if repoMap["type"] == "error" {
		return nil, DecodeError(repoMap)
	}

	var repositorysRes = new(RepositoriesRes)
	err := mapstructure.Decode(repoMap, repositorysRes)
	if err != nil {
		return nil, err
	}

	return repositorysRes, nil
}
