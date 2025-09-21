package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetModuleName(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open go.mod file [%s]: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			module := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if module == "" {
				return "", fmt.Errorf("invalid go.mod: module line is empty")
			}
			return module, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error while reading go.mod file [%s]: %w", path, err)
	}

	return "", fmt.Errorf("no module declaration found in go.mod [%s]", path)
}
