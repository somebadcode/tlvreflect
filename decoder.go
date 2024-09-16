package tlv

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

type Decoder[T Size, L Size] struct {
	r       io.Reader
	options Options
}

type TagCollisionError[T Size] struct {
	Tag T
}

func (e TagCollisionError[T]) Error() string {
	return fmt.Sprintf("tag collision detected: %v", e.Tag)
}

func NewDecoder[T Size, L Size](r io.Reader, options Options) *Decoder[T, L] {
	dec := &Decoder[T, L]{
		r:       r,
		options: options,
	}

	dec.options.SetDefaults()

	return dec
}

func (d *Decoder[T, L]) Decode(v any) error {
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

	err := d.decode(ref)
	if err != nil {
		return err
	}

	return nil
}

func (d *Decoder[T, L]) decode(v reflect.Value) error {
	typ := v.Type()

	mapper := &mappingDecoder[T, L]{
		TagMap: make(map[T]reflect.Value, typ.NumField()),
		Order:  d.options.ByteOrder,
	}

	for i := 0; i < typ.NumField(); i++ {
		field := v.Field(i)
		fieldType := typ.Field(i)

		tag, err := parseStructTag[T](fieldType.Tag, d.options.TagKey)
		if errors.Is(err, &noTagError{}) {
			continue
		} else if err != nil {
			return err
		}

		if _, alreadyExist := mapper.TagMap[tag]; alreadyExist {
			return &TagCollisionError[T]{
				Tag: tag,
			}
		}

		mapper.TagMap[tag] = field
	}

	return mapper.Decode(d.r)
}
