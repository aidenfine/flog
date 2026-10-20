package src

import (
	"bytes"
	"fmt"
	"testing"
)

func TestWriteLogFields(t *testing.T) {
	tests := []struct {
		name   string
		msg    string
		fields []Field
		want   string
	}{
		{name: "message only", msg: "hello", want: "hello\n"},
		{name: "scalars", msg: "values", fields: MakeFields(int32(-42), uint64(18446744073709551615), float32(1.25), float64(2.5), true, false, "hello world", nil), want: `values -42 18446744073709551615 1.25 2.5 true false hello world <nil>` + "\n"},
		{name: "collections", msg: "values", fields: MakeFields([2]int{1, 2}, []string{"a", "b"}, map[string]int{"count": 2}, struct{ Count int }{2}), want: `values [1 2] [a b] map[count:2] {2}` + "\n"},
		{name: "binary", msg: "bytes", fields: MakeFields([]byte{0, 255}), want: `bytes [0 255]` + "\n"},
		{name: "byte string", msg: "text", fields: []Field{{Type: ByteStringType, Any: []byte("hello")}}, want: `text hello` + "\n"},
		{name: "named field", msg: "received", fields: []Field{{Key: "count", Type: Int32Type, Integer: 2}}, want: `received 2` + "\n"},
		{name: "fields without message", fields: MakeFields(1, 2), want: `1 2` + "\n"},
		{name: "escaped string", msg: "text", fields: MakeFields("line\nnext"), want: "text line\nnext" + "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			if _, err := WriteLog(&out, &Entity{Msg: tt.msg, Fields: tt.fields}); err != nil {
				t.Fatal(err)
			}
			if got := out.String(); got != tt.want {
				t.Errorf("output = %q, want %q", got, tt.want)
			}
		})
	}
}

// JSON hooks must never be called by the console formatter.
type consoleJSONTrap struct{}

func (consoleJSONTrap) MarshalJSON() ([]byte, error) {
	panic("console called MarshalJSON")
}
func (consoleJSONTrap) String() string { return "plain value" }

func TestConsoleIndependentEncoding(t *testing.T) {
	var out bytes.Buffer
	fields := MakeFields(consoleJSONTrap{}, testObject{}, testArray{})
	if _, err := WriteLog(&out, &Entity{Fields: fields}); err != nil {
		t.Fatal(err)
	}
	if want := "plain value {name:Alice ids:[1 2]} [1 2]\n"; out.String() != want {
		t.Fatalf("got %q, want %q", out.String(), want)
	}
}

func TestConsoleMarshalerError(t *testing.T) {
	var out bytes.Buffer
	if _, err := WriteLog(&out, &Entity{Fields: MakeFields(failingObject{fmt.Errorf("failed")})}); err == nil {
		t.Fatal("expected marshaler error")
	}
	if out.Len() != 0 {
		t.Fatalf("wrote partial log: %q", out.String())
	}
}

func BenchmarkConsoleValues(b *testing.B) {
	fields := MakeFields("req-123", "GET", "/api/users", 200, 12.5, 42, "1.1.1.1", 1024, false, "api")
	b.ReportAllocs()
	var buf bytes.Buffer
	for b.Loop() {
		buf.Reset()
		if err := writeConsoleValues(&buf, fields); err != nil {
			b.Fatal(err)
		}
	}
}
