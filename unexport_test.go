package gounexport_test

import (
	"go/ast"
	"go/types"
	"log"
	"regexp"
	"testing"

	"github.com/dooman87/gounexport"
	"github.com/dooman87/gounexport/util"
)

func TestGetDefinitionsToHideFunc(t *testing.T) {
	unimportedpkg := pkg + "/testfunc"
	unusedDefs := getDefinitionsToHide(unimportedpkg, 1, t)

	assertDef("github.com/dooman87/gounexport/testdata/testfunc.Unused", unusedDefs, t)
}

func TestGetDefinitionsToHideStruct(t *testing.T) {
	unimportedpkg := pkg + "/teststruct"
	unusedDefs := getDefinitionsToHide(unimportedpkg, 4, t)

	assertDef("github.com/dooman87/gounexport/testdata/teststruct.UnusedStruct", unusedDefs, t)
	assertDef("github.com/dooman87/gounexport/testdata/teststruct.UsedStruct.UnusedField", unusedDefs, t)
	assertDef("github.com/dooman87/gounexport/testdata/teststruct.UsedStruct.UnusedMethod", unusedDefs, t)
	assertDef("github.com/dooman87/gounexport/testdata/teststruct.UsedStruct.UsedInPackageMethod", unusedDefs, t)
}

func TestGetDefinitionsToHideVar(t *testing.T) {
	unimportedpkg := pkg + "/testvar"
	unusedDefs := getDefinitionsToHide(unimportedpkg, 2, t)

	assertDef("github.com/dooman87/gounexport/testdata/testvar.UnusedVar", unusedDefs, t)
	assertDef("github.com/dooman87/gounexport/testdata/testvar.UnusedConst", unusedDefs, t)
}

func TestGetDefinitionsToHideInterface(t *testing.T) {
	unimportedpkg := pkg + "/testinterface"
	unusedDefs := getDefinitionsToHide(unimportedpkg, 1, t)

	assertDef("github.com/dooman87/gounexport/testdata/testinterface.UnusedInterface", unusedDefs, t)
}

func TestGetDefinitionsToHideExclusions(t *testing.T) {
	unimportedpkg := pkg + "/testinterface"
	regex, _ := regexp.Compile("Unused*")
	excludes := []*regexp.Regexp{regex}
	getDefinitionsToHideWithExclusions(unimportedpkg, 0, excludes, t)
}

func getDefinitionsToHide(pkg string, expectedLen int, t *testing.T) []*gounexport.Definition {
	return getDefinitionsToHideWithExclusions(pkg, expectedLen, nil, t)
}

func getDefinitionsToHideWithExclusions(pkg string, expectedLen int, excludes []*regexp.Regexp, t *testing.T) []*gounexport.Definition {
	_, fset, info := parsePackage(pkg, t)
	defs := gounexport.GetDefinitions(info, fset)
	unusedDefs := gounexport.FindUnusedDefinitions(pkg, defs, excludes)

	if expectedLen > 0 && len(unusedDefs) != expectedLen {
		t.Errorf("expected %d unused exported definitions, but found %d", expectedLen, len(unusedDefs))
	}
	return unusedDefs
}

func TestGetDefinitionsToHideThis(t *testing.T) {
	pkg := "github.com/dooman87/gounexport"

	regex, _ := regexp.Compile("Test*")
	excludes := []*regexp.Regexp{regex}

	_, fset, info := parsePackage(pkg, t)
	defs := gounexport.GetDefinitions(info, fset)
	unusedDefs := gounexport.FindUnusedDefinitions(pkg, defs, excludes)

	log.Print("<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	for _, d := range unusedDefs {
		util.Info("DEFINITION %s", d.Name)
		util.Info("\t%s:%d:%d", d.File, d.Line, d.Col)
	}
	log.Print("<<<<<<<<<<<<<<<<<<<<<<<<<<<")

	if len(unusedDefs) != 22 {
		t.Errorf("expected %d unused exported definitions, but found %d", 22, len(unusedDefs))
	}
}

//ExampleGetUnusedDefitions shows how to use gounexport package
//to find all definition that not used in a package. As the result,
//all unused definitions will be printed in console.
func Example() {
	//package to check
	pkg := "github.com/dooman87/gounexport"

	//Regular expression to exclude
	//tests methods from the result.
	regex, _ := regexp.Compile("Test*")
	excludes := []*regexp.Regexp{regex}

	//Internal info structure that required for
	//ParsePackage call
	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	//Parsing package to fill info struct and
	//get file set.
	_, fset, err := gounexport.ParsePackage(pkg, &info)
	if err != nil {
		util.Err("error while parsing package %v", err)
	}

	//Analyze info and extract all definitions with usages.
	defs := gounexport.GetDefinitions(&info, fset)
	//Find all definitions that not used
	unusedDefs := gounexport.FindUnusedDefinitions(pkg, defs, excludes)
	//Print all unused definition to stdout.
	for _, d := range unusedDefs {
		util.Info("DEFINITION %s", d.Name)
	}
}

func assertDef(name string, defs []*gounexport.Definition, t *testing.T) {
	found := false
	for _, d := range defs {
		if name == d.Name {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected [%s] in Definitions", name)
	}
}

func TestUnexport(t *testing.T) {
	_, fset, info := parsePackage(pkg+"/testrename", t)
	defs := gounexport.GetDefinitions(info, fset)
	unusedDefs := gounexport.FindUnusedDefinitions(pkg, defs, nil)

	renamesCount := make(map[string]int)
	renameFunc := func(file string, offset int, source string, target string) error {
		log.Printf("renaming [%s] at %d from [%s] to [%s]", file, offset, source, target)
		renamesCount[source] = renamesCount[source] + 1
		return nil
	}

	for _, d := range unusedDefs {
		err := gounexport.Unexport(d, defs, renameFunc)
		if d.SimpleName == "UnusedStructConflict" && err == nil {
			t.Error("expected conflict error for UnusedStructConflict")
		}
		if d.SimpleName == "UnusedVarConflict" && err == nil {
			t.Error("expected conflict error for UnusedVarConflict")
		}
	}

	assertRename(renamesCount, "UnusedField", 1, t)
	assertRename(renamesCount, "UnusedMethod", 1, t)
	assertRename(renamesCount, "UsedInPackageMethod", 2, t)
}

func assertRename(renamesCount map[string]int, name string, expected int, t *testing.T) {
	if renamesCount[name] != expected {
		t.Errorf("expected [%d] renames of [%s], but was [%d]", expected, name, renamesCount[name])
	}
}
