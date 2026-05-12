package lynx

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"
)

type LogEntry struct {
	ClientID       int64   `json:"client_id"`
	InstallationID int64   `json:"installation_id"`
	Message        string  `json:"msg"`
	Timestamp      float64 `json:"timestamp"`
	Topic          string  `json:"topic"`
	Value          float64 `json:"value"`
}

func (l *LogEntry) Time() time.Time {
	whole, fractals := math.Modf(l.Timestamp)
	return time.Unix(int64(whole), int64(fractals*1000000000))
}

type Status []*LogEntry

type V3Log struct {
	Total    int64      `json:"total"`
	LastTime float64    `json:"last"`
	Count    int        `json:"count"`
	Data     []LogEntry `json:"data"`
}

type LogOptionsV3 struct {
	Limit        int64
	Offset       int64
	From         time.Time
	To           time.Time
	Order        LogOrder
	TopicFilter  []string
	AggrMethod   string
	AggrInterval time.Duration
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
	query := url.Values{
		"topics": topicFilter,
	}
	path := fmt.Sprintf("api/v2/status/%d?%s", installationID, query.Encode())
	req := c.newRequest(http.MethodGet, path, nil)
	err := c.do(req, &status)
	getErr := Error{}
	if errors.As(err, &getErr) && getErr.Code == http.StatusRequestURITooLong {
		path = fmt.Sprintf("api/v2/status/%d", installationID)
		body := requestBody(topicFilter)
		req = c.newRequest(http.MethodPost, path, body)
		if postErr := c.do(req, &status); postErr != nil {
			newErr := Error{}
			ok := errors.As(postErr, &newErr)
			if ok && newErr.Code == http.StatusMethodNotAllowed {
				return nil, err
			}
			return nil, postErr
		}
		return status, nil
	}

	return status, err
}

const (
	defaultV3LogLimit = int64(500)
	defaultV3LogRange = 24 * time.Hour
)

func normalizeLogOptionsV3(opts *LogOptionsV3) (*LogOptionsV3, error) {
	now := time.Now()

	if opts == nil {
		opts = &LogOptionsV3{}
	}

	normalized := *opts

	if normalized.To.IsZero() {
		normalized.To = now
	}
	if normalized.From.IsZero() {
		normalized.From = normalized.To.Add(-defaultV3LogRange)
	}
	if normalized.To.Before(normalized.From) {
		return nil, fmt.Errorf("invalid log time range: from %s is after to %s",
			normalized.From.Format(time.RFC3339),
			normalized.To.Format(time.RFC3339))
	}

	if normalized.Limit <= 0 {
		normalized.Limit = defaultV3LogLimit
	}
	if normalized.Offset < 0 {
		normalized.Offset = 0
	}

	switch normalized.Order {
	case "":
		normalized.Order = LogOrderDesc
	case LogOrderAsc, LogOrderDesc:
	default:
		return nil, fmt.Errorf("invalid log order: %q", normalized.Order)
	}

	if normalized.TopicFilter == nil {
		normalized.TopicFilter = []string{}
	}
	if normalized.AggrInterval < 0 {
		return nil, fmt.Errorf("invalid negative aggregation interval: %s", normalized.AggrInterval)
	}

	return &normalized, nil
}

// Log returns log entries in the V3 format. If opts is nil some default values will be used.
func (c *V3Client) Log(installationID int64, opts *LogOptionsV3) (*V3Log, error) {
	log := &V3Log{}

	normalized, err := normalizeLogOptionsV3(opts)
	if err != nil {
		return nil, err
	}

	query := url.Values{
		"from":   []string{fmt.Sprintf("%d", normalized.From.Unix())},
		"to":     []string{fmt.Sprintf("%d", normalized.To.Unix())},
		"limit":  []string{fmt.Sprintf("%d", normalized.Limit)},
		"offset": []string{fmt.Sprintf("%d", normalized.Offset)},
		"order":  []string{string(normalized.Order)},
		"topics": normalized.TopicFilter,
	}

	if normalized.AggrMethod != "" {
		query["aggr_method"] = []string{normalized.AggrMethod}
		if normalized.AggrInterval > 0 {
			query["aggr_interval"] = []string{normalized.AggrInterval.String()}
		}
	}

	path := fmt.Sprintf("api/v3beta/log/%d?%s", installationID, query.Encode())
	req := c.c.newRequest(http.MethodGet, path, nil)
	err = c.c.do(req, log)

	getErr := Error{}
	if errors.As(err, &getErr) && getErr.Code == http.StatusRequestURITooLong {
		delete(query, "topics")
		path = fmt.Sprintf("api/v3beta/log/%d?%s", installationID, query.Encode())
		body := requestBody(normalized.TopicFilter)
		req = c.c.newRequest(http.MethodPost, path, body)
		if postErr := c.c.do(req, log); postErr != nil {
			newErr := Error{}
			ok := errors.As(postErr, &newErr)
			if ok && newErr.Code == http.StatusMethodNotAllowed {
				return nil, err
			}
			return nil, postErr
		}
		return log, nil
	}

	return log, err
}
