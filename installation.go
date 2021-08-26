package lynx

import (
	"fmt"
	"net/http"
	"net/url"
)

type Installation struct {
	ID           int64    `json:"id"`
	ClientID     int64    `json:"client_id"`
	Name         string   `json:"name"`
	Timezone     string   `json:"timezone"`
	Capabilities []string `json:"capabilities"`
}

type InstallationRow struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	ClientID       int64   `json:"client_id"`
	Created        int64   `json:"created"`
	OrganizationID int64   `json:"organization_id"`
	Notes          string  `json:"notes"`
	Users          []int64 `json:"users"`
	Meta           Meta    `json:"meta"`
}

func (c *Client) GetInstallationRow(installationID int64) (*InstallationRow, error) {
	res := &InstallationRow{}
	path := fmt.Sprintf("api/v2/installation/%d", installationID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) UpdateInstallation(i *InstallationRow) (*InstallationRow, error) {
	res := &InstallationRow{}
	path := fmt.Sprintf("api/v2/installation/%d", i.ID)
	request := c.newRequest(http.MethodPut, path, requestBody(i))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) ListInstallations(filter map[string]string) ([]*InstallationRow, error) {
	res := make([]*InstallationRow, 0, 20)
	query := url.Values{}
	for k, v := range filter {
		query.Set(k, v)
	}
	request := c.newRequest(http.MethodGet, fmt.Sprintf("api/v2/installation?%s", query.Encode()), nil)
	if err := c.do(request, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetInstallations(assignedOnly bool) ([]*Installation, error) {
	res := make([]*Installation, 0, 20)
	query := url.Values{}
	query["assigned"] = []string{fmt.Sprintf("%v", assignedOnly)}
	path := fmt.Sprintf("%s?%s", "api/v2/installationinfo", query.Encode())
	request := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(request, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetInstallation(installationID int64) (*Installation, error) {
	res := make([]*Installation, 0, 20)
	request := c.newRequest(http.MethodGet, "api/v2/installationinfo?assigned=false", nil)
	if err := c.do(request, &res); err != nil {
		return nil, err
	}
	for index, installation := range res {
		if installation.ID == installationID {
			return res[index], nil
		}
	}
	return nil, Error{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}
}

func (c *Client) GetInstallationByClientID(clientID int64, assignedOnly bool) (*Installation, error) {
	res := &Installation{}
	query := url.Values{}
	query["assigned"] = []string{fmt.Sprintf("%v", assignedOnly)}
	path := fmt.Sprintf("api/v2/installationinfo/%d?%s", clientID, query.Encode())
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}
