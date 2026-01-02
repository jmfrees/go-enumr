package enumr

import (
	"path/filepath"
	"testing"
)

func TestGetOutputFilename(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		dir       string
		firstType string
		output    string
		expected  string
	}{
		{
			name:      "Default output",
			dir:       "/tmp",
			firstType: "MyEnum",
			output:    "",
			expected:  filepath.Join("/tmp", "my_enum_enum.go"),
		},
		{
			name:      "Explicit filename",
			dir:       "/tmp",
			firstType: "MyEnum",
			output:    "custom.go",
			expected:  "custom.go",
		},
		{
			name:      "Explicit directory",
			dir:       "/tmp",
			firstType: "MyEnum",
			output:    tmpDir,
			expected:  filepath.Join(tmpDir, "my_enum_enum.go"),
		},
		{
			name:      "Snake case conversion",
			dir:       "/tmp",
			firstType: "UserType",
			output:    "",
			expected:  filepath.Join("/tmp", "user_type_enum.go"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetOutputFilename(tt.dir, tt.firstType, tt.output)
			if got != tt.expected {
				t.Errorf("GetOutputFilename(%q, %q, %q) = %q; want %q", tt.dir, tt.firstType, tt.output, got, tt.expected)
			}
		})
	}
}
