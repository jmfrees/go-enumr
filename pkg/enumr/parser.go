package enumr

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"
)

// processTypeSpec processes the type declarations and ensures the type name matches.
func processTypeSpec(pkg *packages.Package, typeName string) (*TypeSpec, error) {
	decl, err := findTypeDeclaration(pkg, typeName)
	if err != nil {
		return nil, err
	}

	fields := extractFields(pkg, decl.spec)

	doc := decl.genDecl.Doc
	if decl.spec.Doc != nil {
		doc = decl.spec.Doc
	}

	return &TypeSpec{
		PackageName: decl.file.Name.Name,
		TypeSpec:    decl.spec,
		Doc:         doc,
		Fields:      fields,
	}, nil
}

// findTypeDeclaration locates the type declaration in the package.
func findTypeDeclaration(pkg *packages.Package, typeName string) (*typeDeclaration, error) {
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			// Only process type declarations
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			// Check if the declaration matches the specified typeName
			for _, spec := range genDecl.Specs {
				typeSpec, okSpec := spec.(*ast.TypeSpec)
				if !okSpec || typeSpec.Name.Name != typeName {
					continue
				}
				return &typeDeclaration{
					file:    file,
					genDecl: genDecl,
					spec:    typeSpec,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("type %s not found in package", typeName)
}

// extractFields extracts field information from a struct type specification.
func extractFields(pkg *packages.Package, typeSpec *ast.TypeSpec) []FieldInfo {
	var fields []FieldInfo
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	for _, field := range structType.Fields.List {
		typeStr := resolveFieldType(pkg, field)
		for _, name := range field.Names {
			fields = append(fields, FieldInfo{Name: name.Name, Type: typeStr})
		}
	}
	return fields
}

// resolveFieldType resolves the type string for a given field.
func resolveFieldType(pkg *packages.Package, field *ast.Field) string {
	if pkg.TypesInfo != nil {
		if tv, ok := pkg.TypesInfo.Types[field.Type]; ok {
			return tv.Type.String()
		}
	}
	if ident, ok := field.Type.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// collectInstances processes the var declarations and collects instance names.
func collectInstances(pkg *packages.Package, typeName string) []InstanceData {
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
				valueSpec, okSpec := spec.(*ast.ValueSpec)
				if !okSpec {
					continue
				}

				// Collect instances related to the type
				collectVarsOfType(valueSpec, typeName, &instances)
			}
		}
	}
	return instances
}

// collectVarsOfType checks if the variable's type matches the target struct type and adds to instances.
func collectVarsOfType(valueSpec *ast.ValueSpec, typeName string, instances *[]InstanceData) {
	for i, value := range valueSpec.Values {
		v, ok := value.(*ast.CompositeLit)
		if !ok || v.Type == nil {
			continue
		}

		// Check if the type of the literal matches the target typeName
		if typeIdent, okIdent := v.Type.(*ast.Ident); okIdent && typeIdent.Name == typeName {
			// Ensure we don't go out of bounds if Names is shorter
			// than Values (though unlikely in valid code).
			if i >= len(valueSpec.Names) {
				continue
			}
			// Add the instance name to the list
			*instances = append(*instances, InstanceData{
				Name: valueSpec.Names[i].Name,
			})
		}
	}
}
