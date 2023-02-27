package lynx

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TraceEntry struct {
	ID          string          `json:"id"`
	Path        string          `json:"path"`
	Method      string          `json:"method"`
	Timestamp   float64         `json:"timestamp"`
	UserID      int64           `json:"user_id"`
	Action      TraceAction     `json:"action"`
	ObjectType  TraceObjectType `json:"object_type"`
	ObjectID    int64           `json:"object_id"`
	Description string          `json:"description"`
}

type TracePage struct {
	Total    int64        `json:"total"`
	LastTime float64      `json:"last"`
	Count    int          `json:"count"`
	Data     []TraceEntry `json:"data"`
}

type TraceAction string
type TraceObjectType string

const (
	TraceObjectTypeNone                      = TraceObjectType("")
	TraceObjectTypeInstallation              = TraceObjectType("installation")
	TraceObjectTypeGateway                   = TraceObjectType("gateway")
	TraceObjectTypeOrganization              = TraceObjectType("organization")
	TraceObjectTypeUser                      = TraceObjectType("user")
	TraceObjectTypeDevice                    = TraceObjectType("device")
	TraceObjectTypeFunction                  = TraceObjectType("function")
	TraceObjectTypeSchedule                  = TraceObjectType("schedule")
	TraceObjectTypeNotificationOutput        = TraceObjectType("notification_output")
	TraceObjectTypeNotificationMessage       = TraceObjectType("notification_message")
	TraceObjectTypeOutputExecutor            = TraceObjectType("output_executor")
	TraceObjectTypeEdgeApp                   = TraceObjectType("edge_app")
	TraceObjectTypeEdgeAppInstance           = TraceObjectType("edge_app_instance")
	TraceObjectTypeFile                      = TraceObjectType("file")
	TraceObjectTypeRole                      = TraceObjectType("role")
	TraceObjectTypeGatewayRegistrationPolicy = TraceObjectType("gateway_registration_policy")
	TraceObjectTypeUserRegistrationPolicy    = TraceObjectType("user_registration_policy")
	TraceObjectTypeMQTT                      = TraceObjectType("mqtt")
	TraceObjectTypeTrace                     = TraceObjectType("trace")
)

const (
	TraceActionCreate  = TraceAction("create")
	TraceActionDelete  = TraceAction("delete")
	TraceActionUpdate  = TraceAction("update")
	TraceActionView    = TraceAction("view")
	TraceActionFailed  = TraceAction("failed")
	TraceActionExecute = TraceAction("execute")
	TraceActionAuth    = TraceAction("auth")
)

type TraceOptions struct {
	Limit      int64
	Offset     int64
	From       time.Time
	To         time.Time
	Order      LogOrder
	ObjectType TraceObjectType
	ObjectID   int64
	ID         string
	Action     TraceAction
}

func (c *Client) GetTraces(opts *TraceOptions) (*TracePage, error) {
	if opts == nil {
		return nil, fmt.Errorf("options must be specified")
	}
	res := &TracePage{}
	query := url.Values{
		"from":   []string{fmt.Sprintf("%d", opts.From.Unix())},
		"to":     []string{fmt.Sprintf("%d", opts.To.Unix())},
		"limit":  []string{fmt.Sprintf("%d", opts.Limit)},
		"offset": []string{fmt.Sprintf("%d", opts.Offset)},
		"order":  []string{string(opts.Order)},
	}

	if opts.ObjectType != TraceObjectTypeNone {
		query["object_type"] = []string{string(opts.ObjectType)}
		query["object_id"] = []string{string(opts.ObjectID)}
	} else {
		query["id"] = []string{opts.ID}
	}

	path := fmt.Sprintf("api/v2/trace?%s", query.Encode())
	req := c.newRequest(http.MethodGet, path, nil)
	if err := c.do(req, res); err != nil {
		return nil, err
	}
	return res, nil
}
