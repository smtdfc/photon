package dix

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
)

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func ReadConfigFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func MergeConfig(
	rootDir string,
	source *Config,
	target *Config,
) {
	for name, provider := range source.Providers {
		provider.From = path.Join(rootDir, source.Module, provider.From)
		source.Providers[name] = provider
	}

	if target.Providers != nil {
		for name, provider := range source.Providers {
			target.Providers[name] = provider
		}
	} else {
		target.Providers = source.Providers
	}
}

func GenFromConfigFile(path string) (string, error) {
	config, err := ReadConfigFile(path)
	if err != nil {
		return "", err
	}

	if config.Imports != nil {
		for _, filePath := range config.Imports {

			subConfig, err := ReadConfigFile(filePath)

			if err != nil {
				return "", err
			}
			MergeConfig(config.Pkg, subConfig, config)
		}
	}

	generator := NewGenerator(config)
	code, err := generator.Generate()
	if err != nil {
		return "", err
	}
	return code, nil
}
