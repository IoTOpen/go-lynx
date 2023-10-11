package lynx

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type (
	Error struct {
		Code    int    `json:"-"`
		Message string `json:"message"`
	}
	Filter map[string]string
)

func (e Error) Error() string {
	return fmt.Sprintf("%s (%s - %d)", e.Message, http.StatusText(e.Code), e.Code)
}

type Meta map[string]string

func (m Meta) AsInt(key string) (int, error) {
	return strconv.Atoi(m[key])
}

func (m Meta) AsUint(key string) (uint, error) {
	v, err := strconv.ParseUint(m[key], 10, 64)
	return uint(v), err
}

func (m Meta) AsFloat64(key string) (float64, error) {
	return strconv.ParseFloat(m[key], 64)
}

func (m Meta) AsBool(key string) (bool, error) {
	return strconv.ParseBool(m[key])
}

func (m Meta) AsInt64(key string) (int64, error) {
	return strconv.ParseInt(m[key], 10, 64)
}

func (m Meta) AsUint64(key string) (uint64, error) {
	return strconv.ParseUint(m[key], 10, 64)
}

func (m Meta) AsInt32(key string) (int32, error) {
	v, err := strconv.ParseInt(m[key], 10, 32)
	return int32(v), err
}

func (m Meta) AsUint32(key string) (uint32, error) {
	v, err := strconv.ParseUint(m[key], 10, 32)
	return uint32(v), err
}

func (m Meta) AsInt16(key string) (int16, error) {
	v, err := strconv.ParseInt(m[key], 10, 16)
	return int16(v), err
}

func (m Meta) AsUint16(key string) (uint16, error) {
	v, err := strconv.ParseUint(m[key], 10, 16)
	return uint16(v), err
}

func (m Meta) AsInt8(key string) (int8, error) {
	v, err := strconv.ParseInt(m[key], 10, 8)
	return int8(v), err
}

func (m Meta) AsUint8(key string) (uint8, error) {
	v, err := strconv.ParseUint(m[key], 10, 8)
	return uint8(v), err
}

func (f Filter) ToURLValues() url.Values {
	query := make(url.Values, len(f))
	for k, v := range f {
		query.Set(k, v)
	}
	return query
}

type MetaObject struct {
	Value     string
	Protected bool
}
