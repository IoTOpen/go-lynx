package lynx

import "net/http"

type User struct {
	ID        int64   `json:"id"`
	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Role      int64   `json:"role"`
	SmsLogin  bool    `json:"sms_login"`
	Address   Address `json:"address"`
}

type Address struct {
	Address string `json:"address"`
	City    string `json:"city"`
	Country string `json:"country"`
	ZIP     string `json:"zip"`
}

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
