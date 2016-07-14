package metadata

import (
	"fmt"
	"time"
)

// ErrExpired is a structured error type for expressing invalidity via
// expiration.
type ErrExpired struct {
	ValidUntil, EvaluatedAt time.Time
}

func (e *ErrExpired) Error() string {
	return fmt.Sprintf("entity was %2f seconds expired; was valid until %s", e.EvaluatedAt.Sub(e.ValidUntil).Seconds(), e.ValidUntil.String())
}
