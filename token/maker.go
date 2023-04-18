package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// Creates a new token for specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// Verifying Token If It is valid or not valid
	VerifyToken(token string) (*Payload, error)
}
