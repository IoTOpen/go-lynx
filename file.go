package lynx

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
)

type File struct {
	ID             int64  `json:"id"`
	Hash           string `json:"hash"`
	Name           string `json:"name"`
	MIME           string `json:"mime"`
	InstallationID int64  `json:"installation_id"`
	OrganizationID int64  `json:"organization_id"`
	Created        int64  `json:"created"`
	Updated        int64  `json:"updated"`
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func (c *Client) GetFilesInstallation(installationID int64) ([]*File, error) {
	res := make([]*File, 0, 10)
	path := fmt.Sprintf("api/v2/file/installation/%d", installationID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetFileInstallation(installationID int64, fileID int64) (*File, error) {
	res := &File{}
	path := fmt.Sprintf("api/v2/file/installation/%d/%d", installationID, fileID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetFilesOrganization(organizationID int64) ([]*File, error) {
	res := make([]*File, 0, 10)
	path := fmt.Sprintf("api/v2/file/organization/%d", organizationID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetFileOrganization(organizationID int64, fileID int64) (*File, error) {
	res := &File{}
	path := fmt.Sprintf("api/v2/file/organization/%d/%d", organizationID, fileID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) CreateFileInstallation(installationID int64, fileName, mime string, r io.Reader) (*File, error) {
	path := fmt.Sprintf("api/v2/file/installation/%d", installationID)
	res := make([]*File, 0, 1)
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	headers := make(textproto.MIMEHeader)
	headers.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			quoteEscaper.Replace(fileName), quoteEscaper.Replace(filepath.Base(fileName))))
	headers.Set("Content-Type", mime)
	part, err := w.CreatePart(headers)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, r)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	req := c.newRequest(http.MethodPost, path, buf)
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", w.Boundary()))
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	if len(res) > 0 {
		return res[0], nil
	}
	return nil, fmt.Errorf("no response from server")
}

func (c *Client) CreateFileOrganization(organizationID int64, fileName, mime string, r io.Reader) (*File, error) {
	path := fmt.Sprintf("api/v2/file/organization/%d", organizationID)
	res := make([]*File, 0, 1)
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	headers := make(textproto.MIMEHeader)
	headers.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			quoteEscaper.Replace(fileName), quoteEscaper.Replace(filepath.Base(fileName))))
	headers.Set("Content-Type", mime)
	part, err := w.CreatePart(headers)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, r)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	req := c.newRequest(http.MethodPost, path, buf)
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", w.Boundary()))
	if err := c.do(req, &res); err != nil {
		return nil, err
	}

	if len(res) > 0 {
		return res[0], nil
	}
	return nil, fmt.Errorf("no response from server")
}

func (c *Client) UpdateFileInstallation(installationID, fileID int64, fileName, mime string, r io.Reader) (*File, error) {
	path := fmt.Sprintf("api/v2/file/installation/%d/%d", installationID, fileID)
	res := &File{}
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	headers := make(textproto.MIMEHeader)
	headers.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			quoteEscaper.Replace(fileName), quoteEscaper.Replace(filepath.Base(fileName))))
	headers.Set("Content-Type", mime)
	part, err := w.CreatePart(headers)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, r)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	req := c.newRequest(http.MethodPost, path, buf)
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", w.Boundary()))
	if err := c.do(req, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) UpdateFileOrganization(organizationID, fileID int64, fileName, mime string, r io.Reader) (*File, error) {
	path := fmt.Sprintf("api/v2/file/organization/%d/%d", organizationID, fileID)
	res := &File{}
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	headers := make(textproto.MIMEHeader)
	headers.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			quoteEscaper.Replace(fileName), quoteEscaper.Replace(filepath.Base(fileName))))
	headers.Set("Content-Type", mime)
	part, err := w.CreatePart(headers)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, r)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	req := c.newRequest(http.MethodPost, path, buf)
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", w.Boundary()))
	if err := c.do(req, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) DeleteFileInstallation(installationID, fileID int64) error {
	path := fmt.Sprintf("api/v2/file/installation/%d/%d", installationID, fileID)
	req := c.newRequest(http.MethodDelete, path, nil)
	if err := c.do(req, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteFileOrganization(organizationID, fileID int64) error {
	path := fmt.Sprintf("api/v2/file/organization/%d/%d", organizationID, fileID)
	req := c.newRequest(http.MethodDelete, path, nil)
	if err := c.do(req, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) DownloadFile(hash string) (io.ReadCloser, error) {
	path := fmt.Sprintf("api/v2/file/download/%s", hash)
	req := c.newRequest(http.MethodGet, path, nil)
	response, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	if err := requestError(response); err != nil {
		return nil, err
	}

	return response.Body, nil
}
