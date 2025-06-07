package uuid

import (
	"github.com/google/uuid"
)

// NewGoogle generates a new UUID.
func NewGoogle() string {
	return uuid.New().String()
}

// ParseGoogle parses a UUID from a string.
func ParseGoogle(s string) (string, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// MustParseGoogle parses a UUID from a string and panics if there is an error.
func MustParseGoogle(s string) string {
	u, err := uuid.Parse(s)
	if err != nil {
		panic(err)
	}
	return u.String()
}
