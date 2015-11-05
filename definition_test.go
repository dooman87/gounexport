package gounexport_test
import (
	"github.com/dooman87/gounexport"
"testing"
)

func TestGetDefinitionsFunc(t *testing.T) {
	unimportedpkg := pkg + "/testfunc"

	_, fset, info := parsePackage(unimportedpkg, t)
	defs := gounexport.GetDefinitions(info, fset)

	//Used, Unused, main
	if len(defs) != 3 {
		t.Errorf("expected 3 exported definitions, but found %d", len(defs))
	}
}

func TestGetDefinitionsUnimported(t *testing.T) {
	unimportedpkg := pkg + "/unimported"

	_, fset, info := parsePackage(unimportedpkg, t)
	defs := gounexport.GetDefinitions(info, fset)

	//NeverImported
	if len(defs) != 1 {
		t.Errorf("expected 1 exported definitions, but found %d", len(defs))
	}
}

