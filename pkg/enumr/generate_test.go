package enumr

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"Hello_World", "helloWorld"},
		{"HELLO_WORLD", "helloWorld"},
		{"hello world", "helloWorld"},
		{"simple", "simple"},
	}

	for _, test := range tests {
		if result := toCamelCase(test.input); result != test.expected {
			t.Errorf("toCamelCase(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"Hello_World", "HelloWorld"},
		{"HELLO_WORLD", "HelloWorld"},
		{"hello world", "HelloWorld"},
		{"simple", "Simple"},
	}

	for _, test := range tests {
		if result := toPascalCase(test.input); result != test.expected {
			t.Errorf("toPascalCase(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "Hello World"},
		{"Hello_World", "Hello World"},
		{"HELLO_WORLD", "Hello World"},
		{"hello world", "Hello World"},
		{"simple", "Simple"},
	}

	for _, test := range tests {
		if result := toTitleCase(test.input); result != test.expected {
			t.Errorf("toTitleCase(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestGenerateEnumSource(t *testing.T) {
	packageName := "testpkg"
	enums := []EnumInfo{
		{
			TypeName:   "MyEnum",
			CaseFormat: "snake_case",
			Instances: []InstanceData{
				{Name: "ValueOne"},
				{Name: "ValueTwo"},
			},
		},
	}

	source, err := generateEnumSource(packageName, enums)
	if err != nil {
		t.Fatalf("generateEnumSource failed: %v", err)
	}

	goldenFile := filepath.Join("testdata", "myenum_string.go.golden")
	if *update {
		if err := os.WriteFile(goldenFile, source, 0644); err != nil {
			t.Fatalf("failed to update golden file: %v", err)
		}
	}

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	if !bytes.Equal(source, expected) {
		t.Errorf("generated source does not match golden file.\nExpected:\n%s\nGot:\n%s", expected, source)
	}
}
