package lynx

import (
	"fmt"
	"net/http"
	"net/url"
)

type Organization struct {
	ID            int64                `json:"id"`
	Name          string               `json:"name"`
	Address       Address              `json:"address"`
	Email         string               `json:"email"`
	Phone         string               `json:"phone"`
	ForceSMSLogin bool                 `json:"force_sms_login"`
	Parent        int64                `json:"parent"`
	Children      []*OrganizationChild `json:"children"`
	Notes         string               `json:"notes"`
	Meta          Meta                 `json:"meta"`
	ProtectedMeta Meta                 `json:"protected_meta"`
}

type OrganizationList []*Organization

func (ol OrganizationList) MapByID() map[int64]*Organization {
	res := make(map[int64]*Organization, len(ol))
	for i, o := range ol {
		res[o.ID] = ol[i]
	}
	return res
}

type OrganizationChild struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func (c *Client) ListOrganization(minimal bool, filter Filter) ([]*Organization, error) {
	res := make([]*Organization, 0, 5)
	filter["minimal"] = fmt.Sprintf("%t", minimal)
	query := filter.ToURLValues()
	req := c.newRequest(http.MethodGet, fmt.Sprintf("api/v2/organization?%s", query.Encode()), nil)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetOrganization(organizationID int64) (*Organization, error) {
	res := &Organization{}
	path := fmt.Sprintf("api/v2/organization/%d", organizationID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) CreateOrganization(org *Organization) (*Organization, error) {
	organization := &Organization{}
	req := c.newRequest(http.MethodPost, "api/v2/organization", requestBody(org))
	if err := c.do(req, organization); err != nil {
		return nil, err
	}
	return organization, nil
}

func (c *Client) UpdateOrganization(org *Organization) (*Organization, error) {
	organization := &Organization{}
	path := fmt.Sprintf("api/v2/organization/%d", org.ID)
	req := c.newRequest(http.MethodPut, path, requestBody(org))
	if err := c.do(req, organization); err != nil {
		return nil, err
	}
	return organization, nil
}

func (c *Client) DeleteOrganization(org *Organization, force bool) error {
	qs := ""
	if force {
		qs = "&force=true"
	}
	path := fmt.Sprintf("api/v2/organization/%d%s", org.ID, qs)
	req := c.newRequest(http.MethodDelete, path, nil)
	if err := c.do(req, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetOrganizationMeta(organizationID int64, key string) (*MetaObject, error) {
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/organization/%d/meta/%s", organizationID, key)
	request := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) CreateOrganizationMeta(organizationID int64, key string, meta MetaObject, silent bool) (*MetaObject, error) {
	query := url.Values{
		"silent": []string{fmt.Sprintf("%t", silent)},
	}
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/organization/%d/meta/%s?%s", organizationID, key, query.Encode())
	request := c.newRequest(http.MethodPost, path, requestBody(meta))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) UpdateOrganizationMeta(organizationID int64, key string, meta MetaObject, silent, createMissing bool) (*MetaObject, error) {
	query := url.Values{
		"silent":         []string{fmt.Sprintf("%t", silent)},
		"create_missing": []string{fmt.Sprintf("%t", createMissing)},
	}
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/organization/%d/meta/%s?%s", organizationID, key, query.Encode())
	request := c.newRequest(http.MethodPut, path, requestBody(meta))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) DeleteOrganizationMeta(organizationID int64, key string, silent bool) error {
	query := url.Values{
		"silent": []string{fmt.Sprintf("%t", silent)},
	}
	path := fmt.Sprintf("api/v2/organization/%d/meta/%s?%s", organizationID, key, query.Encode())
	request := c.newRequest(http.MethodDelete, path, nil)
	if err := c.do(request, nil); err != nil {
		return err
	}
	return nil
}
