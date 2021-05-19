package lynx

import (
	"encoding/json"
	"fmt"
	"io"
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
	a := &EdgeApp{}
	path := fmt.Sprintf("api/v2/edge/app/%d", id)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (c *Client) CreateEdgeApp(app *EdgeApp) (*EdgeApp, error) {
	a := &EdgeApp{}
	path := "api/v2/edge/app"
	req := c.newRequest(http.MethodPost, path, requestBody(app))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(req, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (c *Client) DownloadEdgeApp(id int64, version string) (*io.Reader, error) {
	path := fmt.Sprintf("api/edge/app/%d/download?version=%s", id, version)
	req := c.newRequest(http.MethodGet, path, nil)
	//TODO: Download file as stream?
}

func (c *Client) GetEdgeAppVersions(appID int64, untagged bool) ([]*EdgeAppVersion, error) {
	versions := make([]*EdgeAppVersion, 0, 10)
	path := fmt.Sprintf("api/edge/app/%d/version?untagged=%t", appID, untagged)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

func (c *Client) CreateEdgeAppVersion() {
	// TODO: Files?
}

func (c *Client) NameEdgeAppVersion(appID int64, version *EdgeAppVersion) (*EdgeAppVersion, error) {
	ver := &EdgeAppVersion{}
	path := fmt.Sprintf("api/edge/app/%d/publish", appID)
	req := c.newRequest(http.MethodPost, path, requestBody(version))
	if err := c.do(req, ver); err != nil {
		return nil, err
	}
	return ver, nil
}

func (c *Client) GetEdgeAppConfigOptions(appID int64) (json.RawMessage, error) {
	path := fmt.Sprintf("api/edge/app/%d/configure", appID)
	req := c.newRequest(http.MethodGet, path, nil)
	//TODO: Do the request
	return nil, nil
}

func (c *Client) GetConfiguredEdgeApps(installationID int64) ([]*EdgeAppConfig, error) {
	res := make([]*EdgeAppConfig, 0, 10)
	path := fmt.Sprintf("api/edge/app/configured/%d", installationID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) CreateEdgeAppInstance(config *EdgeAppConfig) (*EdgeAppConfig, error) {
	con := &EdgeAppConfig{}
	path := fmt.Sprintf("api/edge/app/configured/%d", config.InstallationID)
	req := c.newRequest(http.MethodPost, path, requestBody(con))
	if err := c.do(req, con); err != nil {
		return nil, err
	}
	return con, nil
}
