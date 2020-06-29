package lynx

import (
	"fmt"
	"net/http"
)

type Authentication interface {
	SetHTTPAuth(r *http.Request)
}

type Basic struct {
	User string
	Password string
}

func (b Basic) SetHTTPAuth(r *http.Request) {
	r.SetBasicAuth(b.User, b.Password)
}

type AuthApiKey struct {
	Key string
}

func (a AuthApiKey) SetHTTPAuth(r *http.Request) {
	r.Header.Set("X-API-Key", a.Key)
}

type AuthBearer struct {
	Token string
}

func (a AuthBearer) SetHTTPAuth(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

