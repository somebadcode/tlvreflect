package tlvreflect

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type mappingDecoder[T Size, L Size] struct {
	TagMap map[T]reflect.Value
	Order  binary.ByteOrder
}

func (d mappingDecoder[T, L]) Decode(r io.Reader) error {
	buf := bufio.NewReader(r)

	for {
		_, err := buf.Peek(1)
		if errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}

		err = d.decode(buf)
		if err != nil {
			return err
		}
	}
}

func (d mappingDecoder[T, L]) decode(r io.Reader) error {
	var tag T

	err := binary.Read(r, d.Order, &tag)
	if err != nil {
		return fmt.Errorf("failed to read tag: %w", err)
	}

	var length L

	err = binary.Read(r, d.Order, &length)
	if err != nil {
		return fmt.Errorf("failed to read length: %w", err)
	}

	v, ok := d.TagMap[tag]
	if !ok {
		return fmt.Errorf("tag %v not found", tag)
	}

	switch v.Kind() {
	case reflect.Slice:
		if v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 10))
		}

		b := make([]byte, length)

		err = binary.Read(r, d.Order, b)
		if err != nil {
			return fmt.Errorf("failed to read data: %w", err)
		}

		err = d.decodeIntoSlice(v, b)
		if err != nil {
			return fmt.Errorf("failed to decode data into slice: %w", err)
		}

	case reflect.String:
		b := make([]byte, length)

		err = binary.Read(r, d.Order, b)
		if err != nil {
			return fmt.Errorf("failed to read data: %w", err)
		}

		v.SetString(string(b))

	default:
		err = binary.Read(r, d.Order, v.Interface())
		if err != nil {
			return fmt.Errorf("failed to read value with tag %d: %w", tag, err)
		}
	}

	return nil
}

func (d mappingDecoder[T, L]) decodeIntoSlice(v reflect.Value, data []byte) error {
	switch v.Type().Elem().Kind() {
	case reflect.String:
		v.Set(reflect.Append(v, reflect.ValueOf(string(data))))

	case reflect.Uint8:
		v.Set(reflect.Append(v, reflect.ValueOf(data)))

	default:
		_, err := binary.Decode(data, d.Order, v.Interface())
		if err != nil {
			return err
		}
	}

	return nil
}
