package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/jmfrees/go-enumr/pkg/enumr"
)

var (
	typeNames = flag.String("type", "", "comma-separated list type(s) to generate for (required)")
	format    = flag.String("format", "", "format of the name for each enum instance (default: \"\")")
	output    = flag.String("output", "", "output file name (default: dir/<type>_string.go)")
)

func main() {
	flag.Parse()

	// Ensure the -type argument is provided
	if len(*typeNames) == 0 {
		fmt.Println("Error: -type argument is required")
		os.Exit(2)
	}
	typs := strings.Split(*typeNames, ",")

	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	var dir string

	if len(args) == 1 && isDirectory(args[0]) {
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
		fmt.Printf("Error loading package: %v\n", err)
		os.Exit(1)
	}

	outputName := ""
	if output != nil {
		outputName = *output
	}

	// Process the loaded package and files
	err = enumr.ProcessPackage(pkg, typs, nameFormat, outputName)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Enum generation completed successfully!")
}

func loadPackageFromDir(dir string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode:  packages.LoadSyntax, // Load syntax only
		Tests: false,               // Ignore tests for now
	}

	// Load all Go files in the directory
	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to load package from directory %s: %w", dir, err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found in directory %s", dir)
	}

	return pkgs[0], nil
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
