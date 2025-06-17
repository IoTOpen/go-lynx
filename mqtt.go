package lynx

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

type Message struct {
	Value     float64 `json:"value"`
	Timestamp float64 `json:"timestamp,omitempty"`
	Msg       string  `json:"msg,omitempty"`
}

type MQTTMessage struct {
	Topic string
	QoS   byte
	Msg   Message
}

func (m *Message) Time() time.Time {
	whole, fractals := math.Modf(m.Timestamp)
	return time.Unix(int64(whole), int64(fractals*1000000000))
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

// PublishAllTimeout publishes all messages in the provided slice to their respective topics with a specified timeout.
// It returns a slice of errors for any messages that failed to publish within the timeout.
func (c *Client) PublishAllTimeout(messages []MQTTMessage, timeout time.Duration) []error {
	var queue []mqtt.Token
	var topic []string
	var errors []error
	for _, msg := range messages {
		data, _ := json.Marshal(msg.Msg)
		token := c.Mqtt.Publish(msg.Topic, msg.QoS, false, data)
		queue = append(queue, token)
		topic = append(topic, msg.Topic)
	}
	for i, token := range queue {
		ok := token.WaitTimeout(timeout)
		if !ok {
			errors = append(errors, fmt.Errorf("timeout publishing to topic %s", topic[i]))
		} else if err := token.Error(); err != nil {
			errors = append(errors, fmt.Errorf("error publishing to topic %s: %w", topic[i], err))
		}
	}
	return errors
}

// PublishAll publishes all messages in the provided slice to their respective topics.
// It returns a slice of errors for any messages that failed to publish.
func (c *Client) PublishAll(messages []MQTTMessage) []error {
	var queue []mqtt.Token
	var topic []string
	var errors []error
	for _, msg := range messages {
		data, _ := json.Marshal(msg.Msg)
		token := c.Mqtt.Publish(msg.Topic, msg.QoS, false, data)
		queue = append(queue, token)
		topic = append(topic, msg.Topic)
	}
	for i, token := range queue {
		ok := token.Wait()
		if !ok {
			errors = append(errors, fmt.Errorf("timeout publishing to topic %s", topic[i]))
		} else if err := token.Error(); err != nil {
			errors = append(errors, fmt.Errorf("error publishing to topic %s: %w", topic[i], err))
		}
	}
	return errors
}

// Publish publishes a message to the specified topic with the given QoS level.
// It marshals the payload into JSON format before sending.
func (c *Client) Publish(topic string, payload interface{}, qos byte) error {
	data, _ := json.Marshal(payload)
	token := c.Mqtt.Publish(topic, qos, false, data)
	token.WaitTimeout(time.Second)
	return token.Error()
}

// NewMqttOptions returns default mqtt configuration
// conf is a subset of a viper config which can include:
// broker, the MQTT broker URI
// client_id, id to be used by the client
// connection_log, boolean value for enabling/disabling connection logging
// timeout, the connect-timeout to be used on the client. Defaults to 30s
func NewMqttOptions(conf *viper.Viper, onConnect mqtt.OnConnectHandler, onLost mqtt.ConnectionLostHandler) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(conf.GetString("broker"))
	opts.SetClientID(conf.GetString("client_id"))
	opts.SetCleanSession(true)
	opts.SetConnectRetryInterval(time.Second * 5)
	opts.SetAutoReconnect(true)

	if conf.InConfig("timeout") {
		opts.SetConnectTimeout(conf.GetDuration("timeout"))
	} else {
		opts.SetConnectTimeout(time.Second * 30)
	}

	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		if conf.GetBool("connection_log") {
			log.Println("MQTT: connection lost:", err.Error())
		}
		if onLost != nil {
			onLost(c, err)
		}
	})
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		if conf.GetBool("connection_log") {
			log.Println("MQTT: connected")
		}
		if onConnect != nil {
			onConnect(c)
		}
	})
	return opts
}
