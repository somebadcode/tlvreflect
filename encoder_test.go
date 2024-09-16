package tlv

import (
	"bytes"
	"slices"
	"testing"
)

func ptr[T any](v T) *T {
	return &v
}

func TestEncoder(t *testing.T) {
	type args struct {
		options Options
	}
	type testCase[T Size, L Size] struct {
		name    string
		value   any
		args    args
		want    []byte
		wantErr bool
	}

	tests := []testCase[uint16, uint16]{
		{
			name: "hello_world",
			value: &struct {
				Message string `tlv:"0xFF"`
			}{
				Message: "Hello, World!",
			},
			want: []byte{0x0, 0xFF, 0, 13, 'H', 'e', 'l', 'l', 'o', ',', ' ', 'W', 'o', 'r', 'l', 'd', '!'},
		},
		{
			name: "integers",
			value: &struct {
				Int8  int8  `tlv:"0x02"`
				Int16 int16 `tlv:"0x03"`
				Int32 int32 `tlv:"0x04"`
				Int64 int64 `tlv:"0x05"`
			}{
				Int8:  100,
				Int16: -300,
				Int32: 400,
				Int64: 500,
			},
			want: []byte{
				0x00, 0x02, 0, 1, 100,
				0x00, 0x03, 0, 2, 0xFE, 0xD4,
				0x00, 0x04, 0, 4, 0, 0, 0x01, 0x90,
				0x00, 0x05, 0, 8, 0, 0, 0, 0, 0, 0, 0x01, 0xF4,
			},
		},
		{
			name: "unsigned_integers",
			value: &struct {
				UInt8  uint8  `tlv:"0x02"`
				UInt16 uint16 `tlv:"0x03"`
				UInt32 uint32 `tlv:"0x04"`
				UInt64 uint64 `tlv:"0x05"`
			}{
				UInt8:  255,
				UInt16: 300,
				UInt32: 400,
				UInt64: 500,
			},
			want: []byte{
				0x00, 0x02, 0, 1, 0xFF,
				0x00, 0x03, 0, 2, 0x01, 0x2C,
				0x00, 0x04, 0, 4, 0, 0, 0x01, 0x90,
				0x00, 0x05, 0, 8, 0, 0, 0, 0, 0, 0, 0x1, 0xF4,
			},
		},
		{
			name: "basic_slices",
			value: &struct {
				Strings  []string `tlv:"0x10"`
				Integers []int16  `tlv:"0x11"`
			}{
				Strings:  []string{"abc", "defg", "ijklmno"},
				Integers: []int16{-1, 0, 1},
			},
			want: []byte{
				0x00, 0x10, 0, 0x03, 'a', 'b', 'c',
				0x00, 0x10, 0, 0x04, 'd', 'e', 'f', 'g',
				0x00, 0x10, 0, 0x07, 'i', 'j', 'k', 'l', 'm', 'n', 'o',
				0x00, 0x11, 0, 0x2, 0xFF, 0xFF,
				0x00, 0x11, 0, 0x2, 0x00, 0x00,
				0x00, 0x11, 0, 0x2, 0x00, 0x01,
			},
		},
		{
			name: "vague_integer",
			value: &struct {
				Integer int `tlv:"0x00"`
			}{
				Integer: 1,
			},
			wantErr: true,
		},
		{
			name: "not_a_pointer",
			value: struct {
				Integer int `tlv:"0x00"`
			}{
				Integer: 1,
			},
			wantErr: true,
		},
		{
			name:    "not_a_struct",
			value:   ptr("foo"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}

			encoder := NewEncoder[uint16, uint16](w, tt.args.options)

			err := encoder.Encode(tt.value)
			if err != nil && !tt.wantErr {
				t.Errorf("encoder.Encode() error = %v", err)

				return
			} else if err != nil {
				return
			}

			if got := w.Bytes(); slices.Compare(got, tt.want) != 0 {
				t.Errorf("Encode() got = \n\t%v, want \n\t%v", got, tt.want)
			}
		})
	}
}
