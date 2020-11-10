package bitbucket

import (
	"github.com/mitchellh/mapstructure"
)

type Workspace struct {
	c *Client

	Repositories *Repositories

	UUID       string
	Type       string
	Slug       string
	Is_Private bool
	Name       string
}

type WorkspaceList struct {
	Page       int
	Pagelen    int
	MaxDepth   int
	Size       int
	Next       string
	Workspaces []Workspace
}

func (t *Workspace) List() (*WorkspaceList, error) {
	urlStr := t.c.requestUrl("/workspaces")
	response, err := t.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodeWorkspaceList(response)
}

func (t *Workspace) Get(workspace string) (*Workspace, error) {
	urlStr := t.c.requestUrl("/workspaces/%s", workspace)
	response, err := t.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodeWorkspace(response)
}

func decodeWorkspace(workspace interface{}) (*Workspace, error) {
	var workspaceEntry Workspace
	workspaceResponseMap := workspace.(map[string]interface{})

	if workspaceResponseMap["type"] != nil && workspaceResponseMap["type"] == "error" {
		return nil, DecodeError(workspaceResponseMap)
	}

	err := mapstructure.Decode(workspace, &workspaceEntry)
	return &workspaceEntry, err
}

func decodeWorkspaceList(workspaceResponse interface{}) (*WorkspaceList, error) {
	workspaceResponseMap := workspaceResponse.(map[string]interface{})
	workspaceMapList := workspaceResponseMap["values"].([]interface{})

	var workspaces []Workspace
	for _, workspaceMap := range workspaceMapList {
		workspaceEntry, err := decodeWorkspace(workspaceMap)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, *workspaceEntry)
	}

	page, ok := workspaceResponseMap["page"].(float64)
	if !ok {
		page = 0
	}

	pagelen, ok := workspaceResponseMap["pagelen"].(float64)
	if !ok {
		pagelen = 0
	}
	max_depth, ok := workspaceResponseMap["max_depth"].(float64)
	if !ok {
		max_depth = 0
	}
	size, ok := workspaceResponseMap["size"].(float64)
	if !ok {
		size = 0
	}

	next, ok := workspaceResponseMap["next"].(string)
	if !ok {
		next = ""
	}

	workspacesList := WorkspaceList{
		Page:       int(page),
		Pagelen:    int(pagelen),
		MaxDepth:   int(max_depth),
		Size:       int(size),
		Next:       next,
		Workspaces: workspaces,
	}

	return &workspacesList, nil
}
