package gounexport_test
import (
	"testing"
	"go/types"
	"go/token"
	"go/ast"
	"github.com/dooman87/gounexport"
)

const (
	pkg = "github.com/dooman87/gounexport/testdata"
)

func parsePackage(pkgStr string, t *testing.T) (*types.Package, *token.FileSet, *types.Info) {
	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	packag, fset, err := gounexport.ParsePackage(pkgStr, &info)
	if err != nil {
		t.Errorf("error while parsing package %v", err)
	}
	return packag, fset, &info
}
