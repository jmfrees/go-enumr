package enumr

import (
	"go/ast"
)

// EnumData is used to pass the necessary data to the template.
type EnumData struct {
	PackageName string
	Enums       []EnumInfo
}

// EnumInfo holds data for a specific enum type.
type EnumInfo struct {
	TypeName     string
	Instances    []InstanceData
	CaseFormat   string
	GenerateVars bool
	NameField    string
	StructFields []FieldInfo
}

// TypeSpec holds information about a parsed type definition.
type TypeSpec struct {
	PackageName string
	TypeSpec    *ast.TypeSpec
	Doc         *ast.CommentGroup
	Fields      []FieldInfo
}

// typeDeclaration holds the AST nodes for a type declaration.
type typeDeclaration struct {
	file    *ast.File
	genDecl *ast.GenDecl
	spec    *ast.TypeSpec
}

// FieldInfo holds information about a struct field.
type FieldInfo struct {
	Name string
	Type string
}

// InstanceResolution holds the result of resolving enum instances.
type InstanceResolution struct {
	Instances    []InstanceData
	GenerateVars bool
}

// InstanceData holds information about each constant instance.
type InstanceData struct {
	Name   string
	Fields map[string]string
}
