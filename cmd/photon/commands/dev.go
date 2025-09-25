package commands

import (
	"fmt"
	"path/filepath"

	"github.com/smtdfc/photon/cmd/photon/domain"
	"github.com/smtdfc/photon/cmd/photon/helpers"
	"github.com/urfave/cli/v2"
)

func Dev(c *cli.Context) error {
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

	fmt.Printf("Starting application...\n")
	if _, err := helpers.SpawnCommand("go", []string{"run", entryPoint}, true); err != nil {
		return fmt.Errorf("Application [%s] failed to start: %v", appName, err)
	}

	return nil
}
