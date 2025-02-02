package uuidtools

import "github.com/gofrs/uuid"

func New() uuid.UUID {
	u, err := uuid.NewV6()
	if err != nil {
		// Most likely to not fail
		panic(err)
	}
	return u
}

func NewFromString(text string) (uuid.UUID, error) {
	u, err := uuid.FromString(text)
	if err != nil {
		return uuid.Nil, err
	}
	return u, nil
}
