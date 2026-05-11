package lynx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return NewClient(&Options{
		APIBase:       server.URL,
		Authenticator: AuthNone{},
		HTTPClient:    server.Client(),
	})
}

func writeJSONResponse(t *testing.T, w http.ResponseWriter, payload any) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}

func TestNormalizeLogOptionsV3(t *testing.T) {
	t.Run("applies defaults to partial options", func(t *testing.T) {
		to := time.Unix(1_700_000_000, 0)

		got, err := normalizeLogOptionsV3(&LogOptionsV3{
			To:     to,
			Limit:  0,
			Offset: -10,
		})
		if err != nil {
			t.Fatalf("normalizeLogOptionsV3() error = %v", err)
		}

		if got.To != to {
			t.Fatalf("To = %v, want %v", got.To, to)
		}
		if got.From != to.Add(-defaultV3LogRange) {
			t.Fatalf("From = %v, want %v", got.From, to.Add(-defaultV3LogRange))
		}
		if got.Limit != defaultV3LogLimit {
			t.Fatalf("Limit = %d, want %d", got.Limit, defaultV3LogLimit)
		}
		if got.Offset != 0 {
			t.Fatalf("Offset = %d, want 0", got.Offset)
		}
		if got.Order != LogOrderDesc {
			t.Fatalf("Order = %q, want %q", got.Order, LogOrderDesc)
		}
		if got.TopicFilter == nil {
			t.Fatal("TopicFilter is nil, want empty slice")
		}
		if len(got.TopicFilter) != 0 {
			t.Fatalf("TopicFilter len = %d, want 0", len(got.TopicFilter))
		}
	})

	t.Run("applies defaults for nil options", func(t *testing.T) {
		got, err := normalizeLogOptionsV3(nil)
		if err != nil {
			t.Fatalf("normalizeLogOptionsV3() error = %v", err)
		}

		if got.To.IsZero() {
			t.Fatal("To is zero, want default now")
		}
		if got.From.IsZero() {
			t.Fatal("From is zero, want default range start")
		}
		if got.To.Sub(got.From) != defaultV3LogRange {
			t.Fatalf("range = %v, want %v", got.To.Sub(got.From), defaultV3LogRange)
		}
		if got.Limit != defaultV3LogLimit {
			t.Fatalf("Limit = %d, want %d", got.Limit, defaultV3LogLimit)
		}
		if got.Order != LogOrderDesc {
			t.Fatalf("Order = %q, want %q", got.Order, LogOrderDesc)
		}
	})

	t.Run("preserves explicit values", func(t *testing.T) {
		from := time.Unix(1_700_000_000, 0)
		to := from.Add(2 * time.Hour)

		got, err := normalizeLogOptionsV3(&LogOptionsV3{
			From:         from,
			To:           to,
			Limit:        123,
			Offset:       45,
			Order:        LogOrderAsc,
			TopicFilter:  []string{"a", "b"},
			AggrMethod:   "avg",
			AggrInterval: 10 * time.Minute,
		})
		if err != nil {
			t.Fatalf("normalizeLogOptionsV3() error = %v", err)
		}

		if got.From != from || got.To != to {
			t.Fatalf("time range = %v -> %v, want %v -> %v", got.From, got.To, from, to)
		}
		if got.Limit != 123 {
			t.Fatalf("Limit = %d, want 123", got.Limit)
		}
		if got.Offset != 45 {
			t.Fatalf("Offset = %d, want 45", got.Offset)
		}
		if got.Order != LogOrderAsc {
			t.Fatalf("Order = %q, want %q", got.Order, LogOrderAsc)
		}
		if !reflect.DeepEqual(got.TopicFilter, []string{"a", "b"}) {
			t.Fatalf("TopicFilter = %v, want %v", got.TopicFilter, []string{"a", "b"})
		}
		if got.AggrMethod != "avg" {
			t.Fatalf("AggrMethod = %q, want avg", got.AggrMethod)
		}
		if got.AggrInterval != 10*time.Minute {
			t.Fatalf("AggrInterval = %v, want %v", got.AggrInterval, 10*time.Minute)
		}
	})

	t.Run("rejects inverted range", func(t *testing.T) {
		_, err := normalizeLogOptionsV3(&LogOptionsV3{
			From: time.Unix(200, 0),
			To:   time.Unix(100, 0),
		})
		if err == nil {
			t.Fatal("normalizeLogOptionsV3() error = nil, want error")
		}
	})

	t.Run("rejects invalid order", func(t *testing.T) {
		_, err := normalizeLogOptionsV3(&LogOptionsV3{
			From:  time.Unix(100, 0),
			To:    time.Unix(200, 0),
			Order: LogOrder("sideways"),
		})
		if err == nil {
			t.Fatal("normalizeLogOptionsV3() error = nil, want error")
		}
	})

	t.Run("rejects negative aggregation interval", func(t *testing.T) {
		_, err := normalizeLogOptionsV3(&LogOptionsV3{
			From:         time.Unix(100, 0),
			To:           time.Unix(200, 0),
			AggrInterval: -time.Minute,
		})
		if err == nil {
			t.Fatal("normalizeLogOptionsV3() error = nil, want error")
		}
	})
}

func TestV3ClientLog_BuildsNormalizedQuery(t *testing.T) {
	to := time.Unix(1_700_000_000, 0)
	wantFrom := to.Add(-defaultV3LogRange).Unix()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/api/v3beta/log/123" {
			t.Fatalf("path = %s, want /api/v3beta/log/123", r.URL.Path)
		}

		q := r.URL.Query()
		if got := q.Get("from"); got != fmt.Sprintf("%d", wantFrom) {
			t.Fatalf("from = %q, want %d", got, wantFrom)
		}
		if got := q.Get("to"); got != fmt.Sprintf("%d", to.Unix()) {
			t.Fatalf("to = %q, want %d", got, to.Unix())
		}
		if got := q.Get("limit"); got != fmt.Sprintf("%d", defaultV3LogLimit) {
			t.Fatalf("limit = %q, want %d", got, defaultV3LogLimit)
		}
		if got := q.Get("offset"); got != "0" {
			t.Fatalf("offset = %q, want 0", got)
		}
		if got := q.Get("order"); got != string(LogOrderDesc) {
			t.Fatalf("order = %q, want %q", got, LogOrderDesc)
		}

		if _, ok := q["aggr_method"]; ok {
			t.Fatalf("aggr_method present = %v, want omitted", q["aggr_method"])
		}
		if _, ok := q["aggr_interval"]; ok {
			t.Fatalf("aggr_interval present = %v, want omitted", q["aggr_interval"])
		}

		writeJSONResponse(t, w, &V3Log{
			Total: 1,
			Count: 1,
			Data:  []LogEntry{},
		})
	})

	_, err := client.V3().Log(123, &LogOptionsV3{
		To:     to,
		Limit:  0,
		Offset: -1,
	})
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}
}

func TestV3ClientLog_IncludesAggregationQueryWhenExplicit(t *testing.T) {
	from := time.Unix(1_700_000_000, 0)
	to := from.Add(1 * time.Hour)
	interval := 10 * time.Minute

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if got := q.Get("aggr_method"); got != "avg" {
			t.Fatalf("aggr_method = %q, want avg", got)
		}
		if got := q.Get("aggr_interval"); got != interval.String() {
			t.Fatalf("aggr_interval = %q, want %q", got, interval.String())
		}

		writeJSONResponse(t, w, &V3Log{
			Total: 1,
			Count: 1,
			Data:  []LogEntry{},
		})
	})

	_, err := client.V3().Log(321, &LogOptionsV3{
		From:         from,
		To:           to,
		Limit:        25,
		Order:        LogOrderAsc,
		AggrMethod:   "avg",
		AggrInterval: interval,
	})
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}
}

func TestV3ClientLog_OmitsAggregationIntervalWithoutMethod(t *testing.T) {
	from := time.Unix(1_700_000_000, 0)
	to := from.Add(1 * time.Hour)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if _, ok := q["aggr_method"]; ok {
			t.Fatalf("aggr_method present = %v, want omitted", q["aggr_method"])
		}
		if _, ok := q["aggr_interval"]; ok {
			t.Fatalf("aggr_interval present = %v, want omitted", q["aggr_interval"])
		}

		writeJSONResponse(t, w, &V3Log{
			Total: 1,
			Count: 1,
			Data:  []LogEntry{},
		})
	})

	_, err := client.V3().Log(321, &LogOptionsV3{
		From:         from,
		To:           to,
		AggrInterval: 15 * time.Minute,
	})
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}
}

func TestV3ClientLog_FallsBackToPostWhenURITooLong(t *testing.T) {
	topics := []string{"a/topic", "b/topic"}
	from := time.Unix(1_700_000_000, 0)
	to := from.Add(1 * time.Hour)

	var calls int

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++

		switch calls {
		case 1:
			if r.Method != http.MethodGet {
				t.Fatalf("first request method = %s, want GET", r.Method)
			}
			q := r.URL.Query()
			if !reflect.DeepEqual(q["topics"], topics) {
				t.Fatalf("first request topics = %v, want %v", q["topics"], topics)
			}
			w.WriteHeader(http.StatusRequestURITooLong)
		case 2:
			if r.Method != http.MethodPost {
				t.Fatalf("second request method = %s, want POST", r.Method)
			}

			q := r.URL.Query()
			if _, ok := q["topics"]; ok {
				t.Fatalf("second request query topics = %v, want omitted", q["topics"])
			}
			if got := q.Get("limit"); got != "25" {
				t.Fatalf("second request limit = %q, want 25", got)
			}
			if got := q.Get("order"); got != string(LogOrderAsc) {
				t.Fatalf("second request order = %q, want %q", got, LogOrderAsc)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			defer r.Body.Close()

			var gotTopics []string
			if err := json.Unmarshal(body, &gotTopics); err != nil {
				t.Fatalf("unmarshal body: %v", err)
			}
			if !reflect.DeepEqual(gotTopics, topics) {
				t.Fatalf("post body topics = %v, want %v", gotTopics, topics)
			}

			writeJSONResponse(t, w, &V3Log{
				Total: 2,
				Count: 2,
				Data: []LogEntry{
					{Topic: "a/topic", Timestamp: float64(from.Unix()), Value: 1},
					{Topic: "b/topic", Timestamp: float64(to.Unix()), Value: 2},
				},
			})
		default:
			t.Fatalf("unexpected extra request #%d", calls)
		}
	})

	got, err := client.V3().Log(777, &LogOptionsV3{
		From:        from,
		To:          to,
		Limit:       25,
		Offset:      0,
		Order:       LogOrderAsc,
		TopicFilter: topics,
	})
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}
	if got == nil {
		t.Fatal("Log() returned nil response")
	}
	if calls != 2 {
		t.Fatalf("request count = %d, want 2", calls)
	}
}
