package enumr

import (
	"strings"
	"unicode"
)

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

	if len(words) == 0 {
		return ""
	}

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

// toSnakeCaseUpper converts a name to SNAKE_CASE.
func toSnakeCaseUpper(input string) string {
	return strings.ToUpper(toSnakeCase(input))
}

// toTitleCase converts a name to Title Case.
func toTitleCase(name string) string {
	// Split by underscores or spaces, and capitalize the first letter of each word
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == ' '
	})

	if len(words) == 0 {
		return ""
	}

	for i := range words {
		words[i] = capitalize(words[i])
	}
	return strings.Join(words, " ")
}

// capitalize capitalizes the first letter of a string and lowercases the rest.
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
