package lynx

import (
	"bytes"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Options struct {
	Authenticator Authentication
	ApiBase       string
	MqttOptions   *mqtt.ClientOptions
}

type Client struct {
	opt  *Options
	c    *http.Client
	Mqtt mqtt.Client
}

func NewClient(options *Options) *Client {
	options.ApiBase = strings.TrimSuffix(options.ApiBase, "/")
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
	if err := json.NewDecoder(response.Body).Decode(out); err != nil {
		return err
	}
	return nil
}

func (c *Client) newRequest(method, path string, body io.Reader) *http.Request {
	uri := fmt.Sprintf("%s/%s", c.opt.ApiBase, path)
	r, _ := http.NewRequest(method, uri, body)
	c.opt.Authenticator.SetHTTPAuth(r)
	return r
}

func (c *Client) GetInstallations() ([]*Installation, error) {
	res := make([]*Installation, 0, 20)
	request := c.newRequest(http.MethodGet, "api/v2/installationinfo", nil)
	if err := c.do(request, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetFunctions(installationID int64, filter map[string]string) ([]*Function, error) {
	res := make([]*Function, 0, 20)

	parts := make([]string, 0, len(filter))
	for k, v := range filter {
		key, value := url.QueryEscape(k), url.QueryEscape(v)
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	query := strings.Join(parts, "&")
	request := c.newRequest(http.MethodGet, fmt.Sprintf("api/v2/functionx/%d?%s", installationID, query), nil)
	if err := c.do(request, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetFunction(installationID, functionID int64) (*Function, error) {
	function := &Function{}
	path := fmt.Sprintf("/api/v2/functionx/%d/%d", installationID, functionID)
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
	x, _ := httputil.DumpRequestOut(request, true)
	log.Println(string(x))
	if err := c.do(request, function); err != nil {
		return nil, err
	}
	return function, nil
}

func (c *Client) UpdateFunction(fn *Function) (*Function, error) {
	function := &Function{}
	path := fmt.Sprintf("api/v2/functionx/%d/%d", fn.InstallationID, fn.ID)
	request := c.newRequest(http.MethodPut, path, requestBody(fn))
	request.Header.Set("Content-Type", "application/json; utf-8")
	if err := c.do(request, function); err != nil {
		return nil, err
	}
	return function, nil
}

func (c *Client) GetDevices(installationID int64, filter map[string]string) ([]*Device, error) {
	res := make([]*Device, 0, 20)

	parts := make([]string, 0, len(filter))
	for k, v := range filter {
		key, value := url.QueryEscape(k), url.QueryEscape(v)
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	query := strings.Join(parts, "&")
	request := c.newRequest(http.MethodGet, fmt.Sprintf("api/v2/devicex/%d?%s", installationID, query), nil)
	if err := c.do(request, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetDevice(installationID, deviceID int64) (*Device, error) {
	device := &Device{}
	path := fmt.Sprintf("api/v2/devicex/%d/%d", installationID, deviceID)
	request := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(request, device); err != nil {
		return nil, err
	}
	return device, nil
}

func (c *Client) CreateDevice(dev *Device) (*Device, error) {
	device := &Device{}
	path := fmt.Sprintf("api/v2/devicex/%d", dev.InstallationID)
	request := c.newRequest(http.MethodPost, path, requestBody(dev))
	request.Header.Set("Content-Type", "application/json; utf-8")
	if err := c.do(request, device); err != nil {
		return nil, err
	}
	return device, nil
}

func (c *Client) UpdateDevice(dev *Device) (*Device, error) {
	device := &Device{}
	path := fmt.Sprintf("api/v2/devicex/%d/%d", dev.InstallationID, dev.ID)
	request := c.newRequest(http.MethodPut, path, requestBody(dev))
	request.Header.Set("Content-Type", "application/json; utf-8")
	if err := c.do(request, device); err != nil {
		return nil, err
	}
	return device, nil
}

func (c *Client) Status(installationID int64) (Status, error) {
	status := Status{}
	path := fmt.Sprintf("api/v2/status/%d", installationID)
	request := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(request, &status); err != nil {
		return nil, err
	}
	return status, nil
}

func (c *Client) Me() (*User, error) {
	user := &User{}
	path := "api/v2/user/me"
	request := c.newRequest(http.MethodGet, path, nil)
	request.Header.Set("Content-Type", "application/json; utf-8")
	if err := c.do(request, user); err != nil {
		return nil, err
	}
	return user, nil
}
