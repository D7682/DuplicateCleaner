// main.go

package main

import (
	"DuplicateCleaner/internal/config"
	"DuplicateCleaner/internal/folder"
	"fmt"
	"github.com/fatih/color"
	"log"
	"time"
)

var (
	pastelGreen = color.New(color.FgHiGreen).SprintFunc()
	pastelBlue  = color.New(color.FgHiBlue).SprintFunc()
	pastelPink  = color.New(color.FgHiMagenta).SprintFunc()
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	sourceFolder := cfg.SourceFolder
	backupFolder := cfg.BackupFolder

	// Create a folder instance
	parentFolder := folder.NewFolder(sourceFolder)

	// Scan with concurrency
	startConcurrent := time.Now()
	err = parentFolder.ScanWithConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	displayResults("Scan with Concurrency", sourceFolder, backupFolder, cfg, parentFolder)
	fmt.Printf("Time Elapsed (Concurrency): %s\n", time.Since(startConcurrent))

	// Scan without concurrency
	startSequential := time.Now()
	err = parentFolder.ScanWithoutConcurrency(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	displayResults("Scan without Concurrency", sourceFolder, backupFolder, cfg, parentFolder)
	fmt.Printf("Time Elapsed (Sequential): %s\n", time.Since(startSequential))
}

func displayResults(title, sourceFolder, backupFolder string, cfg config.Config, parentFolder *folder.Folder) {
	// Display results using the logger with color
	color.White("%s:\n", title)
	color.White("Source Folder: %v\n", sourceFolder)
	color.White("Backup Folder: %v\n", backupFolder)
	color.White("Max Scan Depth: %v\n", cfg.MaxScanDepth)

	// Display file information using the logger with color
	color.White("Files in the folder:\n")
	for _, file := range parentFolder.Files {
		color.Cyan("File: %s\n", file.Path)
		color.Yellow("MD5 Checksum: %s\n", file.MD5Checksum)
		// Add more fields as needed
		color.White("------------------------\n")
	}

	// Display file count using pastel colors
	fileCount := len(parentFolder.Files)
	fmt.Printf("File Count: %s\n", pastelGreen(fileCount))

	// Backup to another folder
	parentFolder.BackupToFolder(backupFolder)
}
