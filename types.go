package lynx

import (
	"fmt"
	"net/http"
	"strconv"
)

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%s - %d)", e.Message, http.StatusText(e.Code), e.Code)
}

type Meta map[string]string

type Installation struct {
	ID           int64    `json:"id"`
	ClientID     int64    `json:"client_id"`
	Name         string   `json:"name"`
	Timezone     string   `json:"timezone"`
	Capabilities []string `json:"capabilities"`
}

type Function struct {
	ID             int64  `json:"id"`
	Type           string `json:"type"`
	InstallationID int64  `json:"installation_id"`
	Meta           Meta   `json:"meta"`
	Created        int64  `json:"created"`
	Updated        int64  `json:"updated"`
}

func (m Meta) AsInt(key string) (int, error) {
	return strconv.Atoi(m[key])
}

func (m Meta) AsFloat64(key string) (float64, error) {
	return strconv.ParseFloat(m[key], 64)
}

type Device struct {
	ID             int64  `json:"id"`
	Type           string `json:"type"`
	InstallationID int64  `json:"installation_id"`
	Meta           Meta   `json:"meta"`
	Created        int64  `json:"created"`
	Updated        int64  `json:"updated"`
}

type Address struct {
	Address string `json:"address"`
	City    string `json:"city"`
	Country string `json:"country"`
	ZIP     string `json:"zip"`
}

type User struct {
	ID        int64   `json:"id"`
	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Role      int64   `json:"role"`
	SmsLogin  bool    `json:"sms_login"`
	Address   Address `json:"address"`
}

type LogEntry struct {
	ClientID       int64   `json:"client_id"`
	InstallationID int64   `json:"installation_id"`
	Message        string  `json:"msg"`
	Timestamp      float64 `json:"timestamp"`
	Topic          string  `json:"topic"`
	Value          float64 `json:"value"`
}

type Status []*LogEntry

func (s Status) Map() map[string]*LogEntry {
	res := make(map[string]*LogEntry, len(s))
	for i, v := range s {
		res[v.Topic] = s[i]
	}
	return res
}

type FunctionList []*Function

func (d FunctionList) MapByID() map[int64]*Function {
	res := make(map[int64]*Function, len(d))
	for i, v := range d {
		res[v.ID] = d[i]
	}
	return res
}

func (d FunctionList) MapBy(key string) map[string]*Function {
	res := make(map[string]*Function, len(d))
	for i, v := range d {
		res[v.Meta[key]] = d[i]
	}
	return res
}

type DeviceList []*Device

func (d DeviceList) MapByID() map[int64]*Device {
	res := make(map[int64]*Device, len(d))
	for i, v := range d {
		res[v.ID] = d[i]
	}
	return res
}

func (d DeviceList) MapBy(key string) map[string]*Device {
	res := make(map[string]*Device, len(d))
	for i, v := range d {
		res[v.Meta[key]] = d[i]
	}
	return res
}
