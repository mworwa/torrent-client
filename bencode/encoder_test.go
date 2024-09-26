package bencode

import (
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "int",
			args: args{value: 123},
			want: []byte("i123e"),
		},
		{
			name: "string",
			args: args{value: "hello"},
			want: []byte("5:hello"),
		},
		{
			name: "list",
			args: args{value: []interface{}{123, "hello"}},
			want: []byte("li123e5:helloe"),
		},
		{
			name: "dict",
			args: args{value: map[string]interface{}{"key": "value"}},
			want: []byte("d3:key5:valuee"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
