package enumr

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed enum.tmpl
var enumTemplate string

// generateEnumSource generates the actual code for the type, using a template.
func generateEnumSource(
	packageName string,
	enums []EnumInfo,
) ([]byte, error) {
	// Create the template object with a function map for name transformations
	tmplFuncs := template.FuncMap{
		"transformName": func(name, format string) string {
			return transformName(format)(name)
		},
	}

	tmpl, err := template.New("enumTemplate").Funcs(tmplFuncs).Parse(enumTemplate)
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
