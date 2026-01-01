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
	format := flag.String("format", "", "format of the name for each enum instance (default: \"\")")
	output := flag.String("output", "", "output file name or directory (default: dir/<type>_enum.go)")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.Background()

	// Ensure the -type argument is provided
	if len(*typeNames) == 0 {
		logger.Error("argument is required", "arg", "-type")
		os.Exit(2)
	}
	targetTypes := strings.Split(*typeNames, ",")

	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	var dir string

	if len(args) == 1 && isDirectory(args[0], logger) {
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
	err = generator.Generate(ctx, pkg, targetTypes, nameFormat, outputName)
	if err != nil {
		logger.Error("Error processing file", "error", err)
		os.Exit(1)
	}

	logger.LogAttrs(
		ctx,
		slog.LevelDebug,
		"Enum generation completed successfully",
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
func isDirectory(name string, logger *slog.Logger) bool {
	info, err := os.Stat(name)
	if err != nil {
		logger.Error("Error checking directory", "error", err)
		os.Exit(1)
	}
	return info.IsDir()
}
