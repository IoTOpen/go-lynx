package lynx

import (
	"fmt"
	"net/http"
	"net/url"
)

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
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
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
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, device); err != nil {
		return nil, err
	}
	return device, nil
}
