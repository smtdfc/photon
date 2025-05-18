package gen_api

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Export struct {
	PkgPath string
	PkgName string
	Decl    string
}

func main() {
	root := "pkg"
	output := "photon.go"
	modulePath := "github.com/smtdfc/photon"

	exports := map[string][]Export{}

	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}

		dir := filepath.Dir(path)
		importPath := fmt.Sprintf(`%s/%s`, modulePath, filepath.ToSlash(dir))

		for _, decl := range node.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok && ts.Name.IsExported() {
						exports[importPath] = append(exports[importPath], Export{
							PkgPath: importPath,
							PkgName: node.Name.Name,
							Decl:    fmt.Sprintf("type %s = %s.%s", ts.Name.Name, node.Name.Name, ts.Name.Name),
						})
					}
				}
			case *ast.FuncDecl:
				if d.Name.IsExported() && d.Recv == nil {
					exports[importPath] = append(exports[importPath], Export{
						PkgPath: importPath,
						PkgName: node.Name.Name,
						Decl:    fmt.Sprintf("var %s = %s.%s", d.Name.Name, node.Name.Name, d.Name.Name),
					})
				}
			}
		}
		return nil
	})

	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "package photon\n\nimport (\n")
	alias := map[string]string{}
	for path, exps := range exports {
		if len(exps) == 0 {
			continue
		}
		pkgName := exps[0].PkgName
		alias[pkgName] = path
		fmt.Fprintf(f, "\t%s \"%s\"\n", pkgName, path)
	}
	fmt.Fprintf(f, ")\n\n")

	for _, exps := range exports {
		for _, e := range exps {
			fmt.Fprintln(f, e.Decl)
		}
	}
}