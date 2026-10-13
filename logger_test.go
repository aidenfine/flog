package flog

// func TestLoggerNamedFields(t *testing.T) {
// 	logger := New(false, "app", Config{
// 		LogLevel:   "ERROR",
// 		OutputPath: nil,
// 		ShowTime:   false,
// 	})

// 	var out bytes.Buffer
// 	logger.output = &out
// 	tests := []struct {
// 		name  string
// 		log   func(string, ...any)
// 		level string
// 	}{
// 		{"info", logger.Info, "INFO"},
// 		{"debug", logger.Debug, "DEBUG"},
// 		{"error", logger.Error, "ERROR"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			out.Reset()
// 			tt.log("received", Any("ids", []int{1, 2}), Any("count", int32(2)))
// 			want := "[" + tt.level + "] [app] - received [1,2] 2\n"
// 			if got := out.String(); got != want {
// 				t.Errorf("output = %q, want %q", got, want)
// 			}
// 		})
// 	}
// }
