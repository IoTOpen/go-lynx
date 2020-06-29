package lynx

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
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
