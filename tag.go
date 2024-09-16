package tlv

import (
	"fmt"
	"reflect"
	"strconv"
)

type InvalidTagError struct {
	Tag   string
	Field string
	Err   error
}

func (e *InvalidTagError) Error() string {
	return fmt.Sprintf("invalid tag %q for field %q", e.Tag, e.Field)
}

func (e *InvalidTagError) Unwrap() error { return e.Err }

func parseStructTag[T Size](tag reflect.StructTag, key string) (T, error) {
	s, ok := tag.Lookup(key)
	if !ok {
		return T(0), &noTagError{}
	}

	if s == "-" {
		return T(0), &noTagError{}
	}

	n, err := strconv.ParseUint(s, 0, reflect.TypeFor[T]().Bits())
	if err != nil {
		return T(0), &InvalidTagError{
			Tag:   s,
			Field: key,
			Err:   err,
		}
	}

	return T(n), nil
}
