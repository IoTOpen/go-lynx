package lynx

import (
	"fmt"
	"net/http"
	"net/url"
)

type Device struct {
	ID             int64  `json:"id"`
	Type           string `json:"type"`
	InstallationID int64  `json:"installation_id"`
	Meta           Meta   `json:"meta"`
	ProtectedMeta  Meta   `json:"protected_meta"`
	Created        int64  `json:"created"`
	Updated        int64  `json:"updated"`
}

type DeviceList []*Device

func (d DeviceList) MapByID() map[int64]*Device {
	res := make(map[int64]*Device, len(d))
	for i, v := range d {
		res[v.ID] = d[i]
	}
	return res
}

func (d DeviceList) MapBy(key string) map[string]*Device {
	res := make(map[string]*Device, len(d))
	for i, v := range d {
		res[v.Meta[key]] = d[i]
	}
	return res
}

func (d DeviceList) MapByList(key string) map[string][]*Device {
	res := make(map[string][]*Device, len(d))
	for i, v := range d {
		arr, ok := res[v.Meta[key]]
		if !ok {
			arr = make([]*Device, 0, 10)
		}
		res[v.Meta[key]] = append(arr, d[i])
	}
	return res
}

func (c *Client) GetDevices(installationID int64, filter Filter) (DeviceList, error) {
	res := make(DeviceList, 0, 20)
	query := filter.ToURLValues()
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

func (c *Client) GetDeviceMeta(installationID, deviceID int64, key string) (*MetaObject, error) {
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/devicex/%d/%d/meta/%s", installationID, deviceID, key)
	request := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) CreateDeviceMeta(installationID, deviceID int64, key string, meta MetaObject, silent bool) (*MetaObject, error) {
	query := url.Values{
		"silent": []string{fmt.Sprintf("%t", silent)},
	}
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/devicex/%d/%d/meta/%s?%s", installationID, deviceID, key, query.Encode())
	request := c.newRequest(http.MethodPost, path, requestBody(meta))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) UpdateDeviceMeta(installationID, deviceID int64, key string, meta MetaObject, silent, createMissing bool) (*MetaObject, error) {
	query := url.Values{
		"silent":         []string{fmt.Sprintf("%t", silent)},
		"create_missing": []string{fmt.Sprintf("%t", createMissing)},
	}
	mo := &MetaObject{}
	path := fmt.Sprintf("api/v2/devicex/%d/%d/meta/%s?%s", installationID, deviceID, key, query.Encode())
	request := c.newRequest(http.MethodPut, path, requestBody(meta))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, mo); err != nil {
		return nil, err
	}
	return mo, nil
}

func (c *Client) DeleteDeviceMeta(installationID, deviceID int64, key string, silent bool) error {
	query := url.Values{
		"silent": []string{fmt.Sprintf("%t", silent)},
	}
	path := fmt.Sprintf("api/v2/devicex/%d/%d/meta/%s?%s", installationID, deviceID, key, query.Encode())
	request := c.newRequest(http.MethodDelete, path, nil)
	if err := c.do(request, nil); err != nil {
		return err
	}
	return nil
}
