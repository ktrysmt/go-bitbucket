package bitbucket

import (
	"encoding/json"
	"github.com/k0kubun/pp"
	"os"
)

type Repository struct {
	c *Client
}

func (r *Repository) Create(ro *RepositoryOptions) interface{} {
	data := r.buildRepositoryBody(ro)
	url := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.Repo_slug)
	return r.c.execute("POST", url, data)
}

func (r *Repository) Get(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.Repo_slug)
	return r.c.execute("GET", url, "")
}

func (r *Repository) Delete(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s", ro.Owner, ro.Repo_slug)
	return r.c.execute("DELETE", url, "")
}

func (r *Repository) ListWatchers(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s/watchers", ro.Owner, ro.Repo_slug)
	return r.c.execute("GET", url, "")
}

func (r *Repository) ListForks(ro *RepositoryOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s/forks", ro.Owner, ro.Repo_slug)
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
	if ro.Is_private != "" {
		body["is_private"] = ro.Is_private
	}
	if ro.Description != "" {
		body["description"] = ro.Description
	}
	if ro.Fork_policy != "" {
		body["fork_policy"] = ro.Fork_policy
	}
	if ro.Language != "" {
		body["language"] = ro.Language
	}
	if ro.Has_issues != "" {
		body["has_issues"] = ro.Has_issues
	}
	if ro.Has_wiki != "" {
		body["has_wiki"] = ro.Has_wiki
	}

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	return string(data)
}
