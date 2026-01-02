package enumr

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestGenerateEnumSource(t *testing.T) {
	packageName := "testpkg"
	enums := []enumInfo{
		{
			TypeName:   "MyEnum",
			CaseFormat: "snake_case",
			Instances: []instanceData{
				{Name: "ValueOne"},
				{Name: "ValueTwo"},
			},
		},
	}

	source, err := generateEnumSource(packageName, enums)
	if err != nil {
		t.Fatalf("generateEnumSource failed: %v", err)
	}

	goldenFile := filepath.Join("testdata", "myenum_enum.go.golden")
	if *update {
		if err := os.WriteFile(goldenFile, source, 0o644); err != nil {
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

func TestGenerateEnumSourceWithVars(t *testing.T) {
	packageName := "testpkg"
	fields := []fieldInfo{
		{Name: "Code", Type: "string"},
		{Name: "Desc", Type: "string"},
	}
	enums := []enumInfo{
		{
			TypeName:     "PaymentMethod",
			CaseFormat:   "snake_case",
			GenerateVars: true,
			StructFields: fields,
			Instances: []instanceData{
				{
					Name: "CreditCard",
					Fields: map[string]string{
						"Code": "\"CC\"",
						"Desc": "\"Credit Card\"",
					},
				},
				{
					Name: "PayPal",
					Fields: map[string]string{
						"Code": "\"PP\"",
						"Desc": "\"PayPal\"",
					},
				},
			},
		},
	}

	source, err := generateEnumSource(packageName, enums)
	if err != nil {
		t.Fatalf("generateEnumSource failed: %v", err)
	}

	expectedSnippet := `var (
	CreditCard = PaymentMethod{ Code: "CC", Desc: "Credit Card" }
	PayPal = PaymentMethod{ Code: "PP", Desc: "PayPal" }
)`

	if !strings.Contains(string(source), expectedSnippet) {
		t.Errorf("generated source does not contain expected var block.\nGot:\n%s", source)
	}
}

func TestGenerateEnumSourceWithZero(t *testing.T) {
	packageName := "testpkg"
	enums := []enumInfo{
		{
			TypeName:    "MyEnum",
			CaseFormat:  "snake_case",
			IncludeZero: true,
			Instances: []instanceData{
				{Name: "ValueOne"},
			},
		},
	}

	source, err := generateEnumSource(packageName, enums)
	if err != nil {
		t.Fatalf("generateEnumSource failed: %v", err)
	}

	expectedSnippet := `case "":
		return MyEnum{}, nil`

	if !strings.Contains(string(source), expectedSnippet) {
		t.Errorf("generated source does not contain zero value case.\nGot:\n%s", source)
	}
}
