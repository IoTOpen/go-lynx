package lynx

import (
	"bytes"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io"
	"net/http"
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

type V3Client struct {
	c *Client
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
	uri := fmt.Sprintf("%s/%s", c.opt.ApiBase, path)
	r, _ := http.NewRequest(method, uri, body)
	c.opt.Authenticator.SetHTTPAuth(r)
	return r
}

func (c *Client) GetInstallations() ([]*Installation, error) {
	res := make([]*Installation, 0, 20)
	request := c.newRequest(http.MethodGet, "/api/v2/installationinfo", nil)
	if err := c.do(request, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetInstallation(installationID int64) (*Installation, error) {
	res := make([]*Installation, 0, 20)
	request := c.newRequest(http.MethodGet, "/api/v2/installationinfo", nil)
	if err := c.do(request, &res); err != nil {
		return nil, err
	}
	for index, installation := range res {
		if installation.ID == installationID {
			return res[index], nil
		}
	}
	return nil, Error{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}
}

func (c *Client) GetInstallationByClientID(clientID int64) (*Installation, error) {
	res := &Installation{}
	path := fmt.Sprintf("api/v2/installationinfo/%d", clientID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetFunctions(installationID int64, filter map[string]string) ([]*Function, error) {
	res := make([]*Function, 0, 20)
	query := url.Values{}
	for k, v := range filter {
		query[k] = []string{v}
	}
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
	request.Header.Set("Content-Type", "application/json; utf-8")
	if err := c.do(request, function); err != nil {
		return nil, err
	}
	return function, nil
}

func (c *Client) GetDevices(installationID int64, filter map[string]string) ([]*Device, error) {
	res := make([]*Device, 0, 20)
	query := url.Values{}
	for k, v := range filter {
		query[k] = []string{v}
	}
	request := c.newRequest(http.MethodGet, fmt.Sprintf("api/v2/devicex/%d?%s", installationID, query.Encode()), nil)
	if err := c.do(request, &res); err != nil {
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

func (c *Client) DeleteDevice(dev *Device) error {
	path := fmt.Sprintf("api/v2/devicex/%d/%d", dev.InstallationID, dev.ID)
	request := c.newRequest(http.MethodDelete, path, nil)
	if err := c.do(request, nil); err != nil {
		return err
	}
	return nil
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

func (c *Client) Status(installationID int64, topicFilter []string) (Status, error) {
	status := Status{}
	query := url.Values{
		"topics": topicFilter,
	}
	path := fmt.Sprintf("api/v2/status/%d?%s", installationID, query.Encode())
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

func (c *Client) V3() *V3Client {
	return &V3Client{c: c}
}

// Log returns log entries in the V3 format. If opts is nil some default values will be used.
func (c *V3Client) Log(installationID int64, opts *LogOptionsV3) (*V3Log, error) {
	log := &V3Log{}
	if opts == nil {
		t := time.Now()
		opts = &LogOptionsV3{
			From:        t.Add(-time.Hour * 24),
			To:          t,
			Limit:       500,
			Offset:      0,
			Order:       LogOrderDesc,
			TopicFilter: []string{},
		}
	}
	query := url.Values{
		"from":   []string{fmt.Sprintf("%d", opts.From.Unix())},
		"to":     []string{fmt.Sprintf("%d", opts.To.Unix())},
		"limit":  []string{fmt.Sprintf("%d", opts.Limit)},
		"offset": []string{fmt.Sprintf("%d", opts.Offset)},
		"order":  []string{string(opts.Order)},
		"topics": opts.TopicFilter,
	}

	path := fmt.Sprintf("api/v3beta/log/%d?%s", installationID, query.Encode())
	request := c.c.newRequest(http.MethodGet, path, nil)
	if err := c.c.do(request, log); err != nil {
		return nil, err
	}
	return log, nil
}
