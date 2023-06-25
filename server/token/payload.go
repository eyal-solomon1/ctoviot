package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	errExpiredToken = errors.New("token has expired")
	errInvalidToken = errors.New("token is invalid")
)

// Payload contains the data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuesAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Valid checks if the token payload is valid or not
func (p Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return errExpiredToken
	}
	return nil
}

// NewPayload creates a new token payload with the provided username
func NewPayload(username string, duratiom time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuesAt:  time.Now(),
		ExpiredAt: time.Now().Add(duratiom),
	}

	return payload, nil

}
