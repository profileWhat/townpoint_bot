package common

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gofrs/uuid"
)

func MarshalUUID(u uuid.UUID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, err := io.WriteString(w, strconv.Quote(u.String()))
		if err != nil {
			log.Println("failed to marshal uuid:", err)
		}
	})
}

func UnmarshalUUID(v interface{}) (uuid.UUID, error) {
	s, ok := v.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid type %T, expect string", v)
	}

	u, err := uuid.FromString(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid uuid, %w", err)
	}

	return u, nil
}
