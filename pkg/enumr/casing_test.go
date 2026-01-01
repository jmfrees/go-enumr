package enumr

import (
	"testing"
)

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

func TestTransformName(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		input    string
		expected string
	}{
		{"Default (empty)", "", "MyEnumVal", "MyEnumVal"},
		{"Snake Case", "snake_case", "MyEnumVal", "my_enum_val"},
		{"Camel Case", "camelCase", "MyEnumVal", "myEnumVal"},
		{"Pascal Case", "PascalCase", "my_enum_val", "MyEnumVal"},
		{"Title Case", "Title Case", "my_enum_val", "My Enum Val"},
		{"SNAKE CASE", "SNAKE_CASE", "MyEnumVal", "MY_ENUM_VAL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transform := transformName(tt.format)
			got := transform(tt.input)
			if got != tt.expected {
				t.Errorf("transformName(%q)(%q) = %q; want %q", tt.format, tt.input, got, tt.expected)
			}
		})
	}
}
