package enumr

import (
	"go/ast"
)

// enumData is used to pass the necessary data to the template.
type enumData struct {
	PackageName string
	Enums       []enumInfo
}

// enumInfo holds data for a specific enum type.
type enumInfo struct {
	TypeName     string
	Instances    []instanceData
	CaseFormat   string
	GenerateVars bool
	MarshalField string
	StructFields []fieldInfo
}

// typeSpec holds information about a parsed type definition.
type typeSpec struct {
	PackageName string
	TypeSpec    *ast.TypeSpec
	Doc         *ast.CommentGroup
	Fields      []fieldInfo
}

// typeDeclaration holds the AST nodes for a type declaration.
type typeDeclaration struct {
	file    *ast.File
	genDecl *ast.GenDecl
	spec    *ast.TypeSpec
}

// fieldInfo holds information about a struct field.
type fieldInfo struct {
	Name string
	Type string
}

// instanceResolution holds the result of resolving enum instances.
type instanceResolution struct {
	Instances    []instanceData
	GenerateVars bool
}

// instanceData holds information about each constant instance.
type instanceData struct {
	Name   string
	Fields map[string]string
}
