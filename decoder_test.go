package tlvreflect

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	gocmp "github.com/google/go-cmp/cmp"
)

func TestDecoder(t *testing.T) {
	type args struct {
		r       io.Reader
		options Options
	}

	type data struct {
		String  string   `tlv:"0x00"`
		Strings []string `tlv:"0x01"`
	}

	type testCase[T Size, L Size] struct {
		name string
		args args
		want data
	}

	tests := []testCase[uint16, uint16]{
		{
			args: args{
				r: bytes.NewReader([]byte{
					0x00, 0x00, 0x00, 0x0D, 'H', 'e', 'l', 'l', 'o', ',', ' ', 'W', 'o', 'r', 'l', 'd', '!',
					0x00, 0x01, 0x00, 0x03, 'a', 'b', 'c',
					0x00, 0x01, 0x00, 0x04, 'd', 'e', 'f', 'g',
					0x00, 0x01, 0x00, 0x08, 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
				}),
			},
			want: data{
				String: "Hello, World!",
				Strings: []string{
					"abc",
					"defg",
					"hijklmno",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := NewDecoder[uint16, uint16](tt.args.r, tt.args.options)

			var got data

			err := dec.Decode(&got)
			if err != nil {
				t.Errorf("Decoder.Decode() error = %v", err)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				diff := gocmp.Diff(got, tt.want)
				t.Errorf("NewDecoder() = %v\n", diff)
			}
		})
	}
}
