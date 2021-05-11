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

func (m Meta) AsInt(key string) (int, error) {
	return strconv.Atoi(m[key])
}

func (m Meta) AsFloat64(key string) (float64, error) {
	return strconv.ParseFloat(m[key], 64)
}

func (m Meta) AsBool(key string) (bool, error) {
	return strconv.ParseBool(m[key])
}
