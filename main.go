package main

import (
	"fmt"
	"os"

	"github.com/dhairya13703/docker-cleanup/docker"
	"github.com/dhairya13703/docker-cleanup/utils"
	"github.com/spf13/cobra"
)

var (
	dryRun    bool
	olderThan int
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "docker-cleanup",
	Short: "A CLI tool to cleanup unused Docker resources",
	Long: `docker-cleanup is a CLI tool that helps you maintain your Docker environment
by automatically removing unused containers, images, volumes, and networks.
It can be configured to run periodically and clean resources based on age.`,
}

var containersCmd = &cobra.Command{
	Use:   "containers",
	Short: "Cleanup stopped containers",
	Run: func(cmd *cobra.Command, args []string) {
		cleanupContainers()
	},
}

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Cleanup unused images",
	Run: func(cmd *cobra.Command, args []string) {
		cleanupImages()
	},
}

var volumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "Cleanup unused volumes",
	Run: func(cmd *cobra.Command, args []string) {
		cleanupVolumes()
	},
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Cleanup all unused Docker resources",
	Run: func(cmd *cobra.Command, args []string) {
		cleanupAll()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Print actions without executing them")
	rootCmd.PersistentFlags().IntVar(&olderThan, "older-than", 24, "Remove resources older than specified hours (default 24h)")

	rootCmd.AddCommand(containersCmd)
	rootCmd.AddCommand(imagesCmd)
	rootCmd.AddCommand(volumesCmd)
	rootCmd.AddCommand(allCmd)
}

func cleanupContainers() {
	if err := utils.IsAdmin(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := utils.IsDockerRunning(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cleanup, err := docker.NewCleanup(dryRun, olderThan)
	if err != nil {
		fmt.Printf("Failed to initialize cleanup: %v\n", err)
		os.Exit(1)
	}

	if err := cleanup.CleanContainers(); err != nil {
		fmt.Printf("Failed to cleanup containers: %v\n", err)
		os.Exit(1)
	}
}

func cleanupImages() {
	if err := utils.IsAdmin(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := utils.IsDockerRunning(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cleanup, err := docker.NewCleanup(dryRun, olderThan)
	if err != nil {
		fmt.Printf("Failed to initialize cleanup: %v\n", err)
		os.Exit(1)
	}

	if err := cleanup.CleanImages(); err != nil {
		fmt.Printf("Failed to cleanup images: %v\n", err)
		os.Exit(1)
	}
}

func cleanupVolumes() {
	if err := utils.IsAdmin(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := utils.IsDockerRunning(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cleanup, err := docker.NewCleanup(dryRun, olderThan)
	if err != nil {
		fmt.Printf("Failed to initialize cleanup: %v\n", err)
		os.Exit(1)
	}

	if err := cleanup.CleanVolumes(); err != nil {
		fmt.Printf("Failed to cleanup volumes: %v\n", err)
		os.Exit(1)
	}
}

func cleanupAll() {
	if err := utils.IsAdmin(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := utils.IsDockerRunning(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cleanup, err := docker.NewCleanup(dryRun, olderThan)
	if err != nil {
		fmt.Printf("Failed to initialize cleanup: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Cleaning up containers...")
	if err := cleanup.CleanContainers(); err != nil {
		fmt.Printf("Failed to cleanup containers: %v\n", err)
	}

	fmt.Println("\nCleaning up images...")
	if err := cleanup.CleanImages(); err != nil {
		fmt.Printf("Failed to cleanup images: %v\n", err)
	}

	fmt.Println("\nCleaning up volumes...")
	if err := cleanup.CleanVolumes(); err != nil {
		fmt.Printf("Failed to cleanup volumes: %v\n", err)
	}
}
