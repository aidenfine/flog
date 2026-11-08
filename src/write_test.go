package src

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"
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
		{name: "named field", msg: "received", fields: []Field{{Key: "count", Type: Int32Type, Integer: 2}}, want: `received {"count": 2}` + "\n"},
		{name: "fields without message", fields: MakeFields(1, 2), want: `1 2` + "\n"},
		{name: "escaped string", msg: "text", fields: MakeFields("line\nnext"), want: "text line\nnext" + "\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			if _, err := new(Writer).WriteLog(&out, &Entity{Msg: tt.msg, Fields: tt.fields}, "", "console"); err != nil {
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
	if _, err := new(Writer).WriteLog(&out, &Entity{Fields: fields}, "", "console"); err != nil {
		t.Fatal(err)
	}
	if want := "plain value {name:Alice ids:[1 2]} [1 2]\n"; out.String() != want {
		t.Fatalf("got %q, want %q", out.String(), want)
	}
}

func TestConsoleMarshalerError(t *testing.T) {
	var out bytes.Buffer
	if _, err := new(Writer).WriteLog(&out, &Entity{Fields: MakeFields(failingObject{fmt.Errorf("failed")})}, "", "console"); err == nil {
		t.Fatal("expected marshaler error")
	}
	if out.Len() != 0 {
		t.Fatalf("wrote partial log: %q", out.String())
	}
}

func BenchmarkConsoleValues(b *testing.B) {
	fields := MakeFields("req-123", "GET", "/api/users", 200, 12.5, 42, "1.1.1.1", 1024, false, "api")
	b.ReportAllocs()
	var writer Writer
	for b.Loop() {
		writer.buf.Reset()
		if err := writer.writeConsoleValues(fields); err != nil {
			b.Fatal(err)
		}
	}
}

func TestNamedConsoleFields(t *testing.T) {
	var out bytes.Buffer
	fields := AppendFields(nil, "text", "line\n\"quoted\"", "object", testObject{}, "array", testArray{}, "raw", consoleJSONTrap{})
	if _, err := new(Writer).WriteLog(&out, &Entity{Fields: fields}, "", "console"); err != nil {
		t.Fatal(err)
	}
	want := `{"text": "line\n\"quoted\"", "object": {"name": "Alice", "ids": [1, 2]}, "array": [1, 2], "raw": plain value}` + "\n"
	if out.String() != want {
		t.Fatalf("got %q, want %q", out.String(), want)
	}
}

func TestWriterReuse(t *testing.T) {
	var writer Writer
	var out bytes.Buffer
	now := time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)
	entries := []struct {
		entry        Entity
		caller, want string
	}{
		{Entity{Msg: "long message", LogLevel: "INFO", LoggerName: "app", Time: &now, Fields: AppendFields(nil, "text", "hello", "count", 123)}, " - file.go:12", `2026-11-03T12:00:00Z [INFO] [app] - long message {"text": "hello", "count": 123} - file.go:12` + "\n"},
		{Entity{Msg: "short"}, "", "short\n"},
		{Entity{Fields: MakeFields(1, 2)}, "", "1 2\n"},
		{Entity{Fields: AppendFields(nil, "active", true)}, "", `{"active": true}` + "\n"},
		{Entity{Fields: AppendFields(nil, "nested", testObject{}, "array", testArray{})}, "", `{"nested": {"name": "Alice", "ids": [1, 2]}, "array": [1, 2]}` + "\n"},
		{Entity{Fields: AppendFields(nil, "final", 3)}, "", `{"final": 3}` + "\n"},
	}
	for _, tt := range entries {
		out.Reset()
		got, err := writer.WriteLog(&out, &tt.entry, tt.caller, "console")
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != tt.want || out.String() != tt.want {
			t.Fatalf("got %q, want %q", got, tt.want)
		}
	}
	out.Reset()
	bad := &Entity{Msg: "partial", Fields: AppendFields(nil, "bad", failingObject{fmt.Errorf("failed")})}
	if _, err := writer.WriteLog(&out, bad, "", "console"); err == nil || out.Len() != 0 {
		t.Fatalf("error=%v output=%q", err, out.String())
	}
	if _, err := writer.WriteLog(&out, &Entity{Msg: "recovered", Fields: AppendFields(nil, "ok", true)}, "", "console"); err != nil {
		t.Fatal(err)
	}
	if out.String() != `recovered {"ok": true}`+"\n" {
		t.Fatalf("got %q", out.String())
	}
}

func BenchmarkWriterDiscard(b *testing.B) {
	var writer Writer
	entry := Entity{Msg: "request completed", LogLevel: "INFO", LoggerName: "app", Fields: AppendFields(nil, "request_id", "req-123", "method", "GET", "path", "/api/users", "status", 200, "duration_ms", 12.5, "user_id", 42, "remote_ip", "1.1.1.1", "bytes_sent", 1024, "cached", false, "service", "api")}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := writer.WriteLog(io.Discard, &entry, "", "console"); err != nil {
			b.Fatal(err)
		}
	}
}
