package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type VersionInfo struct {
	Pkg     string `json:"pkg"`
	Version string `json:"version"`
}

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	versionFile := filepath.Join(cwd, "version.json")

	data, err := os.ReadFile(versionFile)
	if err != nil {
		log.Fatalf("Failed to read version.json (%s): %v", versionFile, err)
	}

	var v VersionInfo
	if err := json.Unmarshal(data, &v); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	tag := fmt.Sprintf("%s/v%s", v.Pkg, v.Version)
	fmt.Printf("==> Tag to create: %s\n", tag)

	exists, err := tagExists(tag)
	if err != nil {
		log.Fatalf("Failed to check tag: %v", err)
	}
	if exists {
		fmt.Printf("⚠️ Tag %s already exists, skipping.\n", tag)
		return
	}

	if err := runCmd("git", "tag", tag); err != nil {
		log.Fatalf("Failed to create tag: %v", err)
	}

	if err := runCmd("git", "push", "origin", tag); err != nil {
		log.Fatalf("Failed to push tag: %v", err)
	}

	fmt.Println("✅ Done!")
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func tagExists(tag string) (bool, error) {
	cmd := exec.Command("git", "tag", "--list", tag)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == tag, nil
}
