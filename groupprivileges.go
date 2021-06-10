package bitbucket

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type GroupPrivileges struct {
	c *Client
}

type GroupPrivilegesOptions struct {
	Owner      string
	RepoSlug   string
	Group      string
	GroupOwner string
	Permission string
}

type GroupPrivilege struct {
	Repo      string `mapstructure:"repo"`
	Privilege string `mapstructure:"privilege"`
	Group     struct {
		Owner struct {
			DisplayName string `mapstructure:"display_name"`
			UUID        string `mapstructure:"uuid"`
			IsActive    bool   `mapstructure:"is_active"`
			IsTeam      bool   `mapstructure:"is_team"`
			MentionID   string `mapstructure:"mention_id"`
			Avatar      string `mapstructure:"avatar"`
			Nickname    string `mapstructure:"nickname"`
			AccountID   string `mapstructure:"account_id"`
		} `mapstructure:"owner"`
		Name    string        `mapstructure:"name"`
		Members []interface{} `mapstructure:"members"`
		Slug    string        `mapstructure:"slug"`
	} `mapstructure:"group"`
	Repository struct {
		Owner struct {
			DisplayName string `mapstructure:"display_name"`
			UUID        string `mapstructure:"uuid"`
			IsActive    bool   `mapstructure:"is_active"`
			IsTeam      bool   `mapstructure:"is_team"`
			MentionID   string `mapstructure:"mention_id"`
			Avatar      string `mapstructure:"avatar"`
			Nickname    string `mapstructure:"nickname"`
			AccountID   string `mapstructure:"account_id"`
		} `mapstructure:"owner"`
		Name string `mapstructure:"name"`
		Slug string `mapstructure:"slug"`
	} `mapstructure:"repository"`
}

func (g *GroupPrivileges) List(workspace, repoSlug string) ([]GroupPrivilege, error) {
	urlStr := fmt.Sprintf("%s/1.0/group-privileges/%s/%s", g.c.GetApiHostnameURL(), workspace, repoSlug)
	data, err := g.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return g.decodeGroupPrivileges(data)
}

func (g *GroupPrivileges) Add(gpo GroupPrivilegesOptions) ([]GroupPrivilege, error) {
	groupOwner := gpo.GroupOwner
	if gpo.GroupOwner == "" {
		groupOwner = gpo.Owner

	}
	urlStr := fmt.Sprintf("%s/1.0/group-privileges/%s/%s/%s/%s/", g.c.GetApiHostnameURL(), gpo.Owner, gpo.RepoSlug, groupOwner, gpo.Group)
	data, err := g.c.executeContentType("PUT", urlStr, gpo.Permission, "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}

	return g.decodeGroupPrivileges(data)
}

func (g *GroupPrivileges) Get(gpo GroupPrivilegesOptions) ([]GroupPrivilege, error) {
	groupOwner := gpo.GroupOwner
	if gpo.GroupOwner == "" {
		groupOwner = gpo.Owner

	}
	urlStr := fmt.Sprintf("%s/1.0/group-privileges/%s/%s/%s/%s/", g.c.GetApiHostnameURL(), gpo.Owner, gpo.RepoSlug, groupOwner, gpo.Group)
	data, err := g.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return g.decodeGroupPrivileges(data)
}

func (g *GroupPrivileges) Delete(gpo GroupPrivilegesOptions) (interface{}, error) {
	groupOwner := gpo.GroupOwner
	if gpo.GroupOwner == "" {
		groupOwner = gpo.Owner

	}
	urlStr := fmt.Sprintf("%s/1.0/group-privileges/%s/%s/%s/%s/", g.c.GetApiHostnameURL(), gpo.Owner, gpo.RepoSlug, groupOwner, gpo.Group)
	return g.c.execute("DELETE", urlStr, "")
}

func (g *GroupPrivileges) decodeGroupPrivileges(response interface{}) ([]GroupPrivilege, error) {
	var gp []GroupPrivilege
	err := mapstructure.Decode(response, &gp)
	if err != nil {
		return nil, err
	}
	return gp, nil
}
