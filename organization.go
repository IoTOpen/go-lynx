package lynx

import (
	"fmt"
	"net/http"
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
	Meta          map[string]string    `json:"meta"`
	ProtectedMeta map[string]string    `json:"protected_meta"`
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
