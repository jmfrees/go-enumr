package enumr

import (
	"context"
	"fmt"
	"go/ast"
	"log/slog"
	"strconv"
	"strings"
	"unicode"
)

// parseDirectives parses the comment group for enumr directives.
func parseDirectives(ctx context.Context, logger *slog.Logger, doc *ast.CommentGroup, fields []FieldInfo) []InstanceData {
	if doc == nil {
		return nil
	}

	var instances []InstanceData
	for _, comment := range doc.List {
		if instance, ok := parseDirective(ctx, logger, comment.Text, fields); ok {
			instances = append(instances, instance)
		}
	}

	return instances
}

func parseDirective(ctx context.Context, logger *slog.Logger, text string, fields []FieldInfo) (InstanceData, bool) {
	if !strings.HasPrefix(text, "//") {
		return InstanceData{}, false
	}

	// Normalize: "// enumr:Name" -> "enumr:Name"
	content := strings.TrimSpace(strings.TrimPrefix(text, "//"))

	// Optimization: If it doesn't start with "enumr:", it's likely not for us.
	// This avoids parsing unrelated comments like "//go:generate ..." and logging warnings.
	if !strings.HasPrefix(content, "enumr:") {
		return InstanceData{}, false
	}

	// Split the entire line into arguments
	parts := splitArgs(content)
	if len(parts) == 0 {
		return InstanceData{}, false
	}

	// Parse all arguments into a map
	values := parseArgs(ctx, logger, parts)

	// Check for the 'enumr' key which defines the instance name
	name, ok := values["enumr"]
	if !ok {
		return InstanceData{}, false
	}

	// Remove the 'enumr' key so it doesn't get processed as a field
	delete(values, "enumr")

	return InstanceData{
		Name: name,
		Init: buildInitString(fields, values),
	}, true
}

// parseArgs parses the arguments from a directive string into a map.
func parseArgs(ctx context.Context, logger *slog.Logger, args []string) map[string]string {
	values := make(map[string]string, len(args))
	for _, arg := range args {
		key, val, found := strings.Cut(arg, ":")
		if !found {
			logger.LogAttrs(
				ctx,
				slog.LevelWarn,
				"skipping directive argument without value",
				slog.String("arg", arg),
			)
			continue // Skip arguments without a value
		}

		if unquoted, err := strconv.Unquote(val); err == nil {
			values[key] = unquoted
		} else {
			values[key] = val
		}
	}
	return values
}

// buildInitString constructs the initialization string for the struct instance.
func buildInitString(fields []FieldInfo, values map[string]string) string {
	var parts []string
	for _, field := range fields {
		val, ok := values[field.Name]
		if !ok {
			continue
		}

		if field.Type == "string" {
			val = fmt.Sprintf("%q", val)
		}
		parts = append(parts, fmt.Sprintf("%s: %s", field.Name, val))
	}
	return strings.Join(parts, ", ")
}

// splitArgs splits a string into arguments, respecting quotes.
// It handles shell-style quoting (e.g., key:"value with spaces").
func splitArgs(s string) []string {
	var args []string
	var current strings.Builder
	inQuote := false

	for _, r := range s {
		if r == '"' {
			inQuote = !inQuote
			current.WriteRune(r)
		} else if unicode.IsSpace(r) && !inQuote {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
