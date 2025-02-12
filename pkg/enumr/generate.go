package enumr

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/tools/go/packages"
)

// EnumData is used to pass the necessary data to the template
type EnumData struct {
	PackageName string
	TypeName    string
	Instances   []InstanceData
	CaseFormat  string
}

type TypeSpec struct {
	PackageName string
	TypeSpec    *ast.TypeSpec
}

// InstanceData holds information about each constant instance
type InstanceData struct {
	Name string
}

// ProcessPackage processes a single Go file to find and generate enums for the given type
func ProcessPackage(pkg *packages.Package, typeName string, nameFormat string) error {
	var instances []InstanceData

	// Process type declaration
	var err error
	var typeSpec *TypeSpec
	if typeSpec, err = processTypeSpec(pkg, typeName); err != nil {
		return err
	}

	// Collect instances
	if instances, err = collectInstances(pkg, typeName); err != nil {
		return err
	}

	// If no instances were found for the given type, return an error
	if len(instances) == 0 {
		return fmt.Errorf("failed to find any instances of %s", typeName)
	}

	// Generate the enum code for the type and its instances
	if err := generateEnumCode(pkg, typeSpec.PackageName, typeName, nameFormat, instances); err != nil {
		return fmt.Errorf("error generating enum for type %s in file %s: %w", typeName, pkg.Dir, err)
	}

	return nil
}

// processTypeSpec processes the type declarations and ensures the type name matches
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
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != typeName {
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

// collectInstances processes the var declarations and collects instance names
func collectInstances(pkg *packages.Package, typeName string) ([]InstanceData, error) {
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
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				// Collect instances related to the type
				if err := collectVarsOfType(pkg, valueSpec, typeName, &instances); err != nil {
					return instances, err
				}
			}
		}
	}
	return instances, nil
}

// collectVarsOfType checks if the variable's type matches the target struct type and adds to instances
func collectVarsOfType(pkg *packages.Package, valueSpec *ast.ValueSpec, typeName string, instances *[]InstanceData) error {
	for i, value := range valueSpec.Values {
		v, ok := value.(*ast.CompositeLit)
		if !ok || v.Type == nil {
			continue
		}

		// Check if the type of the literal matches the target typeName
		if typeIdent, ok := v.Type.(*ast.Ident); ok && typeIdent.Name == typeName {
			// Add the instance name to the list
			*instances = append(*instances, InstanceData{
				Name: valueSpec.Names[i].Name,
			})
		}
	}
	return nil
}

// transformName takes a name and a format, then transforms the name accordingly
func transformName(format string) func(string) string {
	switch format {
	case "snake_case":
		return toSnakeCase
	case "SNAKE_CASE":
		return toSnakeCaseUpper
	default:
		// If no format is specified, return the original name as-is
		return func(s string) string { return s }
	}
}

// toCamelCase converts a name to camelCase
func toCamelCase(name string) string {
	// Split by underscores or spaces and lower case the first letter of each word after the first
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == ' '
	})

	// Convert first word to lowercase, others to title case
	for i := 1; i < len(words); i++ {
		words[i] = strings.ToTitle(words[i])
	}
	return strings.ToLower(words[0]) + strings.Join(words[1:], "")
}

// toPascalCase converts a name to PascalCase
func toPascalCase(name string) string {
	// Split by underscores or spaces and capitalize the first letter of each word
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == ' '
	})

	// Capitalize every word
	for i := 0; i < len(words); i++ {
		words[i] = strings.ToTitle(words[i])
	}
	return strings.Join(words, "")
}

// toSnakeCase converts a name to snake_case
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

// toTitleCase converts a name to Title Case
func toTitleCase(name string) string {
	// Split by underscores or spaces, and capitalize the first letter of each word
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == ' '
	})

	for i := 0; i < len(words); i++ {
		words[i] = strings.ToTitle(words[i])
	}
	return strings.Join(words, " ")
}

// generateEnumCode generates the actual code for the type, using a template
func generateEnumCode(pkg *packages.Package, packageName string, typeName string, nameFormat string, instances []InstanceData) error {
	fmt.Println("generating")

	// Create an output file for the generated code
	outFileName := fmt.Sprintf("%s/%s_string.go", pkg.Dir, packageName)
	outFile, err := os.Create(outFileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outFileName, err)
	}
	defer outFile.Close()

	// Prepare the template for generating the code
	codeTemplate := `// Code generated by enumr. DO NOT EDIT.

package {{.PackageName}}

import (
	"fmt"
	"strings"
)

// String converts the enum value to its corresponding name field.
func (t {{.TypeName}}) String() string {
	switch t {
	{{range .Instances}}
	case {{.Name}}:
		return "{{.Name | transformName}}"
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
	case "{{.Name | transformName}}":
		*t = {{.Name}}
	{{end}}
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}
`

	// Create the template object with a function map for name transformations
	tmplFuncs := template.FuncMap{
		"transformName": transformName(nameFormat),
	}

	tmpl, err := template.New("enumTemplate").Funcs(tmplFuncs).Parse(codeTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Define the data to pass to the template
	data := EnumData{
		PackageName: packageName,
		TypeName:    typeName,
		Instances:   instances,
	}

	// Apply the template to the data and write it to the file
	err = tmpl.Execute(outFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("Enum code generated successfully in %s\n", outFileName)
	return nil
}
