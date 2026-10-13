package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
)

type testObject struct{}

func (testObject) MarshalLogObject(enc ObjectEncoder) error {
	if err := enc.AddString("name", "Alice"); err != nil {
		return err
	}
	return enc.AddArray("ids", testArray{})
}

type testArray struct{}

func (testArray) MarshalLogArray(enc ArrayEncoder) error {
	if err := enc.AppendInt32(1); err != nil {
		return err
	}
	return enc.AppendInt64(2)
}

type failingObject struct{ err error }

func (o failingObject) MarshalLogObject(enc ObjectEncoder) error {
	return o.err
}

func TestEncodeCustomMarshalers(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"object", testObject{}, `{"arg0": {"name": "Alice", "ids": [1, 2]}}`},
		{"array", testArray{}, `{"arg0": [1, 2]}`},
		{"explicit field", Field{Key: "user", Type: ObjectType, Any: testObject{}}, `{"user": {"name": "Alice", "ids": [1, 2]}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeFields(MakeFields(tt.value))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
			if !json.Valid(got) {
				t.Errorf("invalid JSON: %s", got)
			}
		})
	}
}

func TestWriteLogEncodingErrors(t *testing.T) {
	sentinel := errors.New("marshaler failed")
	cycle := map[string]any{}
	cycle["self"] = cycle
	tests := []struct {
		name  string
		value any
	}{
		{"channel", make(chan int)},
		{"function", func() {}},
		{"NaN", math.NaN()},
		{"infinity", math.Inf(1)},
		{"cycle", cycle},
		{"marshaler error", failingObject{sentinel}},
		{"unknown field", Field{Type: FieldType(255)}},
		{"invalid float storage", Field{Type: Float32Type, Any: "wrong"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			_, err := WriteLog(&out, &Entity{Msg: "hello", Fields: MakeFields(tt.value)})
			if err == nil {
				t.Fatal("expected encoding error")
			}
			if !strings.Contains(err.Error(), "arg0") {
				t.Errorf("error lacks field name: %v", err)
			}
			if tt.name == "marshaler error" && !errors.Is(err, sentinel) {
				t.Errorf("error not preserved: %v", err)
			}
			if out.Len() != 0 {
				t.Errorf("wrote partial log: %q", out.String())
			}
		})
	}
}

func TestEncodeNumericFields(t *testing.T) {
	type namedFloat float32
	type namedUint uint32
	values := []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10), uintptr(11), float32(12.5), float64(13.5), namedFloat(14.5), namedUint(15)}
	for _, value := range values {
		t.Run(fmt.Sprintf("%T", value), func(t *testing.T) {
			got, err := encodeFields(MakeFields(value))
			if err != nil {
				t.Fatal(err)
			}
			want := `{"arg0": ` + stringMustJSON(t, value) + `}`
			if string(got) != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})
	}
}

func stringMustJSON(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
