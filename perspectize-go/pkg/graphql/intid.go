package graphql

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// IntID is a custom scalar that maps GraphQL ID (string) to Go int
type IntID int

// MarshalIntID marshals an int to a GraphQL ID string
func MarshalIntID(id int) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(strconv.Itoa(id)))
	})
}

// UnmarshalIntID unmarshals a GraphQL ID to an int
func UnmarshalIntID(v interface{}) (int, error) {
	switch v := v.(type) {
	case string:
		return strconv.Atoi(v)
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("IntID must be a string or integer, got %T", v)
	}
}
