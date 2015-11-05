package gounexport_test

import (
	"go/token"
	"testing"
)

func TestParsePackageFunc(t *testing.T) {
	_, fset, _ := parsePackage(pkg+"/testfunc", t)
	fileCounter := 0
	iterator := func(f *token.File) bool {
		fileCounter++
		return true
	}
	fset.Iterate(iterator)

	if fileCounter != 2 {
		t.Errorf("expected 2 files in result file set but found %d", fileCounter)
	}
}
