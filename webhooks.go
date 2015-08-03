package bitbucket

import (
	"encoding/json"
	"github.com/k0kubun/pp"
	"os"
)

type Webhooks struct {
	c *Client
}

func (r *Webhooks) buildWebhooksBody(ro *WebhooksOptions) string {

	body := map[string]interface{}{}

	if ro.Description != "" {
		body["description"] = ro.Description
	}
	if ro.Url != "" {
		body["url"] = ro.Url
	}
	if ro.Active == true || ro.Active == false {
		body["active"] = ro.Active
	}

	if n := len(ro.Events); n > 0 {
		for i, event := range ro.Events {
			body["events"].([]string)[i] = event
		}
	}

	data, err := json.Marshal(body)
	if err != nil {
		pp.Println(err)
		os.Exit(9)
	}

	return string(data)
}

func (r *Webhooks) Gets(ro *WebhooksOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s/hooks/", ro.Owner, ro.Repo_slug)
	return r.c.execute("GET", url, "")
}

func (r *Webhooks) Create(ro *WebhooksOptions) interface{} {
	data := r.buildWebhooksBody(ro)
	url := r.c.requestUrl("/repositories/%s/%s/hooks", ro.Owner, ro.Repo_slug)
	return r.c.execute("POST", url, data)
}

func (r *Webhooks) Get(ro *WebhooksOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s/hooks/%s", ro.Owner, ro.Repo_slug, ro.Uuid)
	return r.c.execute("GET", url, "")
}

func (r *Webhooks) Update(ro *WebhooksOptions) interface{} {
	data := r.buildWebhooksBody(ro)
	url := r.c.requestUrl("/repositories/%s/%s/hooks/%s", ro.Owner, ro.Repo_slug, ro.Uuid)
	return r.c.execute("PUT", url, data)
}

func (r *Webhooks) Delete(ro *WebhooksOptions) interface{} {
	url := r.c.requestUrl("/repositories/%s/%s/hooks/%s", ro.Owner, ro.Repo_slug, ro.Uuid)
	return r.c.execute("DELETE", url, "")
}

//
