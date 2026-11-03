package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

func TestEpochTimeEncoders(t *testing.T) {
	tests := []struct {
		name   string
		value  time.Time
		secs   float64
		millis float64
		nanos  int64
	}{
		{"epoch", time.Unix(0, 0), 0, 0, 0},
		{"whole seconds", time.Unix(42, 0), 42, 42000, 42000000000},
		{"fractional seconds", time.Unix(1, 250000000), 1.25, 1250, 1250000000},
		{"fractional milliseconds", time.Unix(0, 125000), 0.000125, 0.125, 125000},
		{"nanosecond precision", time.Unix(0, 1), 0.000000001, 0.000001, 1},
		{"before epoch", time.Unix(-1, 750000000), -0.25, -250, -250000000},
		{"timezone independent", time.Unix(1, 250000000).In(time.FixedZone("UTC+0530", 19800)), 1.25, 1250, 1250000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EpochTimeEncoder(tt.value); got != tt.secs {
				t.Errorf("EpochTimeEncoder() = %v, want %v", got, tt.secs)
			}
			if got := EpochMillisTimeEncoder(tt.value); got != tt.millis {
				t.Errorf("EpochMillisTimeEncoder() = %v, want %v", got, tt.millis)
			}
			if got := EpochNanosTimeEncoder(tt.value); got != tt.nanos {
				t.Errorf("EpochNanosTimeEncoder() = %v, want %v", got, tt.nanos)
			}
		})
	}
}

func TestFormattedTimeEncoders(t *testing.T) {
	tests := []struct {
		name    string
		value   time.Time
		rfc3339 string
		nano    string
		kitchen string
	}{
		{"zero", time.Time{}, "0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z", "12:00AM"},
		{"UTC with nanoseconds", time.Date(2026, time.November, 2, 15, 4, 5, 123456789, time.UTC), "2026-11-02T15:04:05Z", "2026-11-02T15:04:05.123456789Z", "3:04PM"},
		{"trim trailing zeros", time.Date(2026, time.November, 2, 9, 7, 6, 120000000, time.UTC), "2026-11-02T09:07:06Z", "2026-11-02T09:07:06.12Z", "9:07AM"},
		{"positive offset", time.Date(2026, time.November, 2, 12, 0, 0, 0, time.FixedZone("UTC+0530", 19800)), "2026-11-02T12:00:00+05:30", "2026-11-02T12:00:00+05:30", "12:00PM"},
		{"negative offset", time.Date(2026, time.November, 2, 0, 0, 0, 1, time.FixedZone("UTC-0500", -18000)), "2026-11-02T00:00:00-05:00", "2026-11-02T00:00:00.000000001-05:00", "12:00AM"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RFC3339TimeEncoder(tt.value); got != tt.rfc3339 {
				t.Errorf("RFC3339TimeEncoder() = %q, want %q", got, tt.rfc3339)
			}
			if got := RFC3339NanoEncoder(tt.value); got != tt.nano {
				t.Errorf("RFC3339NanoEncoder() = %q, want %q", got, tt.nano)
			}
			if got := KitchenEncoder(tt.value); got != tt.kitchen {
				t.Errorf("KitchenEncoder() = %q, want %q", got, tt.kitchen)
			}
		})
	}
}

func TestEncodeTime(t *testing.T) {
	value := time.Date(2026, time.November, 2, 15, 4, 5, 0, time.UTC)
	for _, tt := range []struct {
		layout string
		want   string
	}{
		{"2006/01/02 15:04:05", "2026/11/02 15:04:05"},
		{"", ""},
	} {
		t.Run(tt.layout, func(t *testing.T) {
			if got := encodeTime(value, tt.layout); got != tt.want {
				t.Errorf("encodeTime() = %q, want %q", got, tt.want)
			}
		})
	}
}

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

func TestJSONEncodingErrors(t *testing.T) {
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
			_, err := encodeFields(MakeFields(tt.value))
			if err == nil {
				t.Fatal("expected encoding error")
			}
			if !strings.Contains(err.Error(), "arg0") {
				t.Errorf("error lacks field name: %v", err)
			}
			if tt.name == "marshaler error" && !errors.Is(err, sentinel) {
				t.Errorf("error not preserved: %v", err)
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
