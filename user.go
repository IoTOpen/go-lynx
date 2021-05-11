package lynx

import "net/http"

func (c *Client) Me() (*User, error) {
	user := &User{}
	path := "api/v2/user/me"
	request := c.newRequest(http.MethodGet, path, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, user); err != nil {
		return nil, err
	}
	return user, nil
}
