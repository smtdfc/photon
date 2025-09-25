package dix

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
)

type ImportManager struct {
	Sources map[string]string
}

var counter = 0

func GenUniqueName() string {
	counter++
	return fmt.Sprintf("id%d", counter)
}

func (m *ImportManager) Import(path string) string {
	if _, ok := m.Sources[path]; !ok {
		m.Sources[path] = `import_` + GenUniqueName()
	}
	return m.Sources[path]
}

func (m *ImportManager) GenImportDecls() *ast.GenDecl {
	decls := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{},
	}

	for path, alias := range m.Sources {
		decls.Specs = append(
			decls.Specs,
			&ast.ImportSpec{
				Name: ast.NewIdent(alias),
				Path: &ast.BasicLit{Kind: token.STRING, Value: "\"" + path + "\""},
			},
		)
	}
	return decls
}

func NewImportManager() *ImportManager {
	return &ImportManager{
		Sources: make(map[string]string),
	}
}

type Generator struct {
	Graph         DependencyGraph
	Config        *Config
	Singleton     map[string]string
	ImportManager *ImportManager
	Stmts         []ast.Stmt // all statements generated
}

func NewGenerator(config *Config) *Generator {
	return &Generator{
		Graph:         BuildDepGraph(config),
		Config:        config,
		Singleton:     make(map[string]string),
		ImportManager: NewImportManager(),
		Stmts:         []ast.Stmt{},
	}
}

func isDepOfOther(name string, graph DependencyGraph) bool {
	for parent, deps := range graph {
		if parent == name {
			continue
		}
		for _, d := range deps {
			if d.Name == name {
				return true
			}
		}
	}
	return false
}

func (g *Generator) GenExprForDep(dep *Dependency) ast.Expr {
	dInfo := g.Config.Providers[dep.Name]
	pkgAlias := g.ImportManager.Import(dInfo.From)

	args := []ast.Expr{}
	for _, sub := range g.Graph[dep.Name] {
		args = append(args, g.GenExprForDep(sub))
	}

	call := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent(pkgAlias),
			Sel: ast.NewIdent(dInfo.Factory),
		},
		Args: args,
	}

	if dep.Transient {
		// No cache
		return call
	}

	// Reuse instance
	if id, ok := g.Singleton[dep.Name]; ok {
		return ast.NewIdent(id)
	}

	// create assign statement

	id := GenUniqueName()

	g.Singleton[dep.Name] = id
	assign := &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent(id)},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{call},
	}
	g.Stmts = append(g.Stmts, assign)

	return ast.NewIdent(id)
}

func (g *Generator) Generate() (string, error) {
	order, err := BuildInitOrder(g.Graph)
	if err != nil {
		return "", err
	}

	if g.Config.Providers["Root"] == nil {
		return "", errors.New("cannot find root provider")
	}

	// clear statements
	g.Stmts = []ast.Stmt{}

	for _, dep := range order {
		_ = g.GenExprForDep(dep)
	}

	for name, id := range g.Singleton {
		if name == "Root" {
			continue
		}
		if !isDepOfOther(name, g.Graph) {
			g.Stmts = append(g.Stmts, &ast.AssignStmt{
				Lhs: []ast.Expr{ast.NewIdent("_")},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{ast.NewIdent(id)},
			})
		}
	}

	// return Root
	g.Stmts = append(g.Stmts, &ast.ReturnStmt{
		Results: []ast.Expr{ast.NewIdent(g.Singleton["Root"])},
	})

	// build ast.File
	file := &ast.File{
		Name: ast.NewIdent("dix"),
		Decls: []ast.Decl{
			g.ImportManager.GenImportDecls(),
			&ast.FuncDecl{
				Name: ast.NewIdent("Init"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("any"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{List: g.Stmts},
			},
		},
	}

	fset := token.NewFileSet()
	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, file)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
