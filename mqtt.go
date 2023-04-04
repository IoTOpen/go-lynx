package lynx

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Message struct {
	Value     float64 `json:"value"`
	Timestamp float64 `json:"timestamp,omitempty"`
	Msg       string  `json:"msg,omitempty"`
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

// NewMqttOptions returns default mqtt configuration
// conf is a subset of a viper config which can include:
// broker, the MQTT broker URI
// client_id, id to be used by the client
// connection_log, boolean value for enabling/disabling connection logging
func NewMqttOptions(conf *viper.Viper, onConnect mqtt.OnConnectHandler, onLost mqtt.ConnectionLostHandler) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(conf.GetString("broker"))
	opts.SetClientID(conf.GetString("client_id"))
	opts.SetCleanSession(true)
	opts.SetConnectRetryInterval(time.Second * 5)
	opts.SetConnectTimeout(time.Second)
	opts.SetAutoReconnect(true)

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
