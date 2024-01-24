package main

import (
	"DuplicateCleaner/internal/config"
	"DuplicateCleaner/internal/folder"
	"DuplicateCleaner/logger"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	sourceFolder := cfg.SrcFolder
	// backupFolder := cfg.BkpFolder

	// Create a folder instance
	parentFolder := folder.NewFolder(sourceFolder, &cfg)

	// Initialize logger with both console and file outputs
	logFilePath := cfg.LogFilePath
	lgr := logger.NewLogger(logFilePath)
	defer lgr.Close()

	// Scan with concurrency
	err = parentFolder.ScanWithConfig()
	if err != nil {
		lgr.Fatal(err)
	}

	// Scan without concurrency
	err = parentFolder.ScanWithoutConcurrency()
	if err != nil {
		lgr.Fatal(err)
	}
}
