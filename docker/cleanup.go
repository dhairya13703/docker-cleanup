package docker

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Cleanup struct {
	dryRun    bool
	olderThan int
}

func NewCleanup(dryRun bool, olderThan int) (*Cleanup, error) {
	return &Cleanup{
		dryRun:    dryRun,
		olderThan: olderThan,
	}, nil
}

func executeCommand(cmd *exec.Cmd) (string, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %v\nOutput: %s", err, string(output))
	}
	return string(output), nil
}

func (c *Cleanup) CleanContainers() error {
	// Get all stopped containers
	cmd := exec.Command("docker", "ps", "-a", "--filter", "status=exited", "--filter", "status=dead", "--format", "{{.ID}}\t{{.State}}")
	output, err := executeCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to list containers: %v", err)
	}

	containers := strings.Split(strings.TrimSpace(output), "\n")
	if len(containers) == 0 || (len(containers) == 1 && containers[0] == "") {
		fmt.Println("No stopped containers found")
		return nil
	}

	for _, container := range containers {
		if container == "" {
			continue
		}

		parts := strings.Split(container, "\t")
		containerID := parts[0]

		// Get container creation time
		timeCmd := exec.Command("docker", "inspect", "-f", "{{.State.FinishedAt}}", containerID)
		timeOutput, err := executeCommand(timeCmd)
		if err != nil {
			fmt.Printf("Failed to get finish time for container %s: %v\n", containerID, err)
			continue
		}

		finishTime, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(timeOutput))
		if err != nil {
			fmt.Printf("Failed to parse finish time for container %s: %v\n", containerID, err)
			continue
		}

		if time.Since(finishTime).Hours() > float64(c.olderThan) {
			if c.dryRun {
				fmt.Printf("[DRY RUN] Would remove container: %s\n", containerID)
				continue
			}

			removeCmd := exec.Command("docker", "rm", "-f", "-v", containerID)
			if out, err := executeCommand(removeCmd); err != nil {
				fmt.Printf("Failed to remove container %s: %v\n", containerID, err)
			} else {
				fmt.Print(out)
			}
		}
	}

	return nil
}

func (c *Cleanup) CleanImages() error {
	// First, get a list of all images with their details
	cmd := exec.Command("docker", "images", "--format", "{{.ID}}\t{{.Repository}}\t{{.Tag}}\t{{.CreatedAt}}")
	output, err := executeCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to list images: %v", err)
	}

	// Get list of images used by containers (both running and stopped)
	usedCmd := exec.Command("docker", "ps", "-a", "--format", "{{.Image}}")
	usedOutput, err := executeCommand(usedCmd)
	if err != nil {
		return fmt.Errorf("failed to get used images: %v", err)
	}

	usedImages := make(map[string]bool)
	for _, img := range strings.Split(strings.TrimSpace(usedOutput), "\n") {
		if img != "" {
			usedImages[img] = true
		}
	}

	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 4 {
			continue
		}

		imageID := parts[0]
		repo := parts[1]
		tag := parts[2]
		createdStr := parts[3]

		// Parse the creation time
		created, err := time.Parse("2006-01-02 15:04:05 -0700 MST", createdStr)
		if err != nil {
			fmt.Printf("Failed to parse creation time for image %s: %v\n", imageID, err)
			continue
		}

		imageName := fmt.Sprintf("%s:%s", repo, tag)
		isUsed := usedImages[imageName]

		if !isUsed && time.Since(created).Hours() > float64(c.olderThan) {
			if c.dryRun {
				fmt.Printf("[DRY RUN] Would remove image: %s (%s)\n", imageName, imageID[:12])
				continue
			}

			removeCmd := exec.Command("docker", "rmi", "-f", imageID)
			if out, err := executeCommand(removeCmd); err != nil {
				fmt.Printf("Failed to remove image %s: %v\n", imageID[:12], err)
			} else {
				fmt.Print(out)
			}
		}
	}

	// Clean up any remaining dangling images
	if !c.dryRun {
		pruneCmd := exec.Command("docker", "image", "prune", "-f")
		if out, err := executeCommand(pruneCmd); err != nil {
			fmt.Printf("Failed to remove dangling images: %v\n", err)
		} else {
			fmt.Print(out)
		}
	}

	return nil
}

func (c *Cleanup) CleanVolumes() error {
	// First list volumes with their details
	cmd := exec.Command("docker", "volume", "ls", "--format", "{{.Name}}")
	output, err := executeCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to list volumes: %v", err)
	}

	// Get list of volumes used by containers
	usedCmd := exec.Command("docker", "ps", "-a", "--format", "{{.Mounts}}")
	usedOutput, err := executeCommand(usedCmd)
	if err != nil {
		return fmt.Errorf("failed to get used volumes: %v", err)
	}

	usedVolumes := make(map[string]bool)
	for _, mounts := range strings.Split(strings.TrimSpace(usedOutput), "\n") {
		for _, mount := range strings.Split(mounts, ",") {
			if strings.Contains(mount, "volume") {
				parts := strings.Fields(mount)
				if len(parts) > 0 {
					usedVolumes[parts[len(parts)-1]] = true
				}
			}
		}
	}

	for _, volume := range strings.Split(strings.TrimSpace(output), "\n") {
		if volume == "" {
			continue
		}

		if !usedVolumes[volume] {
			if c.dryRun {
				fmt.Printf("[DRY RUN] Would remove volume: %s\n", volume)
				continue
			}

			removeCmd := exec.Command("docker", "volume", "rm", "-f", volume)
			if out, err := executeCommand(removeCmd); err != nil {
				fmt.Printf("Failed to remove volume %s: %v\n", volume, err)
			} else {
				fmt.Print(out)
			}
		}
	}

	return nil
}

func (c *Cleanup) CleanAll() error {
	fmt.Println("Cleaning up containers...")
	if err := c.CleanContainers(); err != nil {
		fmt.Printf("Error cleaning containers: %v\n", err)
	}

	fmt.Println("\nCleaning up images...")
	if err := c.CleanImages(); err != nil {
		fmt.Printf("Error cleaning images: %v\n", err)
	}

	fmt.Println("\nCleaning up volumes...")
	if err := c.CleanVolumes(); err != nil {
		fmt.Printf("Error cleaning volumes: %v\n", err)
	}

	return nil
}
