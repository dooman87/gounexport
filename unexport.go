package gounexport

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dooman87/gounexport/fs"
	"github.com/dooman87/gounexport/util"
)

//FindUnusedDefinitions returns list of definitions that could be
//moved to private e.g. renamed. Criteria for renaming:
// - Definition should be exported
// - Definition should be in target package
// - Definition is not implementing external interfaces
// - Definition is not used in external packages
func FindUnusedDefinitions(pkg string, defs map[string]*Definition, excludes []*regexp.Regexp) []*Definition {
	var unused []*Definition
	for _, def := range defs {
		if !def.Exported {
			continue
		}

		if strings.HasPrefix(def.Name, pkg) && !isExcluded(def, excludes) && !isUsed(def) {
			util.Info("adding [%s] to unexport list", def.Name)
			unused = append(unused, def)
		}
	}

	return unused
}

func isExcluded(def *Definition, excludes []*regexp.Regexp) bool {
	if excludes == nil || len(excludes) == 0 {
		return false
	}

	for _, exc := range excludes {
		if exc.MatchString(def.Name) {
			util.Info("definition [%s] excluded, because matched [%s]", def.Name, exc.String())
			return true
		}
	}
	return false
}

func isUsed(def *Definition) bool {
	used := true

	if len(def.Usages) == 0 {
		used = false
	} else {
		//Checking pathes of usages to not count internal
		hasExternalUsages := false
		util.Debug("checking [%s]", def.Name)
		for _, u := range def.Usages {
			pkgPath := ""
			if def.Pkg != nil {
				pkgPath = def.Pkg.Path()
			} else if dotIdx := strings.LastIndex(def.Name, "."); dotIdx >= 0 {
				pkgPath = def.Name[0:dotIdx]
			}
			util.Debug("checking [%v]", u.Pos)
			if u.Pos.IsValid() && fs.GetPackagePath(u.Pos.Filename) != pkgPath {
				hasExternalUsages = true
				break
			}
		}
		used = hasExternalUsages
	}

	if !used {
		//Check all interfaces
		for _, i := range def.Interfaces {
			if isUsed(i) {
				used = true
				break
			}
		}
	}
	return used
}

//Unexport hides definition by changing first letter
//to lower case. It won't rename if there is already existing
//unexported symbol with the same name.
//renameFunc is a func that accepts four arguments: full path to file,
//offset in a file to replace, original string, string to replace. It will
//be called when renaming is possible.
func Unexport(def *Definition, allDefs map[string]*Definition,
	renameFunc func(string, int, string, string) error) error {
	util.Info("unexporting %s in %s:%d:%d", def.SimpleName, def.File, def.Line, def.Col)
	newName := strings.ToLower(def.SimpleName[0:1]) + def.SimpleName[1:]

	//Searching for conflict
	lastIdx := strings.LastIndex(def.Name, def.SimpleName)
	newFullName := def.Name[0:lastIdx] + newName + def.Name[lastIdx+len(newName):]
	if allDefs[newFullName] != nil {
		return fmt.Errorf("can't unexport %s because it conflicts with existing member", def.Name)
	}

	//rename definitions and usages
	err := renameFunc(def.File, def.Offset, def.SimpleName, newName)
	for _, u := range def.Usages {
		if err != nil {
			break
		}
		err = renameFunc(u.Pos.Filename, u.Pos.Offset, def.SimpleName, newName)
	}

	return err
}
