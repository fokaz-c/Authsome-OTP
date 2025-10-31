package otp_type

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Metadata struct {
	TenantUsername   string `json:"tenant_username"`
	IdentityType     string `json:"identity_type"`
	IdentitySource   string `json:"identity_source"`
	IdentityPassword string `json:"identity_password"`
}

func (m Metadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *Metadata) Scan(value any) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("unsupported type for Metadata scan")
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, m)
}
