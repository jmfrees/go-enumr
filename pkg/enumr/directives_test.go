package enumr

import (
	"go/ast"
	"log/slog"
	"os"
	"reflect"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"key:value", []string{"key:value"}},
		{"key:value flag:true", []string{"key:value", "flag:true"}},
		{"key:\"quoted value\"", []string{"key:\"quoted value\""}},
		{"key:\"quoted value\" flag:true", []string{"key:\"quoted value\"", "flag:true"}},
		{"key:\"value with spaces\" another:val", []string{"key:\"value with spaces\"", "another:val"}},
	}

	for _, test := range tests {
		result := splitArgs(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("splitArgs(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestParseDirectives(t *testing.T) {
	fields := []FieldInfo{
		{Name: "Code", Type: "string"},
		{Name: "Desc", Type: "string"},
		{Name: "IsActive", Type: "bool"},
		{Name: "Int", Type: "int"},
		{Name: "Float", Type: "float64"},
		{Name: "Slice", Type: "[]string"},
		{Name: "Const", Type: "MyConstType"},
	}

	tests := []struct {
		name      string
		directive string
		wantCount int
		wantName  string
		wantInit  string
	}{
		{
			name:      "Full fields",
			directive: "//enumr:Item1 Code:I1 Desc:\"Item One\" IsActive:true",
			wantCount: 1,
			wantName:  "Item1",
			wantInit:  "Code: \"I1\", Desc: \"Item One\", IsActive: true",
		},
		{
			name:      "Partial fields",
			directive: "// enumr:Item2 Code:I2 Desc:ItemTwo",
			wantCount: 1,
			wantName:  "Item2",
			wantInit:  "Code: \"I2\", Desc: \"ItemTwo\"",
		},
		{
			name:      "No fields",
			directive: "// enumr:Item3",
			wantCount: 1,
			wantName:  "Item3",
			wantInit:  "",
		},
		{
			name:      "Ignored comment",
			directive: "// regular comment",
			wantCount: 0,
		},
		{
			name:      "Int and Float",
			directive: "//enumr:Item4 Int:42 Float:3.14",
			wantCount: 1,
			wantName:  "Item4",
			wantInit:  "Int: 42, Float: 3.14",
		},
		{
			name:      "Slice no spaces",
			directive: "//enumr:Item5 Slice:[]string{\"a\",\"b\"}",
			wantCount: 1,
			wantName:  "Item5",
			wantInit:  "Slice: []string{\"a\",\"b\"}",
		},
		{
			name:      "Slice with spaces (quoted)",
			directive: "//enumr:Item6 Slice:\"[]string{\\\"a\\\", \\\"b\\\"}\"",
			wantCount: 1,
			wantName:  "Item6",
			wantInit:  "Slice: []string{\"a\", \"b\"}",
		},
		{
			name:      "Constant reference",
			directive: "//enumr:Item7 Const:SomeConstant",
			wantCount: 1,
			wantName:  "Item7",
			wantInit:  "Const: SomeConstant",
		},
		{
			name:      "Quoted string with spaces",
			directive: "//enumr:Item8 Desc:\"two words\"",
			wantCount: 1,
			wantName:  "Item8",
			wantInit:  "Desc: \"two words\"",
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &ast.CommentGroup{
				List: []*ast.Comment{{Text: tt.directive}},
			}

			instances := parseDirectives(t.Context(), logger, doc, fields)

			if len(instances) != tt.wantCount {
				t.Fatalf("got %d instances, want %d", len(instances), tt.wantCount)
			}

			if tt.wantCount > 0 {
				inst := instances[0]
				if inst.Name != tt.wantName {
					t.Errorf("Name = %q; want %q", inst.Name, tt.wantName)
				}
				gotInit := renderInit(inst, fields)
				if gotInit != tt.wantInit {
					t.Errorf("Init = %q; want %q", gotInit, tt.wantInit)
				}
			}
		})
	}
}
