package aspect

import "testing"

// TestParse exercises Parse on valid and invalid aspect ratio strings.
func TestParse(t *testing.T) {
	// Table covers common aspect ratios plus malformed inputs that should error.
	cases := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"square", "1:1", 1.0, false},
		{"landscape", "3:2", 1.5, false},
		{"portrait", "2:3", 2.0 / 3.0, false},
		{"bad format", "12", 0, true},
		{"zero height", "3:0", 0, true},
		{"non number", "a:2", 0, true},
	}

	for _, tc := range cases {
		// Each case runs as its own subtest for clearer failures in `go test`.
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q, got none", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("Parse(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
