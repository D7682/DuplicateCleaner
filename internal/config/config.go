package config

import (
	"fmt"
	"github.com/spf13/viper"
	"runtime"
	"time"
)

// Config struct holds configuration parameters for the duplicate cleaner application
type Config struct {
	SourceFolder       string        `mapstructure:"source_folder"`        // Folder to scan for duplicate files (required)
	BackupFolder       string        `mapstructure:"backup_folder"`        // Backup folder where duplicate files will be moved (required)
	MaxScanDepth       int           `mapstructure:"max_scan_depth"`       // Maximum depth level to scan recursively (optional, default: infinite)
	ExcludedFileTypes  []string      `mapstructure:"excluded_file_types"`  // File types to exclude from scanning (optional, default: none)
	ConcurrentScan     bool          `mapstructure:"concurrent_scan"`      // Enable concurrency for faster scanning (optional, default: true)
	MaxConcurrentScans int           `mapstructure:"max_concurrent_scans"` // Maximum number of concurrent scans (optional, default: runtime.NumCPU())
	DuplicateThreshold time.Duration `mapstructure:"duplicate_threshold"`  // Time threshold for considering files as duplicates (optional, default: 1 hour)
	DryRun             bool          `mapstructure:"dry_run"`              // Whether to dry-run (display duplicates without removing them) (optional, default: false)
	// Add other fields as needed...
}

const (
	InfiniteScanDepth = -1
)

// LoadConfig loads configuration from the specified file path
func LoadConfig(filePath string) (Config, error) {
	var config Config
	viper.SetConfigFile(filePath) // Set the configuration file path

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read configuration file: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration: %v", err)
	}

	// Check and set default for MaxScanDepth
	if config.MaxScanDepth <= InfiniteScanDepth {
		config.MaxScanDepth = InfiniteScanDepth
	}

	// Check and set default for MaxConcurrentScans
	if config.MaxConcurrentScans <= 0 {
		// Set default value based on available CPUs
		config.MaxConcurrentScans = runtime.NumCPU()
	}

	// You can add validation logic for other fields here

	return config, nil
}
