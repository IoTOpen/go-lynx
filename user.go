package lynx

import (
	"fmt"
	"net/http"
	"net/url"
)

const userMePath = "api/v2/user/me"

type User struct {
	ID            int64   `json:"id"`
	Email         string  `json:"email"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Role          int64   `json:"role"`
	SmsLogin      bool    `json:"sms_login"`
	Address       Address `json:"address"`
	Mobile        string  `json:"mobile"`
	Organizations []int64 `json:"organisations"`
	Meta          Meta    `json:"meta"`
	ProtectedMeta Meta    `json:"protected_meta"`
}

type Address struct {
	Address string `json:"address"`
	City    string `json:"city"`
	Country string `json:"country"`
	ZIP     string `json:"zip"`
}

func (c *Client) Me() (*User, error) {
	user := &User{}
	request := c.newRequest(http.MethodGet, userMePath, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *Client) UpdateMe(u *User) (*User, error) {
	user := &User{}
	request := c.newRequest(http.MethodPut, userMePath, requestBody(u))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *Client) UpdateUser(u *User) (*User, error) {
	user := &User{}
	path := fmt.Sprintf("api/v2/user/%d", u.ID)
	request := c.newRequest(http.MethodPut, path, requestBody(u))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *Client) GetUsers(filter Filter) ([]*User, error) {
	res := make([]*User, 0, 5)
	query := filter.ToURLValues()
	req := c.newRequest(http.MethodGet, fmt.Sprintf("api/v2/user?%s", query.Encode()), nil)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetUserMeta(userID int64, key string) (*MetaObject, error) {
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/user/%d/meta/%s", userID, key)
	request := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) CreateUserMeta(userID int64, key string, meta MetaObject, silent bool) (*MetaObject, error) {
	query := url.Values{
		"silent": []string{fmt.Sprintf("%t", silent)},
	}
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/user/%d/meta/%s?%s", userID, key, query.Encode())
	request := c.newRequest(http.MethodPost, path, requestBody(meta))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) UpdateUserMeta(userID int64, key string, meta MetaObject, silent, createMissing bool) (*MetaObject, error) {
	query := url.Values{
		"silent":         []string{fmt.Sprintf("%t", silent)},
		"create_missing": []string{fmt.Sprintf("%t", createMissing)},
	}
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/user/%d/meta/%s?%s", userID, key, query.Encode())
	request := c.newRequest(http.MethodPut, path, requestBody(meta))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) DeleteUserMeta(userID int64, key string, silent bool) error {
	query := url.Values{
		"silent": []string{fmt.Sprintf("%t", silent)},
	}
	path := fmt.Sprintf("api/v2/user/%d/meta/%s?%s", userID, key, query.Encode())
	request := c.newRequest(http.MethodDelete, path, nil)
	if err := c.do(request, nil); err != nil {
		return err
	}
	return nil
}
