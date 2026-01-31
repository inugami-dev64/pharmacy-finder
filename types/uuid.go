package types

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type UUID uuid.UUID

func (uuid *UUID) String() string {
	bytes := [16]byte(*uuid)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", bytes[:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
}

func (uuid UUID) MarshalText() ([]byte, error) {
	return []byte(uuid.String()), nil
}

func (u *UUID) UnmarshalText(text []byte) error {
	val, err := uuid.Parse(string(text))
	if err != nil {
		return err
	}
	*u = UUID(val)
	return nil
}

func (u *UUID) Scan(value interface{}) error {
	val, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprintf("failed to unmarshal flight ID:", value))
	}

	*u = UUID(uuid.MustParse(val))
	return nil
}
