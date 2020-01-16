package bitbucket

import (
	"encoding/json"
	"os"

	"github.com/k0kubun/pp"
)

type Ref struct {
	c *Client
}

func (r *Ref) Create(ro *RefsOptions) (*RepositoryBranch, error) {
	data := r.buildRefsBody(ro)
	urlStr := r.c.requestUrl("/repositories/%s/%s/refs/branches", ro.Owner, ro.RepoSlug)

	response, err := r.c.execute("POST", urlStr, data)
	if err != nil {
		return nil, err
	}

	resultMap, isMap := response.(map[string]interface{})
	var rb RepositoryBranch
	if isMap  {
		res, err := json.Marshal(resultMap)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(res, &rb)
		if err != nil {
			return nil, err
		}
	}
	return &rb, nil
} 

func (r *Ref) Delete(ro *RefsOptions) (interface{}, error) {
	urlStr := r.c.requestUrl("/repositories/%s/%s/refs/branches/%s", ro.Owner, ro.RepoSlug, ro.BranchName)
	return r.c.execute("DELETE", urlStr, "")
} 

func (r *Ref) Get(ro *RefsOptions) (*RepositoryBranch, error) {
	urlStr := r.c.requestUrl("/repositories/%s/%s/refs/branches/%s", ro.Owner, ro.RepoSlug, ro.BranchName)
	response, err := r.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}	
	resultMap, isMap := response.(map[string]interface{})
	var rb RepositoryBranch
	if isMap  {
		res, err := json.Marshal(resultMap)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(res, &rb)
		if err != nil {
			return nil, err
		}
	}
	return &rb, nil
} 

func (r *Ref) buildRefsBody(ro *RefsOptions) string {
	body := map[string]interface{}{}
	body["name"] = ro.BranchName
	if ro.TargetHash != "" {
		body["target"] = map[string]string{
			"hash": ro.TargetHash,
		}
	} else {
		body["target"] = map[string]string{
			"hash" : "default",
		}
	}

	return r.buildJsonBody(body)
}

func (r *Ref) buildJsonBody(body map[string]interface{}) string {

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	return string(data)
}

