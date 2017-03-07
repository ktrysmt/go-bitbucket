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

type branchRestrictionsBody struct {
	Kind    string `json:"kind"`
	Pattern string `json:"pattern"`
	Links   struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
	Value  int `json:"value"`
	Id     int `json:"id"`
	Users  []branchRestrictionsBodyUser
	Groups []branchRestrictionsBodyGroup
}

type branchRestrictionsBodyGroup struct {
	Name  string
	Links struct {
		Self      struct{ href string }
		Html      struct{ href string }
		Full_slug string
		Members   int
		Slug      string
	}
}

type branchRestrictionsBodyUser struct {
	Username     string `json:"username"`
	Website      string
	Display_name string
	Uuid         string
	Created_on   string
	Links        struct {
		Self         struct{ href string }
		Repositories struct{ href string }
		Html         struct{ href string }
		Followers    struct{ href string }
		Avatar       struct{ href string }
		Following    struct{ href string }
	}
}

func (b *BranchRestrictions) buildBranchRestrictionsBody(bo *BranchRestrictionsOptions) string {

	body := branchRestrictionsBody{
		Kind:    bo.Kind,
		Pattern: bo.Pattern,
		Users: []branchRestrictionsBodyUser{
			{
				Username: bo.Users[0],
			},
		},
		Value: 0,
	}

	// if bo.Kind == "push" {
	// 	if n := len(bo.Users); n > 0 {
	// 		for i, user := range bo.Users {
	// 			body["users"].([]map[string]string)[i] = map[string]string{"username": user}
	// 		}
	// 	}
	// 	if n := len(bo.Groups); n > 0 {
	// 		i := 0
	// 		for username, slug := range bo.Groups {
	// 			body["groups"].([]map[string]interface{})[i] = map[string]interface{}{"slug": slug}
	// 			body["groups"].([]map[string]interface{})[i]["owner"] = map[string]string{"username": username}
	// 			i++
	// 		}
	// 	}
	// }

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	pp.Println(body)
	pp.Println(string(data))

	return string(data)
}
