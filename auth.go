package lynx

import (
	"fmt"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Authentication interface {
	SetHTTPAuth(r *http.Request)
	SetMQTTAuth(o *mqtt.ClientOptions)
}

type Basic struct {
	User     string
	Password string
}

func (b Basic) SetHTTPAuth(r *http.Request) {
	r.SetBasicAuth(b.User, b.Password)
}
func (b Basic) SetMQTTAuth(o *mqtt.ClientOptions) {
	o.SetUsername(b.User)
	o.SetPassword(b.Password)
}

type AuthAPIKey struct {
	Key string
}

func (a AuthAPIKey) SetHTTPAuth(r *http.Request) {
	r.Header.Set("X-API-Key", a.Key)
}

func (a AuthAPIKey) SetMQTTAuth(o *mqtt.ClientOptions) {
	o.SetUsername("apikey")
	o.SetPassword(a.Key)
}

type AuthBearer struct {
	Token string
}

func (a AuthBearer) SetHTTPAuth(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func (a AuthBearer) SetMQTTAuth(o *mqtt.ClientOptions) {
	o.SetUsername("bearer")
	o.SetPassword(a.Token)
}

type AuthNone struct{}

func (a AuthNone) SetHTTPAuth(r *http.Request) {
}

func (a AuthNone) SetMQTTAuth(o *mqtt.ClientOptions) {
}
