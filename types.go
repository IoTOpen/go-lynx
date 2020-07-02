package lynx

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%s - %d)", e.Message, http.StatusText(e.Code), e.Code)
}

type Installation struct {
	ID           int64    `json:"id"`
	ClientID     int64    `json:"client_id"`
	Name         string   `json:"name"`
	Timezone     string   `json:"timezone"`
	Capabilities []string `json:"capabilities"`
}

type Function struct {
	ID             int64             `json:"id"`
	Type           string            `json:"type"`
	InstallationID int64             `json:"installation_id"`
	Meta           map[string]string `json:"meta"`
	Created        int64             `json:"created"`
	Updated        int64             `json:"updated"`
}

type Device struct {
	ID             int64             `json:"id"`
	Type           string            `json:"type"`
	InstallationID int64             `json:"installation_id"`
	Meta           map[string]string `json:"meta"`
	Created        int64             `json:"created"`
	Updated        int64             `json:"updated"`
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
