package tlv

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

const (
	DefaultTagKey = "tlv"
)

type Size interface {
	~uint8 | ~uint16
}

type FieldTooLongError struct {
	Length int
	Bits   int
}

func (e *FieldTooLongError) Error() string {
	return fmt.Sprintf("field too long: length=%d bits=%d", e.Length, e.Bits)
}

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e UnsupportedTypeError) Error() string {
	return fmt.Sprintf("unsupported type %s", e.Type)
}

type UnexpectedTypeError struct {
	Type reflect.Type
}

func (e UnexpectedTypeError) Error() string {
	return fmt.Sprintf("unexpected type for variable `%s`, expected pointer to struct", e.Type)
}

type ValueTooLong struct {
	Length int
	Bits   int
}

func (e ValueTooLong) Error() string {
	return fmt.Sprintf("value of length %d is too long for %d bits", e.Length, e.Bits)
}

type noTagError struct{}

func (e *noTagError) Error() string {
	return "no tag"
}

type Options struct {
	TagKey    string
	ByteOrder binary.ByteOrder
}

func (o *Options) SetDefaults() {
	if o.TagKey == "" {
		o.TagKey = DefaultTagKey
	}

	if o.ByteOrder == nil {
		o.ByteOrder = binary.BigEndian
	}
}
