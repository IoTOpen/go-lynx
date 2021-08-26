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

func (c *Client) GetUsers(filter map[string]string) ([]*User, error) {
	res := make([]*User, 0, 5)
	query := url.Values{}
	for k, v := range filter {
		query[k] = []string{v}
	}
	req := c.newRequest(http.MethodGet, fmt.Sprintf("api/v2/user?%s", query.Encode()), nil)

	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}
