package tokens

import "time"

type Provider interface {
	Generate(payload TokenPayload, exp int) (*Token, error)
	Validate(token string) (*TokenPayload, error)
}

type Token struct {
	Token          string    `json:"token"`
	CreatedAt      time.Time `json:"createdAt"`
	AvailableUntil int       `json:"availableUntil"`
}

type TokenPayload struct {
	UserId uint `json:"userId"`
	RoleId uint `json:"roleId"`
}
