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
func (g *Generator) Generate(
	ctx context.Context,
	pkg *packages.Package,
	typeNames []string,
	nameFormat string,
	output string,
) error {
	var enums []EnumInfo

	for _, typeName := range typeNames {
		// Process type declaration
		typeSpec, err := processTypeSpec(pkg, typeName)
		if err != nil {
			return err
		}

		// Resolve instances (either from directives or by scanning vars)
		resolution, err := g.resolveInstances(ctx, pkg, typeSpec)
		if err != nil {
			return err
		}

		enums = append(enums, EnumInfo{
			TypeName:     typeName,
			Instances:    resolution.Instances,
			CaseFormat:   nameFormat,
			GenerateVars: resolution.GenerateVars,
		})
	}

	if len(enums) == 0 {
		return nil
	}

	g.Logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"Generating enum source",
		slog.String("package", pkg.Name),
		slog.Any("types", typeNames),
	)

	// Generate the enum code for the type and its instances
	source, err := generateEnumSource(pkg.Name, enums)
	if err != nil {
		return fmt.Errorf("error generating enum source: %w", err)
	}

	// Determine output filename
	outFileName := getOutputFilename(pkg.Dir, typeNames[0], output)

	// Write the generated source to a file
	if err = os.WriteFile(outFileName, source, 0o644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outFileName, err)
	}

	g.Logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"Enum code generated successfully",
		slog.String("file", outFileName),
	)

	return nil
}

// getOutputFilename determines the output filename based on the directory, type name, and output flag.
func getOutputFilename(dir, firstType, output string) string {
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
func (g *Generator) resolveInstances(ctx context.Context, pkg *packages.Package, typeSpec *TypeSpec) (InstanceResolution, error) {
	// 1. Try Directives
	instances := parseDirectives(ctx, g.Logger, typeSpec.Doc, typeSpec.Fields)
	if len(instances) > 0 {
		return InstanceResolution{Instances: instances, GenerateVars: true}, nil
	}

	// 2. Fallback to Scanning
	instances = collectInstances(pkg, typeSpec.TypeSpec.Name.Name)
	if len(instances) > 0 {
		return InstanceResolution{Instances: instances, GenerateVars: false}, nil
	}

	return InstanceResolution{}, fmt.Errorf("failed to find any instances of %s", typeSpec.TypeSpec.Name.Name)
}
