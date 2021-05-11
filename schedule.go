package lynx

import (
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetSchedules(installationID int64, executor string) ([]*Schedule, error) {
	res := make([]*Schedule, 0, 20)
	query := url.Values{}
	if executor != "" {
		query["executor"] = []string{executor}
	}
	path := fmt.Sprintf("api/v2/schedule/%d?%s", installationID, query.Encode())
	request := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(request, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetSchedule(installationID, scheduleID int64) (*Schedule, error) {
	schedule := &Schedule{}
	path := fmt.Sprintf("api/v2/schedule/%d/%d", installationID, scheduleID)
	request := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(request, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (c *Client) CreateSchedule(s *Schedule) (*Schedule, error) {
	schedule := &Schedule{}
	path := fmt.Sprintf("api/v2/schedule/%d", s.InstallationID)
	request := c.newRequest(http.MethodPost, path, requestBody(s))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (c *Client) UpdateSchedule(s *Schedule) (*Schedule, error) {
	schedule := &Schedule{}
	path := fmt.Sprintf("api/v2/schedule/%d/%d", s.InstallationID, s.ID)
	request := c.newRequest(http.MethodPut, path, requestBody(s))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if err := c.do(request, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (c *Client) DeleteSchedule(s *Schedule) error {
	path := fmt.Sprintf("api/v2/schedule/%d/%d", s.InstallationID, s.ID)
	request := c.newRequest(http.MethodDelete, path, nil)
	if err := c.do(request, nil); err != nil {
		return err
	}
	return nil
}
