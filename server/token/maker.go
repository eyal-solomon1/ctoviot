package token

import "time"

// Maker is an Interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for the requested username
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	// VerifyToken verifies the given token and checks if it's valid
	VerifyToken(token string) (*Payload, error)
}
