package uuid

import (
	"github.com/google/uuid"
)

// NewGoogle returns a newly generated UUID as a string.
func NewGoogle() string {
	return uuid.New().String()
}

// ParseGoogle attempts to parse a UUID from the provided string and returns its canonical string representation.
// Returns an error if the input is not a valid UUID.
func ParseGoogle(s string) (string, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// MustParseGoogle parses a UUID from the provided string and returns it as a string, panicking if parsing fails.
func MustParseGoogle(s string) string {
	u, err := uuid.Parse(s)
	if err != nil {
		panic(err)
	}
	return u.String()
}
