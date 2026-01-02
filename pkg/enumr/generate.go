// Package enumr provides functionality to generate string representations for enums.
package enumr

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

// Generator handles the enum generation process.
type Generator struct {
	Logger *slog.Logger
}

// NewGenerator creates a new Generator with the given logger.
func NewGenerator(logger *slog.Logger) *Generator {
	return &Generator{Logger: logger}
}

// Generate processes a single Go file to find and generate enums for the given type.
// It returns the generated source code as a byte slice.
func (g *Generator) Generate(
	ctx context.Context,
	pkg *packages.Package,
	typeNames []string,
	nameFormat string,
	marshalField string,
	includeZero bool,
) ([]byte, error) {
	var enums []enumInfo

	for _, typeName := range typeNames {
		// Process type declaration
		typeSpec, err := processTypeSpec(pkg, typeName)
		if err != nil {
			return nil, err
		}

		// Resolve instances (either from directives or by scanning vars)
		resolution, err := g.resolveInstances(ctx, pkg, typeSpec)
		if err != nil {
			return nil, err
		}

		// Validate that if marshalField is specified, all instances have it
		if marshalField != "" {
			for _, instance := range resolution.Instances {
				if _, ok := instance.Fields[marshalField]; !ok {
					return nil, fmt.Errorf(
						"instance %s does not have field %q",
						instance.Name,
						marshalField,
					)
				}
			}
		}

		enums = append(enums, enumInfo{
			TypeName:     typeName,
			Instances:    resolution.Instances,
			CaseFormat:   nameFormat,
			GenerateVars: resolution.GenerateVars,
			IncludeZero:  includeZero,
			MarshalField: marshalField,
			StructFields: typeSpec.Fields,
		})
	}

	if len(enums) == 0 {
		return nil, nil
	}

	g.Logger.LogAttrs(
		ctx,
		slog.LevelDebug,
		"Generating enum source",
		slog.String("package", pkg.Name),
		slog.Any("types", typeNames),
	)

	// Generate the enum code for the type and its instances
	source, err := generateEnumSource(pkg.Name, enums)
	if err != nil {
		return nil, fmt.Errorf("error generating enum source: %w", err)
	}

	return source, nil
}

// GetOutputFilename determines the output filename based on the directory, type name, and output flag.
func GetOutputFilename(dir, firstType, output string) string {
	if output == "" {
		// Default: <first_type>_enum.go in the package directory
		return filepath.Join(dir, fmt.Sprintf("%s_enum.go", toSnakeCase(firstType)))
	}

	// If output is a directory, join with default filename
	if info, err := os.Stat(output); err == nil && info.IsDir() {
		return filepath.Join(output, fmt.Sprintf("%s_enum.go", toSnakeCase(firstType)))
	}

	return output
}

// resolveInstances determines the instances for a type, prioritizing directives over manual scanning.
func (g *Generator) resolveInstances(
	ctx context.Context,
	pkg *packages.Package,
	typeSpec *typeSpec,
) (instanceResolution, error) {
	// 1. Try Directives
	instances := parseDirectives(ctx, g.Logger, typeSpec.Doc, typeSpec.Fields)
	if len(instances) > 0 {
		return instanceResolution{Instances: instances, GenerateVars: true}, nil
	}

	// 2. Fallback to Scanning
	instances = collectInstances(pkg, typeSpec.TypeSpec.Name.Name, typeSpec.Fields)
	if len(instances) > 0 {
		return instanceResolution{Instances: instances, GenerateVars: false}, nil
	}

	return instanceResolution{}, fmt.Errorf(
		"failed to find any instances of %s",
		typeSpec.TypeSpec.Name.Name,
	)
}
