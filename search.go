package lynx

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type SearchEntry struct {
	Type           string `json:"type"`
	ID             int64  `json:"id"`
	InstallationID int64  `json:"installation_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Created        int64  `json:"created"`
	Updated        int64  `json:"updated"`
	MetaMatches    *Meta  `json:"meta_matches"`
	Meta           Meta   `json:"metadata"`
}

type SearchResult struct {
	Total   int64         `json:"total"`
	Limit   int64         `json:"limit"`
	Offset  int64         `json:"offset"`
	Results []SearchEntry `json:"results"`
}

type SearchOptions struct {
	Query    string
	Types    []string
	Limit    int
	Offset   int
	Metadata Meta
}

func (c *V3Client) Search(opts *SearchOptions) (*SearchResult, error) {
	query := url.Values{}

	if opts.Query != "" {
		query.Set("q", opts.Query)
	}

	if opts.Types != nil {
		query.Set("types", strings.Join(opts.Types, ","))
	}

	if opts.Limit != 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}

	if opts.Offset != 0 {
		query.Set("offset", strconv.Itoa(opts.Offset))
	}

	if opts.Metadata != nil {
		for k, v := range opts.Metadata {
			query.Set(fmt.Sprintf("metadata.%s", k), v)
		}
	}

	path := fmt.Sprintf("api/v3beta/search?%s", query.Encode())
	req := c.c.newRequest(http.MethodGet, path, nil)

	res := &SearchResult{}
	if err := c.c.do(req, res); err != nil {
		return nil, err
	}

	return res, nil
}
