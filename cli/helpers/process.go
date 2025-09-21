package helpers

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func SpawnCommand(name string, args []string, inheritIO bool) (int, error) {
	cmd := exec.Command(name, args...)

	if inheritIO {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	} else {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return -1, fmt.Errorf("failed to get stdout pipe: %w", err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return -1, fmt.Errorf("failed to get stderr pipe: %w", err)
		}

		// Stream outputs asynchronously
		go func() { _, _ = io.Copy(os.Stdout, stdout) }()
		go func() { _, _ = io.Copy(os.Stderr, stderr) }()
	}

	if err := cmd.Start(); err != nil {
		return -1, fmt.Errorf("cannot spawn command %q: %w", name, err)
	}

	err := cmd.Wait()
	if exitError, ok := err.(*exec.ExitError); ok {
		// Process ran but exited with non-zero
		return exitError.ExitCode(), err
	} else if err != nil {
		// Other failure (e.g. process start failed, I/O issues)
		return -1, fmt.Errorf("command %q failed: %w", name, err)
	}

	return 0, nil
}
