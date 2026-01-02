package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/jmfrees/go-enumr/pkg/enumr"
)

func main() {
	typeNames := flag.String("type", "", "comma-separated list type(s) to generate for (required)")
	format := flag.String(
		"format",
		"",
		"format of the name for each enum instance (default: preserve case)",
	)
	output := flag.String(
		"output",
		"",
		"output file name or directory (default: dir/<type>_enum.go)",
	)
	marshalField := flag.String(
		"marshal-field",
		"",
		"field to use for marshaling (String/MarshalText)",
	)
	zero := flag.Bool("zero", false, "allow zero value (empty string) during parsing")
	dryRun := flag.Bool("dry-run", false, "perform a trial run with no changes made")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.Background()

	// Ensure the -type argument is provided
	if len(*typeNames) == 0 {
		logger.ErrorContext(ctx, "argument is required", "arg", "-type")
		os.Exit(2)
	}
	targetTypes := strings.Split(*typeNames, ",")

	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	var dir string

	if len(args) == 1 && isDirectory(ctx, args[0], logger) {
		dir = args[0]
	} else {
		dir = filepath.Dir(args[0])
	}

	nameFormat := ""
	if format != nil && len(*format) > 0 {
		nameFormat = *format
	}

	// Load the package
	pkg, err := loadPackageFromDir(dir)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"Error loading package",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	outputName := ""
	if output != nil {
		outputName = *output
	}

	// Process the loaded package and files
	generator := enumr.NewGenerator(logger)
	source, err := generator.Generate(ctx, pkg, targetTypes, nameFormat, *marshalField, *zero)
	if err != nil {
		logger.ErrorContext(ctx, "Error processing file", "error", err)
		os.Exit(1)
	}

	if source == nil {
		logger.LogAttrs(ctx, slog.LevelInfo, "No enums found to generate")
		return
	}

	// Determine output filename
	outFileName := enumr.GetOutputFilename(pkg.Dir, targetTypes[0], outputName)

	if *dryRun {
		logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"Dry run enabled, no files written",
			slog.String("file", outFileName),
		)
		fmt.Println(string(source))
		return
	}

	// Write the generated source to a file
	if err = os.WriteFile(outFileName, source, 0o644); err != nil {
		logger.ErrorContext(ctx, "Error writing file", "file", outFileName, "error", err)
		os.Exit(1)
	}

	logger.LogAttrs(
		ctx,
		slog.LevelDebug,
		"Enum generation completed successfully",
		slog.String("file", outFileName),
	)
}

func loadPackageFromDir(dir string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode:  packages.LoadSyntax, // Load syntax only
		Tests: false,               // Ignore tests for now
	}

	// Load all Go files in the directory
	packagesList, err := packages.Load(cfg, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to load package from directory %s: %w", dir, err)
	}

	if len(packagesList) == 0 {
		return nil, fmt.Errorf("no packages found in directory %s", dir)
	}

	return packagesList[0], nil
}

// isDirectory reports whether the named file is a directory.
func isDirectory(ctx context.Context, name string, logger *slog.Logger) bool {
	info, err := os.Stat(name)
	if err != nil {
		logger.ErrorContext(ctx, "Error checking directory", "error", err)
		os.Exit(1)
	}
	return info.IsDir()
}
