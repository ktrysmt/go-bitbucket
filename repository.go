package bitbucket

import (
	"encoding/json"
	"os"

	"github.com/k0kubun/pp"
)

type Repository struct {
	c *Client
}

func (r *Repository) Create(ro *RepositoryOptions) interface{} {
	data := r.buildRepositoryBody(ro)
	url := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.RepoSlug)
	return r.c.execute("POST", url, data)
}

func (r *Repository) Get(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.RepoSlug)
	return r.c.execute("GET", url, "")
}

func (r *Repository) Delete(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.RepoSlug)
	return r.c.execute("DELETE", url, "")
}

func (r *Repository) ListWatchers(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s/watchers", ro.Owner, ro.RepoSlug)
	return r.c.execute("GET", url, "")
}

func (r *Repository) ListForks(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s/forks", ro.Owner, ro.RepoSlug)
	return r.c.execute("GET", url, "")
}

func (r *Repository) buildRepositoryBody(ro *RepositoryOptions) string {

	body := map[string]interface{}{}

	if ro.Scm != "" {
		body["scm"] = ro.Scm
	}
	//if ro.Scm != "" {
	//		body["name"] = ro.Name
	//}
	if ro.IsPrivate != "" {
		body["is_private"] = ro.IsPrivate
	}
	if ro.Description != "" {
		body["description"] = ro.Description
	}
	if ro.ForkPolicy != "" {
		body["fork_policy"] = ro.ForkPolicy
	}
	if ro.Language != "" {
		body["language"] = ro.Language
	}
	if ro.HasIssues != "" {
		body["has_issues"] = ro.HasIssues
	}
	if ro.HasWiki != "" {
		body["has_wiki"] = ro.HasWiki
	}

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	return string(data)
}
