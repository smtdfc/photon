package photonGenerator

import (
	"os"
	"path/filepath"
	"text/template"
)

type ModuleTemplateData struct {
	PackageName string
	ModuleName  string
}

func GenerateFile(path string, tmplPath string, data any) error {
	absTmplPath := filepath.Join("../", "templates", tmplPath)

	tmplContent, err := os.ReadFile(absTmplPath)
	if err != nil {
		return err
	}

	t, err := template.New(filepath.Base(tmplPath)).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, data)
}