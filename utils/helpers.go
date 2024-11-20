package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// IsAdmin checks if the program is running with administrator privileges
func IsAdmin() error {
	switch runtime.GOOS {
	case "windows":
		// Check if running as Administrator on Windows
		_, err := exec.Command("net", "session").Output()
		if err != nil {
			return fmt.Errorf("this program must be run as Administrator on Windows")
		}
	case "linux", "darwin":
		// Check if running as root on Linux/macOS
		if os.Geteuid() != 0 {
			if runtime.GOOS == "darwin" {
				return fmt.Errorf("this program must be run with sudo on macOS")
			}
			return fmt.Errorf("this program must be run with sudo on Linux")
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	return nil
}

// IsDockerRunning checks if Docker daemon is running
func IsDockerRunning() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("docker", "info")
		// Set up specific Windows options if needed
	case "darwin":
		cmd = exec.Command("docker", "info")
		// For Docker Desktop on macOS
	case "linux":
		cmd = exec.Command("docker", "info")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		switch runtime.GOOS {
		case "windows":
			return fmt.Errorf("docker Desktop is not running on Windows. Please start Docker Desktop")
		case "darwin":
			return fmt.Errorf("docker Desktop is not running on macOS. Please start Docker Desktop")
		default:
			return fmt.Errorf("docker daemon is not running. Please start the Docker service")
		}
	}

	return nil
}
