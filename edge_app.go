package lynx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type EdgeApp struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	Category         string     `json:"category"`
	Tags             []string   `json:"tags"`
	ShortDescription string     `json:"short_description"`
	Description      string     `json:"description"`
	Publisher        *Publisher `json:"publisher,omitempty"`
	Official         bool       `json:"official"`
	Public           bool       `json:"public"`
	SourceURL        string     `json:"source_url"`
	Created          int64      `json:"created"`
	Updated          int64      `json:"updated"`
}

type Publisher struct {
	ID   int64      `json:"id"`
	Name string     `json:"name"`
	Apps []*EdgeApp `json:"apps,omitempty"`
}

type EdgeAppVersion struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type EdgeAppConfig struct {
	ID             int64                  `json:"id"`
	AppID          int64                  `json:"app_id"`
	InstallationID int64                  `json:"installation_id"`
	Version        string                 `json:"version"`
	Config         map[string]interface{} `json:"config"`
	Name           string                 `json:"name"`
	Created        int64                  `json:"created"`
	Updated        int64                  `json:"updated"`
}

func (c *Client) GetEdgeApps() ([]*EdgeApp, error) {
	res := make([]*EdgeApp, 0, 10)
	path := "api/v2/edge/app"
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetEdgeApp(id int64) (*EdgeApp, error) {
	res := &EdgeApp{}
	path := fmt.Sprintf("api/v2/edge/app/%d", id)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) CreateEdgeApp(app *EdgeApp) (*EdgeApp, error) {
	res := &EdgeApp{}
	path := "api/v2/edge/app"
	req := c.newRequest(http.MethodPost, path, requestBody(app))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) DownloadEdgeApp(id int64, version string) ([]byte, error) {
	path := fmt.Sprintf("api/v2/edge/app/%d/download?version=%s", id, version)
	req := c.newRequest(http.MethodGet, path, nil)
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := requestError(resp); err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func (c *Client) GetEdgeAppVersions(appID int64, untagged bool) ([]*EdgeAppVersion, error) {
	res := make([]*EdgeAppVersion, 0, 10)
	path := fmt.Sprintf("api/v2/edge/app/%d/version?untagged=%t", appID, untagged)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) CreateEdgeAppVersion(appID int64, luaFile, jsonFile io.Reader) (string, error) {
	res := &struct {
		Hash string `json:"hash"`
	}{}
	path := fmt.Sprintf("api/v2/edge/app/%d/version", appID)
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	if luaFile != nil {
		fw, err := w.CreateFormFile("app_lua", "app.lua")
		if err != nil {
			return "", err
		}
		io.Copy(fw, luaFile)
	}
	if jsonFile != nil {
		fw, err := w.CreateFormFile("app_json", "app.json")
		if err != nil {
			return "", err
		}
		io.Copy(fw, jsonFile)
	}
	w.Close()
	req := c.newRequest(http.MethodPost, path, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if err := c.do(req, res); err != nil {
		return "", err
	}
	return res.Hash, nil
}

func (c *Client) NameEdgeAppVersion(appID int64, version *EdgeAppVersion) (*EdgeAppVersion, error) {
	res := &EdgeAppVersion{}
	path := fmt.Sprintf("api/v2/edge/app/%d/publish", appID)
	req := c.newRequest(http.MethodPost, path, requestBody(version))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetEdgeAppConfigOptions(appID int64, version string) (json.RawMessage, error) {
	path := fmt.Sprintf("api/v2/edge/app/%d/configure?version=%s", appID, version)
	req := c.newRequest(http.MethodGet, path, nil)
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	err = requestError(resp)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(bodyBytes), nil
}

func (c *Client) GetConfiguredEdgeApps(installationID int64) ([]*EdgeAppConfig, error) {
	res := make([]*EdgeAppConfig, 0, 10)
	path := fmt.Sprintf("api/v2/edge/app/configured/%d", installationID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) CreateEdgeAppInstance(config *EdgeAppConfig) (*EdgeAppConfig, error) {
	res := &EdgeAppConfig{}
	path := fmt.Sprintf("api/v2/edge/app/configured/%d", config.InstallationID)
	req := c.newRequest(http.MethodPost, path, requestBody(config))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetEdgeAppInstance(InstallationID, instanceID int64) (*EdgeAppConfig, error) {
	res := &EdgeAppConfig{}
	path := fmt.Sprintf("api/v2/edge/app/configured/%d/%d", InstallationID, instanceID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) UpdateEdgeAppInstance(config *EdgeAppConfig) (*EdgeAppConfig, error) {
	res := &EdgeAppConfig{}
	path := fmt.Sprintf("api/v2/edge/app/configured/%d/%d", config.InstallationID, config.ID)
	req := c.newRequest(http.MethodPut, path, requestBody(config))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) DeleteEdgeAppInstance(config *EdgeAppConfig) error {
	path := fmt.Sprintf("/api/v2/edge/app/configured/%d/%d", config.InstallationID, config.ID)
	req := c.newRequest(http.MethodDelete, path, nil)
	return c.do(req, nil)
}
