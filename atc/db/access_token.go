package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gopkg.in/square/go-jose.v2/jwt"
)

//go:generate counterfeiter . AccessToken

type AccessToken interface {
	Token() string
	Claims() Claims
}

type accessToken struct {
	token  string
	claims Claims
}

func (a accessToken) Token() string  { return a.token }
func (a accessToken) Claims() Claims { return a.claims }

func scanAccessToken(rcv *accessToken, scan scannable) error {
	return scan.Scan(&rcv.token, &rcv.claims)
}

type Claims struct {
	jwt.Claims
	RawClaims map[string]interface{} `json:"-"`
}

func (c Claims) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.RawClaims)
}

func (c *Claims) UnmarshalJSON(data []byte) error {
	type target Claims
	var t target
	if err := json.Unmarshal(data, &t.Claims); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &t.RawClaims); err != nil {
		return err
	}

	*c = Claims(t)
	return nil
}

func (c Claims) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Claims) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid claims: expected []byte, got %T", value)
	}

	return json.Unmarshal(b, c)
}
