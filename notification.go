package lynx

import (
	"fmt"
	"net/http"
)

type NotificationMessage struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Text string `json:"text"`
}

type NotificationOutput struct {
	ID                           int64             `json:"id"`
	Name                         string            `json:"name"`
	InstallationID               int64             `json:"installation_id"`
	NotificationOutputExecutorID int64             `json:"notification_output_executor_id"`
	NotificationMessageID        int64             `json:"notification_message_id"`
	Config                       map[string]string `json:"config"`
}

type NotificationOutputExecutor struct {
	ID             int64             `json:"id"`
	Type           string            `json:"type"`
	Name           string            `json:"name"`
	OrganizationID int64             `json:"organization_id"`
	Config         map[string]string `json:"config"`
	Secret         string            `json:"secret"`
}

type NotificationExecutorPayload struct {
	Message        string            `json:"message"`
	OutputConfig   map[string]string `json:"output_config"`
	ExecutorConfig map[string]string `json:"executor_config"`
	Organization   Organization      `json:"organization"`
	Installation   Installation      `json:"installation"`
	Payload        map[string]any    `json:"payload"`
}

func (c *Client) GetNotificationMessages(installationID int64) ([]*NotificationMessage, error) {
	res := make([]*NotificationMessage, 0, 20)
	path := fmt.Sprintf("api/v2/notification/%d/message", installationID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetNotificationMessage(installationID, messageID int64) (*NotificationMessage, error) {
	msg := &NotificationMessage{}
	path := fmt.Sprintf("api/v2/notification/%d/message/%d", installationID, messageID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Client) CreateNotificationMessage(installationID int64, message *NotificationMessage) (*NotificationMessage, error) {
	msg := &NotificationMessage{}
	path := fmt.Sprintf("api/v2/notification/%d/message", installationID)
	req := c.newRequest(http.MethodPost, path, requestBody(message))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(req, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Client) UpdateNotificationMessage(installationID int64, message *NotificationMessage) (*NotificationMessage, error) {
	msg := &NotificationMessage{}
	path := fmt.Sprintf("api/v2/notification/%d/message/%d", installationID, message.ID)
	req := c.newRequest(http.MethodPut, path, requestBody(message))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(req, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Client) DeleteNotificationMessage(installationID int64, message *NotificationMessage) error {
	path := fmt.Sprintf("api/v2/notification/%d/message/%d", installationID, message.ID)
	req := c.newRequest(http.MethodDelete, path, nil)
	return c.do(req, nil)
}

func (c *Client) GetNotificationOutputs(installationID int64) ([]*NotificationOutput, error) {
	res := make([]*NotificationOutput, 0, 20)
	path := fmt.Sprintf("api/v2/notification/%d/output", installationID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetNotificationOutput(installationID, outputID int64) (*NotificationOutput, error) {
	o := &NotificationOutput{}
	path := fmt.Sprintf("api/v2/notification/%d/output/%d", installationID, outputID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (c *Client) CreateNotificationOutput(output *NotificationOutput) (*NotificationOutput, error) {
	o := &NotificationOutput{}
	path := fmt.Sprintf("api/v2/notification/%d/output", output.InstallationID)
	req := c.newRequest(http.MethodPost, path, requestBody(output))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(req, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (c *Client) UpdateNotificationOutput(output *NotificationOutput) (*NotificationOutput, error) {
	o := &NotificationOutput{}
	path := fmt.Sprintf("api/v2/notification/%d/output/%d", output.InstallationID, output.ID)
	req := c.newRequest(http.MethodPut, path, requestBody(output))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(req, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (c *Client) DeleteNotificationOutput(output *NotificationOutput) error {
	path := fmt.Sprintf("api/v2/notification/%d/output/%d", output.InstallationID, output.ID)
	req := c.newRequest(http.MethodDelete, path, nil)
	return c.do(req, nil)
}

func (c *Client) GetNotificationOutputExecutors(installationID int64) ([]*NotificationOutputExecutor, error) {
	executors := make([]*NotificationOutputExecutor, 0, 5)
	path := fmt.Sprintf("api/v2/notification/%d/executor", installationID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, &executors); err != nil {
		return nil, err
	}
	return executors, nil
}

func (c *Client) GetNotificationOutputExecutor(installationID, executorID int64) (*NotificationOutputExecutor, error) {
	ex := &NotificationOutputExecutor{}
	path := fmt.Sprintf("api/v2/notification/%d/executor/%d", installationID, executorID)
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, ex); err != nil {
		return nil, err
	}
	return ex, nil
}

func (c *Client) SendNotification(installationID, outputID int64, data interface{}) error {
	path := fmt.Sprintf("api/v2/notification/%d/output/%d/send", installationID, outputID)
	req := c.newRequest(http.MethodPost, path, requestBody(data))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return c.do(req, nil)
}
