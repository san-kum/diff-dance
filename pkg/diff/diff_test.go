package diff

import (
	"strings"
	"testing"
)

func TestFiles(t *testing.T) {
	tests := []struct {
		name    string
		input1  string
		input2  string
		want    []Diff
		wantErr bool
	}{
		{
			name:   "identical files",
			input1: "line1\nline2",
			input2: "line1\nline2",
			want: []Diff{
				{Line: "line1", Type: "same"},
				{Line: "line2", Type: "same"},
			},
			wantErr: false,
		},
		{
			name:   "one line added",
			input1: "line1",
			input2: "line1\nline2",
			want: []Diff{
				{Line: "line1", Type: "same"},
				{Line: "line2", Type: "add"},
			},
			wantErr: false,
		},
		{
			name:   "one line removed",
			input1: "line1\nline2",
			input2: "line1",
			want: []Diff{
				{Line: "line1", Type: "same"},
				{Line: "line2", Type: "remove"},
			},
			wantErr: false,
		},
		{
			name:   "one line changed",
			input1: "line1\nline2",
			input2: "line1\nline3",
			want: []Diff{
				{Line: "line1", Type: "same"},
				{Line: "line2", Type: "remove"},
				{Line: "line3", Type: "add"},
			},
			wantErr: false,
		},
		{
			name:    "empty files",
			input1:  "",
			input2:  "",
			want:    []Diff{},
			wantErr: false,
		},
		{
			name:   "added to the beginning", //Added new test case!
			input1: "line2",
			input2: "line1\nline2",
			want: []Diff{
				{Line: "line1", Type: "add"},
				{Line: "line2", Type: "same"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader1 := strings.NewReader(tt.input1)
			reader2 := strings.NewReader(tt.input2)

			got, err := Files(reader1, reader2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Files() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !diffsEqual(got, tt.want) { // Helper function for comparing diffs
				t.Errorf("Files() = %v, want %v", got, tt.want)
			}
		})
	}
}

// diffsEqual compares two slices of Diff structs.
func diffsEqual(a, b []Diff) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
