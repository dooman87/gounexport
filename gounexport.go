//Package gounexport provides functionality to unexport unused public symbols.
//
//In detail, what you can do using this package:
//
//* parse package and get package's definition
//
//* get information about all definitions in packages
//
//* get unused packages
//
//* unexport unused definitions
//
//For result Definition struct is used. It's also includes Definition.Usage array
//with all usages (internal and external) across the package.
//
package gounexport

import (
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
)

//Definition of symbol in package
type Definition struct {
	//Full file path for current defintion
	File string
	//Full name of the definition
	Name string
	//Simple name of the definition. We are using it
	//for renaming purposes.
	SimpleName string
	//List of interfaces that implemented current definition
	//It will be interfaces definitions for a type definition and
	// methods definition for a function
	Interfaces []*Definition
	//type of definition
	TypeOf reflect.Type
	//Number of line in the file where definition is declared
	Line int
	//Column in the file where definition is declared
	Col int
	//Offset in source file
	Offset int
	//True, if definition is exported
	Exported bool
	//Package of the definition
	Pkg *types.Package
	//List of usages of the definition
	Usages []*Usage
}

func (def *Definition) addUsage(pos token.Position) {
	u := new(Usage)
	u.Pos = pos
	def.Usages = append(def.Usages, u)
}

//Usage is a struct that define one usage of a definition
type Usage struct {
	//Pos is a position of usage: file, line, col
	Pos token.Position
}

type objectWithIdent struct {
	obj   types.Object
	ident *ast.Ident
}

func newObjectWithIdent(obj types.Object, ident *ast.Ident) *objectWithIdent {
	result := new(objectWithIdent)
	result.obj = obj
	result.ident = ident
	return result
}

type defWithInterface struct {
	def      *Definition
	interfac *types.Interface
	named    *types.Named
}

type context struct {
	structs    map[string]string
	vars       []*objectWithIdent
	funcs      []*objectWithIdent
	interfaces []*defWithInterface
	fset       *token.FileSet
	defs       map[string]*Definition
}

func newContext(fset *token.FileSet) *context {
	ctx := new(context)
	ctx.fset = fset
	ctx.structs = make(map[string]string, 0)
	ctx.interfaces = make([]*defWithInterface, 0)
	ctx.vars = make([]*objectWithIdent, 0)
	ctx.funcs = make([]*objectWithIdent, 0)
	return ctx
}
