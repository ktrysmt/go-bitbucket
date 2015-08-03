package bitbucket

import (
	"encoding/json"
	"github.com/k0kubun/pp"
	"os"
)

type BranchRestrictions struct {
	c *Client
}

func (b *BranchRestrictions) Gets(bo *BranchRestrictionsOptions) interface{} {
	url := b.c.requestUrl("/repositories/%s/%s/branch-restrictions", bo.Owner, bo.Repo_slug)
	return b.c.execute("GET", url, "")
}

func (b *BranchRestrictions) Create(bo *BranchRestrictionsOptions) interface{} {
	data := b.buildBranchRestrictionsBody(bo)
	url := b.c.requestUrl("/repositories/%s/%s/branch-restrictions", bo.Owner, bo.Repo_slug)
	return b.c.execute("POST", url, data)
}

func (b *BranchRestrictions) Get(bo *BranchRestrictionsOptions) interface{} {
	url := b.c.requestUrl("/repositories/%s/%s/branch-restrictions/%s", bo.Owner, bo.Repo_slug, bo.Id)
	return b.c.execute("GET", url, "")
}

func (b *BranchRestrictions) Update(bo *BranchRestrictionsOptions) interface{} {
	data := b.buildBranchRestrictionsBody(bo)
	url := b.c.requestUrl("/repositories/%s/%s/branch-restrictions/%s", bo.Owner, bo.Repo_slug, bo.Id)
	return b.c.execute("PUT", url, data)
}

func (b *BranchRestrictions) Delete(bo *BranchRestrictionsOptions) interface{} {
	url := b.c.requestUrl("/repositories/%s/%s/branch-restrictions/%s", bo.Owner, bo.Repo_slug, bo.Id)
	return b.c.execute("DELETE", url, "")
}

func (b *BranchRestrictions) buildBranchRestrictionsBody(bo *BranchRestrictionsOptions) string {

	body := map[string]interface{}{}
	body["groups"] = map[string]string{}
	body["users"] = map[string]string{}

	if bo.Pattern != "" {
		body["pattern"] = bo.Pattern
	}
	if bo.Kind != "" {
		body["kind"] = bo.Kind
	}

	if bo.Kind == "push" {
		if n := len(bo.Users); n > 0 {
			for i, user := range bo.Users {
				body["users"].([]map[string]string)[i] = map[string]string{"username": user}
			}
		}
		if n := len(bo.Groups); n > 0 {
			i := 0
			for username, slug := range bo.Groups {
				body["groups"].([]map[string]interface{})[i] = map[string]interface{}{"slug": slug}
				body["groups"].([]map[string]interface{})[i]["owner"] = map[string]string{"username": username}
				i++
			}
		}
	}

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	return string(data)
}
