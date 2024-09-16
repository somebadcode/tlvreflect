package tlv

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type Encoder[T Size, L Size] struct {
	w       io.Writer
	options Options
}

func NewEncoder[T Size, L Size](w io.Writer, options Options) *Encoder[T, L] {
	enc := &Encoder[T, L]{
		w:       w,
		options: options,
	}

	enc.options.SetDefaults()

	return enc
}

func (e *Encoder[T, L]) Encode(v any) error {
	ptr := reflect.ValueOf(v)

	if ptr.Kind() != reflect.Pointer {
		return fmt.Errorf("parse: %w", &UnexpectedTypeError{
			Type: ptr.Type(),
		})
	}

	ref := ptr.Elem()

	if ref.Kind() != reflect.Struct {
		return fmt.Errorf("parse: %w", &UnexpectedTypeError{
			Type: ref.Type(),
		})
	}

	err := e.encode(ref)
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder[T, L]) encode(v reflect.Value) error {
	typ := v.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := v.Field(i)
		fieldType := typ.Field(i)

		tag, err := parseStructTag[T](fieldType.Tag, e.options.TagKey)
		if errors.Is(err, &noTagError{}) {
			// skip fields that do not have a struct tag.
			continue
		} else if err != nil {
			return err
		}

		switch field.Kind() {
		case reflect.Slice:
			err = e.encodeSlice(field, tag)
		default:
			err = e.encodeData(field, tag)
		}

		if err != nil {
			return fmt.Errorf("failed to encode field %q: %w", fieldType.Name, err)
		}
	}

	return nil
}

func (e *Encoder[T, L]) encodeSlice(v reflect.Value, tag T) error {
	for i := 0; i < v.Len(); i++ {
		err := e.encodeData(v.Index(i), tag)
		if err != nil {

		}
	}

	return nil
}

func (e *Encoder[T, L]) encodeData(v reflect.Value, tag T) error {
	var data any
	switch v.Kind() {
	case reflect.String:
		data = []byte(v.Interface().(string))

	case reflect.Invalid:
		return &UnsupportedTypeError{
			Type: v.Type(),
		}

	default:
		data = v.Interface()
	}

	b, err := binary.Append(nil, e.options.ByteOrder, data)
	if err != nil {
		return err
	}

	l := L(len(b))
	if int(l) != len(b) {
		return &FieldTooLongError{
			Length: len(b),
			Bits:   reflect.TypeFor[L]().Bits(),
		}
	}

	err = binary.Write(e.w, e.options.ByteOrder, tag)
	if err != nil {
		return err
	}

	err = binary.Write(e.w, e.options.ByteOrder, l)
	if err != nil {
		return err
	}

	_, err = e.w.Write(b)
	if err != nil {
		return err
	}

	return nil
}
