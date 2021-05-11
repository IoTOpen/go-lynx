package lynx

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

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
