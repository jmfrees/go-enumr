package test

import (
	"testing"
)

// This test file assumes that go generate has been run and test_string.go exists.
// It verifies that the generated code behaves as expected.

func TestEnumString(t *testing.T) {
	tests := []struct {
		input    Type
		expected string
	}{
		{Foo, "foo"},
		{Bar, "bar"},
		{Baz, "baz"},
		{LongerName, "longer_name"},
	}

	for _, test := range tests {
		if result := test.input.String(); result != test.expected {
			t.Errorf("Type.String() = %q; want %q", result, test.expected)
		}
	}
}

func TestEnumUnmarshalText(t *testing.T) {
	tests := []struct {
		input    string
		expected Type
		wantErr  bool
	}{
		{"foo", Foo, false},
		{"bar", Bar, false},
		{"baz", Baz, false},
		{"longer_name", LongerName, false},
		{"invalid", Type{}, true},
	}

	for _, test := range tests {
		var result Type
		err := result.UnmarshalText([]byte(test.input))
		if (err != nil) != test.wantErr {
			t.Errorf("Type.UnmarshalText(%q) error = %v, wantErr %v", test.input, err, test.wantErr)
			continue
		}
		if !test.wantErr && result != test.expected {
			t.Errorf("Type.UnmarshalText(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
