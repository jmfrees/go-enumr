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
