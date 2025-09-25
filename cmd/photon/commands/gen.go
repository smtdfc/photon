package commands

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/smtdfc/photon/cmd/photon/domain"
	"github.com/smtdfc/photon/cmd/photon/helpers"
	"github.com/smtdfc/photon/dix"
	"github.com/urfave/cli/v2"
)

func GenModule(config *domain.Config, appName, moduleName, pkgName, projPath string) error {
	fmt.Printf("[@%s] Generating module '%s'...\n", appName, moduleName)

	if !helpers.IsPascalCase(moduleName) {
		return fmt.Errorf("[@%s] Invalid module name: '%s'. Module names must follow CamelCase format.", appName, moduleName)
	}

	callDir := GetCallerDir()
	normalizedName := strings.ToLower(moduleName)
	appNormalizedName := strings.ToLower(appName)

	projectFilesMap := map[string]string{
		"modules/" + normalizedName + "/gen.go":       "../templates/module/gen.go.tmpl",
		"modules/" + normalizedName + "/lifecycle.go": "../templates/module/lifecycle.go.tmpl",
		"modules/" + normalizedName + "/http.go":      "../templates/module/http.go.tmpl",
		"modules/" + normalizedName + "/dix.json":     "../templates/module/dix.json.tmpl",
		"test/modules/" + normalizedName + "_test.go": "../templates/test/module_test.go.tmpl",
	}

	data := map[string]any{
		"PkgName":              pkgName,
		"ModuleNormalizedName": normalizedName,
		"ModuleName":           moduleName,
		"AppNormaizedName":     appNormalizedName,
	}

	for fileName, tmplFile := range projectFilesMap {
		realFilePath, err := helpers.EnsureDirAndResolve(filepath.Join(projPath, fileName))
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

	dixConfig, err := dix.ReadConfigFile(config.DixConfig)
	if err != nil {
		return fmt.Errorf("Failed to read Dix configuration [%s]: %v", config.DixConfig, err)
	}

	dixConfig.Imports = append(dixConfig.Imports, "./modules/"+normalizedName+"/dix.json")
	dataPretty, err := json.MarshalIndent(dixConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to serialize Dix configuration: %v", err)
	}

	if err := helpers.WriteFile(filepath.Join(projPath, config.DixConfig), string(dataPretty)); err != nil {
		return fmt.Errorf("Failed to update Dix configuration [%s]: %v", config.DixConfig, err)
	}

	fmt.Printf("[@%s] Module '%s' created successfully.\n", appName, moduleName)
	return nil
}

func Gen(c *cli.Context) error {
	if c.NArg() < 2 {
		return fmt.Errorf("Invalid arguments. Usage: photon gen <type> <name>")
	}

	cwd := helpers.GetCWD()

	fmt.Println("Reading go.mod...")
	goModFilePath := filepath.Join(cwd, "go.mod")
	if !helpers.FileExists(goModFilePath) {
		return fmt.Errorf("go.mod file not found at: %s", goModFilePath)
	}

	fmt.Println("Reading configuration file...")
	configFilePath := filepath.Join(cwd, "photon.config.json")
	if !helpers.FileExists(configFilePath) {
		return fmt.Errorf("Configuration file not found at: %s", configFilePath)
	}

	config, err := helpers.LoadJSONFile[*domain.Config](configFilePath)
	if err != nil {
		return fmt.Errorf("Failed to read configuration file [%s]: %v", configFilePath, err)
	}

	appName := config.Name
	targetType := c.Args().Get(0)
	targetName := c.Args().Get(1)

	moduleName, err := helpers.GetModuleName(goModFilePath)
	if err != nil || moduleName == "" {
		// fallback to app name if go.mod module cannot be determined
		moduleName = appName
	}

	switch targetType {
	case "module":
		return GenModule(config, appName, targetName, moduleName, cwd)
	default:
		return fmt.Errorf("Unknown target type: %s. Supported types: module", targetType)
	}
}
