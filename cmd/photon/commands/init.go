package commands

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/smtdfc/photon/cmd/photon/helpers"
	"github.com/urfave/cli/v2"
)

func GetCallerDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Dir(filename)
}

func Init(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("Missing project name. Usage: photon init <project-name>")
	}

	cwd := helpers.GetCWD()

	fmt.Println("Reading go.mod...")
	goModFilePath := filepath.Join(cwd, "go.mod")
	if !helpers.FileExists(goModFilePath) {
		return fmt.Errorf("go.mod file not found at: %s\nMake sure you run 'go mod init <module>' before initializing a Photon project.", goModFilePath)
	}

	projectName := c.Args().Get(0)
	fmt.Printf("Initializing project '%s'...\n", projectName)

	pkgName, err := helpers.GetModuleName(goModFilePath)
	if err != nil {
		return fmt.Errorf("Failed to read module name from go.mod: %v", err)
	}

	callDir := GetCallerDir()
	projectFilesMap := map[string]string{
		"photon.config.json": "../templates/project/photon.config.json.tmpl",
		"main.go":            "../templates/project/main.go.tmpl",
		"app/app.go":         "../templates/project/app/app.go.tmpl",
		"dix.json":           "../templates/project/dix.json.tmpl",
		"app/dix.json":       "../templates/project/app/dix.json.tmpl",
	}

	data := map[string]any{
		"PkgName":     pkgName,
		"ProjectName": projectName,
	}

	for fileName, tmplFile := range projectFilesMap {
		realFilePath, err := helpers.EnsureDirAndResolve(filepath.Join(cwd, fileName))
		if err != nil {
			return fmt.Errorf("Failed to prepare path for %s: %v", fileName, err)
		}

		realTemplPath := filepath.Join(callDir, tmplFile)
		code, err := helpers.RenderTemplateFile(realTemplPath, data)
		if err != nil {
			return fmt.Errorf("Failed to render template %s: %v", tmplFile, err)
		}

		if err := helpers.WriteFile(realFilePath, code); err != nil {
			return fmt.Errorf("Failed to write file %s: %v", realFilePath, err)
		}
	}

	fmt.Println("Project created successfully.")
	return nil
}
