package lynx

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

func (c *Client) MQTTConnect() error {
	x := c.Mqtt.Connect()
	timedOut := x.WaitTimeout(time.Second)
	if x.Error() != nil {
		return x.Error()
	} else if !timedOut {
		return fmt.Errorf("connection timeout")
	}
	return nil
}

func (c *Client) MQTTDisconnect() {
	c.Mqtt.Disconnect(1000)
}

func (c *Client) Publish(topic string, payload interface{}, qos byte) error {
	data, _ := json.Marshal(payload)
	token := c.Mqtt.Publish(topic, qos, false, data)
	token.WaitTimeout(time.Second)
	return token.Error()
}
