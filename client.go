package lynx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
)

// Options is connection options for the client
type Options struct {
	Authenticator Authentication
	APIBase       string
	MqttOptions   *mqtt.ClientOptions
	HTTPClient    *http.Client
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
	if options.HTTPClient == nil {
		options.HTTPClient = &http.Client{
			Timeout: time.Second * 5,
		}
		options.HTTPClient.Transport = &http.Transport{
			TLSHandshakeTimeout: time.Second * 5,
		}
		if tmp, err := url.Parse(options.APIBase); err == nil && (tmp.Scheme == "h2c" || tmp.Scheme == "h2") {
			tr := http2.Transport{}
			if tmp.Scheme == "h2c" {
				tmp.Scheme = "http"
				tr.AllowHTTP = true
				tr.DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
					var d net.Dialer
					return d.Dial(network, addr)
				}
				tr.DialTLSContext = func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
					var d net.Dialer
					return d.DialContext(ctx, network, tmp.Host)
				}
			} else {
				tmp.Scheme = "https"
			}
			options.APIBase = tmp.String()
			options.HTTPClient.Transport = &tr
		}
	}
	return &Client{
		c:    options.HTTPClient,
		opt:  options,
		Mqtt: mq,
	}
}

func requestError(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		err := Error{
			Code: response.StatusCode,
		}
		if response.StatusCode == http.StatusRequestURITooLong {
			return err
		}
		if response.StatusCode != http.StatusNoContent {
			jsonError := json.NewDecoder(response.Body).Decode(&err)
			if jsonError != nil {
				return jsonError
			}
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
	if out != nil {
		if err := json.NewDecoder(response.Body).Decode(out); err != nil {
			return err
		}
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
