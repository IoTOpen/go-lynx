package lynx

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const maxTopicQueryLength = 30

type LogEntry struct {
	ClientID       int64   `json:"client_id"`
	InstallationID int64   `json:"installation_id"`
	Message        string  `json:"msg"`
	Timestamp      float64 `json:"timestamp"`
	Topic          string  `json:"topic"`
	Value          float64 `json:"value"`
}

type Status []*LogEntry

type V3Log struct {
	Total    int64      `json:"total"`
	LastTime float64    `json:"last"`
	Count    int        `json:"count"`
	Data     []LogEntry `json:"data"`
}

type LogOptionsV3 struct {
	Limit       int64
	Offset      int64
	From        time.Time
	To          time.Time
	Order       LogOrder
	TopicFilter []string
}

type LogOrder string

const (
	LogOrderDesc = LogOrder("desc")
	LogOrderAsc  = LogOrder("asc")
)

func (s Status) Map() map[string]*LogEntry {
	res := make(map[string]*LogEntry, len(s))
	for i, v := range s {
		res[v.Topic] = s[i]
	}
	return res
}
func (c *Client) Status(installationID int64, topicFilter []string) (Status, error) {
	status := Status{}
	query := url.Values{}
	method := http.MethodGet
	var body io.Reader
	if len(topicFilter) > maxTopicQueryLength {
		method = http.MethodPost
		body = requestBody(topicFilter)
	} else {
		query["topics"] = topicFilter
	}
	path := fmt.Sprintf("api/v2/status/%d?%s", installationID, query.Encode())
	request := c.newRequest(method, path, body)
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
	}
	method := http.MethodGet
	var body io.Reader
	if len(opts.TopicFilter) > maxTopicQueryLength {
		method = http.MethodPost
		body = requestBody(opts.TopicFilter)
	} else {
		query["topics"] = opts.TopicFilter
	}

	path := fmt.Sprintf("api/v3beta/log/%d?%s", installationID, query.Encode())
	request := c.c.newRequest(method, path, body)
	if err := c.c.do(request, log); err != nil {
		return nil, err
	}
	return log, nil
}
