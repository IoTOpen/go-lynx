package lynx

import (
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetFunctions(installationID int64, filter map[string]string) ([]*Function, error) {
	res := make([]*Function, 0, 20)
	query := url.Values{}
	for k, v := range filter {
		query[k] = []string{v}
	}
	request := c.newRequest(http.MethodGet, fmt.Sprintf("api/v2/functionx/%d?%s", installationID, query.Encode()), nil)
	if err := c.do(request, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetFunction(installationID, functionID int64) (*Function, error) {
	function := &Function{}
	path := fmt.Sprintf("api/v2/functionx/%d/%d", installationID, functionID)
	request := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(request, function); err != nil {
		return nil, err
	}
	return function, nil
}

func (c *Client) CreateFunction(fn *Function) (*Function, error) {
	function := &Function{}
	path := fmt.Sprintf("api/v2/functionx/%d", fn.InstallationID)
	request := c.newRequest(http.MethodPost, path, requestBody(fn))
	if err := c.do(request, function); err != nil {
		return nil, err
	}
	return function, nil
}

func (c *Client) DeleteFunction(fn *Function) error {
	path := fmt.Sprintf("api/v2/functionx/%d/%d", fn.InstallationID, fn.ID)
	request := c.newRequest(http.MethodDelete, path, nil)
	if err := c.do(request, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateFunction(fn *Function) (*Function, error) {
	function := &Function{}
	path := fmt.Sprintf("api/v2/functionx/%d/%d", fn.InstallationID, fn.ID)
	request := c.newRequest(http.MethodPut, path, requestBody(fn))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, function); err != nil {
		return nil, err
	}
	return function, nil
}
