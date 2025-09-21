package commands

import (
	"fmt"
	"github.com/smtdfc/photon/cli/domain"
	"github.com/smtdfc/photon/cli/helpers"
	"github.com/urfave/cli/v2"
	"path/filepath"
)

func Build(c *cli.Context) error {
	fmt.Println("Reading configuration file...")
	cwd := helpers.GetCWD()
	configFilePath := filepath.Join(cwd, "photon.config.json")

	if !helpers.FileExists(configFilePath) {
		return fmt.Errorf("Configuration file not found at: %s\nMake sure you are running the command from the project root.", configFilePath)
	}

	config, err := helpers.LoadJSONFile[*domain.Config](configFilePath)
	if err != nil {
		return fmt.Errorf("Failed to read configuration file [%s]: %v\nCheck if the JSON syntax is valid and that you have the correct permissions.", configFilePath, err)
	}

	appName := config.Name
	entryPoint := config.EntryPoint
	dixConfig := config.DixConfig

	if appName == "" {
		return fmt.Errorf("Missing 'name' field in configuration file: %s", configFilePath)
	}
	if entryPoint == "" {
		return fmt.Errorf("Missing 'entryPoint' field in configuration file: %s", configFilePath)
	}

	if err := helpers.GenerateDixDI(cwd, dixConfig); err != nil {
		return fmt.Errorf("Failed to generate Dix dependency injection code: %v", err)
	}

	fmt.Printf("[@%s] Building application...\n", appName)
	if _, err := helpers.SpawnCommand("go", []string{"build", "-o", "bin/" + appName, entryPoint}, true); err != nil {
		return fmt.Errorf("Build failed for application [%s]: %v", appName, err)
	}

	fmt.Printf("Build completed successfully. Output: bin/%s\n", appName)
	return nil
}
