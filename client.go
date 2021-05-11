package lynx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Options is connection options for the client
type Options struct {
	Authenticator Authentication
	APIBase       string
	MqttOptions   *mqtt.ClientOptions
}

// Client is the main client for Lynx integration
type Client struct {
	opt  *Options
	c    *http.Client
	Mqtt mqtt.Client
}

// V3Client is a client implementing the V3 endpoints
type V3Client struct {
	c *Client
}

// V3 Returns the V3 client
func (c *Client) V3() *V3Client {
	return &V3Client{c: c}
}

// NewClient create a new client for V2 API:s with specified options
func NewClient(options *Options) *Client {
	options.APIBase = strings.TrimSuffix(options.APIBase, "/")
	var mq mqtt.Client
	if options.MqttOptions != nil {
		options.Authenticator.SetMQTTAuth(options.MqttOptions)
		mq = mqtt.NewClient(options.MqttOptions)
	}
	return &Client{
		c: &http.Client{
			Timeout: time.Second * 5,
		},
		opt:  options,
		Mqtt: mq,
	}
}

func requestError(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		err := Error{
			Code: response.StatusCode,
		}
		jsonError := json.NewDecoder(response.Body).Decode(&err)
		if jsonError != nil {
			return jsonError
		}
		return err
	}
	return nil
}

func requestBody(data interface{}) io.Reader {
	bin, _ := json.Marshal(data)
	body := bytes.NewReader(bin)
	return body
}

func (c *Client) do(r *http.Request, out interface{}) error {
	response, err := c.c.Do(r)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if err := requestError(response); err != nil {
		return err
	}
	if response.StatusCode == http.StatusOK {
		if out != nil {
			if err := json.NewDecoder(response.Body).Decode(out); err != nil {
				return err
			}
		}
	} else {
		err := Error{Message: ""}
		if err := json.NewDecoder(response.Body).Decode(&err); err != nil {
			return err
		}
		err.Code = response.StatusCode
		return err
	}
	return nil
}

func (c *Client) newRequest(method, path string, body io.Reader) *http.Request {
	uri := fmt.Sprintf("%s/%s", c.opt.APIBase, path)
	r, _ := http.NewRequest(method, uri, body)
	c.opt.Authenticator.SetHTTPAuth(r)
	return r
}

// Ping verify that the api is responding
func (c *Client) Ping() error {
	request := c.newRequest(http.MethodGet, "/api/v2/ping", nil)
	if err := c.do(request, nil); err != nil {
		return err
	}
	return nil
}
