// Package enumr provides functionality to generate string representations for enums.
package enumr

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/tools/go/packages"
)

// EnumData is used to pass the necessary data to the template.
type EnumData struct {
	PackageName string
	Enums       []EnumInfo
}

// EnumInfo holds data for a specific enum type.
type EnumInfo struct {
	TypeName   string
	Instances  []InstanceData
	CaseFormat string
}

type TypeSpec struct {
	PackageName string
	TypeSpec    *ast.TypeSpec
}

// InstanceData holds information about each constant instance.
type InstanceData struct {
	Name string
}

// ProcessPackage processes a single Go file to find and generate enums for the given type.
func ProcessPackage(
	pkg *packages.Package,
	typeNames []string,
	nameFormat string,
	output string,
) error {
	var enums []EnumInfo

	for _, typeName := range typeNames {
		// Process type declaration
		if _, err := processTypeSpec(pkg, typeName); err != nil {
			return err
		}

		// Collect instances
		instances := collectInstances(pkg, typeName)
		if len(instances) == 0 {
			// If no instances were found for the given type, return an error
			return fmt.Errorf("failed to find any instances of %s", typeName)
		}

		enums = append(enums, EnumInfo{
			TypeName:   typeName,
			Instances:  instances,
			CaseFormat: nameFormat,
		})
	}

	if len(enums) == 0 {
		return nil
	}

	// Generate the enum code for the type and its instances
	source, err := generateEnumSource(pkg.Name, enums)
	if err != nil {
		return fmt.Errorf("error generating enum source: %w", err)
	}

	// Determine output filename
	var outFileName string
	if output != "" {
		outFileName = output
		// If output is a directory, join with default filename
		if info, err := os.Stat(outFileName); err == nil && info.IsDir() {
			outFileName = filepath.Join(
				outFileName,
				fmt.Sprintf("%s_string.go", strings.ToLower(typeNames[0])),
			)
		}
	} else {
		// Default: <first_type>_string.go in the package directory
		outFileName = filepath.Join(pkg.Dir, fmt.Sprintf("%s_string.go", strings.ToLower(typeNames[0])))
	}

	// Write the generated source to a file
	if err := os.WriteFile(outFileName, source, 0o600); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outFileName, err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Enum code generated successfully", "file", outFileName)

	return nil
}

// processTypeSpec processes the type declarations and ensures the type name matches.
func processTypeSpec(pkg *packages.Package, typeName string) (*TypeSpec, error) {
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			// Only process type declarations
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			// Check if the declaration matches the specified typeName
			for _, spec := range genDecl.Specs {
				typeSpec, okSpec := spec.(*ast.TypeSpec)
				if !okSpec || typeSpec.Name.Name != typeName {
					continue
				}

				return &TypeSpec{
					PackageName: file.Name.Name,
					TypeSpec:    typeSpec,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("type %s not found in package", typeName)
}

// collectInstances processes the var declarations and collects instance names.
func collectInstances(pkg *packages.Package, typeName string) []InstanceData {
	var instances []InstanceData
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			// Only process var declarations
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}

			// Process variable specifications
			for _, spec := range genDecl.Specs {
				valueSpec, okSpec := spec.(*ast.ValueSpec)
				if !okSpec {
					continue
				}

				// Collect instances related to the type
				collectVarsOfType(valueSpec, typeName, &instances)
			}
		}
	}
	return instances
}

// collectVarsOfType checks if the variable's type matches the target struct type and adds to instances.
func collectVarsOfType(valueSpec *ast.ValueSpec, typeName string, instances *[]InstanceData) {
	for i, value := range valueSpec.Values {
		v, ok := value.(*ast.CompositeLit)
		if !ok || v.Type == nil {
			continue
		}

		// Check if the type of the literal matches the target typeName
		if typeIdent, okIdent := v.Type.(*ast.Ident); okIdent && typeIdent.Name == typeName {
			// Add the instance name to the list
			*instances = append(*instances, InstanceData{
				Name: valueSpec.Names[i].Name,
			})
		}
	}
}

// transformName takes a name and a format, then transforms the name accordingly.
func transformName(format string) func(string) string {
	switch format {
	case "snake_case":
		return toSnakeCase
	case "SNAKE_CASE":
		return toSnakeCaseUpper
	case "camelCase":
		return toCamelCase
	case "PascalCase":
		return toPascalCase
	case "Title Case":
		return toTitleCase
	default:
		// If no format is specified, return the original name as-is
		return func(s string) string { return s }
	}
}

// toCamelCase converts a name to camelCase.
func toCamelCase(name string) string {
	// Split by underscores or spaces and lower case the first letter of each word after the first
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == ' '
	})

	if len(words) == 0 {
		return ""
	}

	// Convert first word to lowercase, others to title case
	for i := 1; i < len(words); i++ {
		words[i] = capitalize(words[i])
	}
	return strings.ToLower(words[0]) + strings.Join(words[1:], "")
}

// toPascalCase converts a name to PascalCase.
func toPascalCase(name string) string {
	// Split by underscores or spaces and capitalize the first letter of each word
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == ' '
	})

	// Capitalize every word
	for i := range words {
		words[i] = capitalize(words[i])
	}
	return strings.Join(words, "")
}

// toSnakeCase converts a name to snake_case.
func toSnakeCase(input string) string {
	var result []rune
	for i, r := range input {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func toSnakeCaseUpper(input string) string {
	return strings.ToUpper(toSnakeCase(input))
}

// toTitleCase converts a name to Title Case.
func toTitleCase(name string) string {
	// Split by underscores or spaces, and capitalize the first letter of each word
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == ' '
	})

	for i := range words {
		words[i] = capitalize(words[i])
	}
	return strings.Join(words, " ")
}

func capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	for i := 1; i < len(r); i++ {
		r[i] = unicode.ToLower(r[i])
	}
	return string(r)
}

// generateEnumSource generates the actual code for the type, using a template.
func generateEnumSource(
	packageName string,
	enums []EnumInfo,
) ([]byte, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("generating")

	// Create the template object with a function map for name transformations
	tmplFuncs := template.FuncMap{
		"transformName": func(name, format string) string {
			return transformName(format)(name)
		},
	}

	// Prepare the template for generating the code
	codeTemplate := `// Code generated by enumr. DO NOT EDIT.

package {{.PackageName}}

import (
	"fmt"
	"strings"
)
{{range .Enums}}
{{$format := .CaseFormat}}
// String converts the enum value to its corresponding name field.
func (t {{.TypeName}}) String() string {
	switch t {
	{{range .Instances}}
	case {{.Name}}:
		return "{{transformName .Name $format}}"
	{{end}}
	}
	return ""
}

// MarshalText converts the enum value to a string
func (t {{.TypeName}}) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

// UnmarshalText converts a string to the appropriate enum value
func (t *{{.TypeName}}) UnmarshalText(text []byte) error {
	trimmedText := strings.ReplaceAll(strings.ToLower(string(text)), "\"", "")
	switch trimmedText {
	{{range .Instances}}
	case "{{transformName .Name $format}}":
		*t = {{.Name}}
	{{end}}
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}
{{end}}`

	tmpl, err := template.New("enumTemplate").Funcs(tmplFuncs).Parse(codeTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Define the data to pass to the template
	data := EnumData{
		PackageName: packageName,
		Enums:       enums,
	}

	// Apply the template to the data and write it to the buffer
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}
