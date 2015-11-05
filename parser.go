package gounexport
import (
	"github.com/dooman87/gounexport/util"
	"github.com/dooman87/gounexport/fs"
	"go/token"
	"go/types"
	"github.com/dooman87/gounexport/importer"
)

//ParsePackage parses package and filling info structure.
//It's filling info about all internal packages even if they
//are not imported in the root package.
func ParsePackage(pkgName string, info *types.Info) (*types.Package, *token.FileSet, error) {
	collectImporter := new(importer.CollectInfoImporter)
	collectImporter.Info = info

	var resultPkg *types.Package
	var resultFset *token.FileSet
	parsedPackages := make(map[string]bool)

	notParsedPackage := pkgName
	for len(notParsedPackage) > 0 {
		collectImporter.Pkg = notParsedPackage
		pkg, fset, err := collectImporter.Collect()
		if err != nil {
			return nil, nil, err
		}

		//Filling results only from first package
		//that was passed as argument to function
		if resultPkg == nil {
			resultPkg = pkg
			resultFset = fset
		}
		parsedPackages[notParsedPackage] = true

		//Searching for a new package that was not parsed before
		notParsedPackage = ""
		files, err := fs.GetUnusedSources(pkgName, fset)
		if err != nil {
			return nil, nil, err
		}
		for _, f := range files {
			newNotParsedPackage := fs.GetPackagePath(f)
			if !parsedPackages[newNotParsedPackage] {
				notParsedPackage = newNotParsedPackage
				break
			} else {
				util.Info("package %s has been already parsed, however %s file is still unused", newNotParsedPackage, f)
			}
		}
	}

	return resultPkg, resultFset, nil
}

