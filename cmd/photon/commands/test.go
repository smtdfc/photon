package commands

import (
	"fmt"
	"path/filepath"

	"github.com/smtdfc/photon/cmd/photon/domain"
	"github.com/smtdfc/photon/cmd/photon/helpers"
	"github.com/urfave/cli/v2"
)

func Test(c *cli.Context) error {
	fmt.Println("Reading configuration file...")
	cwd := helpers.GetCWD()
	configFilePath := filepath.Join(cwd, "photon.config.json")

	if !helpers.FileExists(configFilePath) {
		return fmt.Errorf("Configuration file not found at: %s\nMake sure you are running the command from the project root.", configFilePath)
	}

	config, err := helpers.LoadJSONFile[*domain.Config](configFilePath)
	if err != nil {
		return fmt.Errorf("Failed to read configuration file [%s]: %v", configFilePath, err)
	}

	appName := config.Name
	if appName == "" {
		return fmt.Errorf("Missing 'name' field in configuration file: %s", configFilePath)
	}

	fmt.Printf("[@%s] Running tests...\n", appName)
	if _, err := helpers.SpawnCommand("go", []string{"test", "./test/..."}, true); err != nil {
		return fmt.Errorf("Tests failed for application [%s]: %v", appName, err)
	}

	fmt.Printf("[@%s] Tests completed successfully.\n", appName)
	return nil
}
