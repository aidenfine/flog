package flog

import (
	"bytes"
	"testing"
)

func TestArraysWithSlicesAndStructs(t *testing.T) {
	type item struct {
		Name string
		IDs  []int
	}
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"array of slices", [2][]int{{1, 2}, {3, 4}}, "[[1 2] [3 4]]"},
		{"slice of slices", [][]int{{1, 2}, {3}}, "[[1 2] [3]]"},
		{"array of structs", [2]item{{Name: "alice", IDs: []int{1, 2}}, {Name: "bob", IDs: []int{3}}}, "[{alice [1 2]} {bob [3]}]"},
		{"slice of structs", []item{{Name: "alice", IDs: []int{1, 2}}}, "[{alice [1 2]}]"},
		{"array of struct slices", [2][]item{{{Name: "alice", IDs: []int{1}}}, {{Name: "bob", IDs: []int{2, 3}}}}, "[[{alice [1]}] [{bob [2 3]}]]"},
		{"mixed slice", []any{[]int{1, 2}, item{Name: "alice", IDs: []int{3}}, nil}, "[[1 2] {alice [3]} <nil>]"},
		{"empty array", [0]item{}, "[]"},
		{"nil slice", []item(nil), "[]"},
		{"nil and empty inner slices", [2][]int{nil, {}}, "[[] []]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, mode := range []string{"array helper", "any helper", "loose pairs"} {
				t.Run(mode, func(t *testing.T) {
					var out bytes.Buffer
					logger := &Logger{name: "app", output: &out}
					switch mode {
					case "array helper":
						logger.Info("items", Array("values", tt.value))
					case "any helper":
						logger.Info("items", Any("values", tt.value))
					case "loose pairs":
						logger.Info("items", "values", tt.value)
					}
					want := "[INFO] [app] - items {\"values\": " + tt.want + "}\n"
					if out.String() != want {
						t.Fatalf("got %q, want %q", out.String(), want)
					}
				})
			}
		})
	}
}
