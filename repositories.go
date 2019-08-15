package bitbucket

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
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
	Downloads          *Downloads
	repositories
}

type RepositoriesRes struct {
	Page    int32
	Pagelen int32
	Size    int32
	Items   []Repository
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
	repos, err := r.c.executeRaw("GET", urlStr, "")
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
	var reposResponseMap map[string]interface{}
	err := json.Unmarshal(reposResponse.([]byte), &reposResponseMap)
	if err != nil {
		return nil, err
	}

	repoArray := reposResponseMap["values"].([]interface{})
	var repos []Repository
	for _, repoEntry := range repoArray {
		var repo Repository
		err = mapstructure.Decode(repoEntry, &repo)
		if err == nil {
			repos = append(repos, repo)
		}
	}

	page, ok := reposResponseMap["page"].(float64)
	if !ok {
		page = 0
	}

	pagelen, ok := reposResponseMap["pagelen"].(float64)
	if !ok {
		pagelen = 0
	}
	size, ok := reposResponseMap["size"].(float64)
	if !ok {
		size = 0
	}

	repositories := RepositoriesRes{
		Page:    int32(page),
		Pagelen: int32(pagelen),
		Size:    int32(size),
		Items:   repos,
	}
	return &repositories, nil
}
