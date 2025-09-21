package helpers

import (
	"github.com/smtdfc/photon/dix"
	"os"
	"path"
)

func GenerateDixDI(rootDir string, configFile string) error {
	err := os.MkdirAll(path.Join(rootDir, "dix"), 0755)
	if err != nil {
		return err
	}
	code, err := dix.GenFromConfigFile(configFile)
	if err != nil {
		return err
	}

	content := []byte(code)
	err = os.WriteFile(path.Join(rootDir, "dix/main.go"), content, 0644)
	if err != nil {
		return err
	}

	return nil
}
