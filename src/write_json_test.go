package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestJSONWriter(t *testing.T) {
	var writer Writer
	var out bytes.Buffer
	now := time.Date(2026, 11, 8, 12, 30, 0, 0, time.UTC)
	fields := AppendFields(nil, "msg", "custom", "nested", testObject{}, "array", testArray{}, "duration", 2*time.Second, "nil", nil)
	fields = append(fields, MakeFields(42)...)
	entry := Entity{Msg: "line\n\"quoted\"", LogLevel: "INFO", LoggerName: "app", Time: &now, TimeFormat: "15:04", Fields: fields}
	got, err := writer.WriteLog(&out, &entry, " - file.go:12", "json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, out.Bytes()) || bytes.Count(got, []byte{'\n'}) != 1 {
		t.Fatalf("unexpected record: %q", got)
	}
	var record map[string]any
	if err := json.Unmarshal(got, &record); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"msg": entry.Msg, "level": "INFO", "logger": "app", "time": "12:30", "caller": "file.go:12",
		"fields": map[string]any{"msg": "custom", "nested": map[string]any{"name": "Alice", "ids": []any{float64(1), float64(2)}}, "array": []any{float64(1), float64(2)}, "duration": "2s", "nil": nil, "arg5": float64(42)},
	}
	if !reflect.DeepEqual(record, want) {
		t.Fatalf("got %#v, want %#v", record, want)
	}
	out.Reset()
	sentinel := errors.New("failed")
	if _, err := writer.WriteLog(&out, &Entity{Fields: MakeFields(failingObject{sentinel})}, "", "json"); !errors.Is(err, sentinel) || out.Len() != 0 {
		t.Fatalf("error=%v output=%q", err, out.String())
	}
	for _, encoding := range []string{"json", "console", "json"} {
		out.Reset()
		if _, err := writer.WriteLog(&out, &Entity{Msg: "ok"}, "", encoding); err != nil {
			t.Fatal(err)
		}
		want := "ok\n"
		if encoding == "json" {
			want = "{\"msg\": \"ok\"}\n"
		}
		if out.String() != want {
			t.Fatalf("got %q, want %q", out.String(), want)
		}
	}
}

type objectArray struct{}

func (objectArray) MarshalLogArray(enc ArrayEncoder) error { return enc.AppendObject(testObject{}) }

func TestJSONObjectInsideArray(t *testing.T) {
	got, err := encodeFields(MakeFields(objectArray{}))
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(got) {
		t.Fatalf("invalid JSON: %s", got)
	}
}
