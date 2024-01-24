package lynx

import (
	"fmt"
	"net/http"
	"net/url"
)

type Function struct {
	ID             int64  `json:"id"`
	Type           string `json:"type"`
	InstallationID int64  `json:"installation_id"`
	Meta           Meta   `json:"meta"`
	ProtectedMeta  Meta   `json:"protected_meta"`
	Created        int64  `json:"created"`
	Updated        int64  `json:"updated"`
}

type FunctionList []*Function

func (f FunctionList) MapByID() map[int64]*Function {
	res := make(map[int64]*Function, len(f))
	for i, v := range f {
		res[v.ID] = f[i]
	}
	return res
}

func (f FunctionList) MapBy(key string) map[string]*Function {
	res := make(map[string]*Function, len(f))
	for i, v := range f {
		res[v.Meta[key]] = f[i]
	}
	return res
}

func (f FunctionList) MapByList(key string) map[string][]*Function {
	res := make(map[string][]*Function, len(f))
	for i, v := range f {
		arr, ok := res[v.Meta[key]]
		if !ok {
			arr = make([]*Function, 0, 10)
		}
		res[v.Meta[key]] = append(arr, f[i])
	}
	return res
}

func (c *Client) GetFunctions(installationID int64, filter Filter) (FunctionList, error) {
	res := make([]*Function, 0, 20)
	query := filter.ToURLValues()
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

func (c *Client) GetFunctionMeta(installationID, functionID int64, key string) (*MetaObject, error) {
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/functionx/%d/%d/meta/%s", installationID, functionID, key)
	request := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) CreateFunctionMeta(installationID, functionID int64, key string, meta MetaObject, silent bool) (*MetaObject, error) {
	query := url.Values{
		"silent": []string{fmt.Sprintf("%t", silent)},
	}
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/functionx/%d/%d/meta/%s?%s", installationID, functionID, key, query.Encode())
	request := c.newRequest(http.MethodPost, path, requestBody(meta))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) UpdateFunctionMeta(installationID, functionID int64, key string, meta MetaObject, silent, createMissing bool) (*MetaObject, error) {
	query := url.Values{
		"silent":         []string{fmt.Sprintf("%t", silent)},
		"create_missing": []string{fmt.Sprintf("%t", createMissing)},
	}
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/functionx/%d/%d/meta/%s?%s", installationID, functionID, key, query.Encode())
	request := c.newRequest(http.MethodPut, path, requestBody(meta))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) DeleteFunctionMeta(installationID, functionID int64, key string, silent bool) error {
	query := url.Values{
		"silent": []string{fmt.Sprintf("%t", silent)},
	}
	path := fmt.Sprintf("api/v2/functionx/%d/%d/meta/%s?%s", installationID, functionID, key, query.Encode())
	request := c.newRequest(http.MethodDelete, path, nil)
	if err := c.do(request, nil); err != nil {
		return err
	}
	return nil
}
