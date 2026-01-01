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
	words := splitIntoWords(name)
	if len(words) == 0 {
		return ""
	}

	// Lowercase the first word
	words[0] = strings.ToLower(words[0])

	// Capitalize subsequent words
	for i := 1; i < len(words); i++ {
		words[i] = capitalize(words[i])
	}

	return strings.Join(words, "")
}

// toPascalCase converts a name to PascalCase.
func toPascalCase(name string) string {
	words := splitIntoWords(name)
	if len(words) == 0 {
		return ""
	}

	for i := range words {
		words[i] = capitalize(words[i])
	}
	return strings.Join(words, "")
}

// toSnakeCase converts a name to snake_case.
func toSnakeCase(input string) string {
	words := splitIntoWords(input)
	if len(words) == 0 {
		return ""
	}

	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "_")
}

// toSnakeCaseUpper converts a name to SNAKE_CASE.
func toSnakeCaseUpper(input string) string {
	return strings.ToUpper(toSnakeCase(input))
}

// toTitleCase converts a name to Title Case.
func toTitleCase(name string) string {
	words := splitIntoWords(name)
	if len(words) == 0 {
		return ""
	}

	for i := range words {
		words[i] = capitalize(words[i])
	}
	return strings.Join(words, " ")
}

// splitIntoWords splits a string into words based on casing and delimiters.
func splitIntoWords(s string) []string {
	var words []string
	var currentWord []rune

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '_' || r == ' ' {
			if len(currentWord) > 0 {
				words = append(words, string(currentWord))
				currentWord = nil
			}
			continue
		}

		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			// Case 1: aB (lower -> upper)
			if unicode.IsLower(prev) {
				if len(currentWord) > 0 {
					words = append(words, string(currentWord))
					currentWord = nil
				}
			} else if unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				// Case 2: ABc (upper -> upper -> lower) - split before the last upper
				// e.g. JSONParser -> JSON, Parser. We are at 'P'. Prev is 'N'. Next is 'a'.
				if len(currentWord) > 0 {
					words = append(words, string(currentWord))
					currentWord = nil
				}
			}
		}
		currentWord = append(currentWord, r)
	}
	if len(currentWord) > 0 {
		words = append(words, string(currentWord))
	}
	return words
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
