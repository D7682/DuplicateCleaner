package config

import (
	"fmt"
	"github.com/spf13/viper"
	"runtime"
	"time"
)

// Config struct holds configuration parameters for the duplicate cleaner application
type Config struct {
	SrcFolder      string        `mapstructure:"src_folder"`          // SrcFolder to scan for duplicates (required)
	BkpFolder      string        `mapstructure:"bkp_folder"`          // BkpFolder to specify the path to the backup folder where non-duplicates will be saved. (required)
	MaxScanDepth   int           `mapstructure:"max_scan_depth"`      // MaxScanDepth specifies the maximum depth level to scan recursively. (optional, default: infinite)
	ExclExt        []string      `mapstructure:"excluded_extensions"` // ExclExt specifies the list of file extensions to exclude from scanning. (optional, default: none)
	ConcurrentScan bool          `mapstructure:"concurrent_scan"`     // ConcurrentScan enables concurrency for faster scanning. (optional, default: true)
	MaxConcurrent  int           `mapstructure:"max_concurrent"`      // MaxConcurrent specifies the maximum number of concurrent scans. (optional, default: runtime.NumCPU())
	DuplicateThr   time.Duration `mapstructure:"duplicate_threshold"` // DuplicateThr specifies the time threshold for considering files as duplicates. (optional, default: 1 hour)
	DryRun         bool          `mapstructure:"dry_run"`             // DryRun enables dry-run mode, displaying duplicates without removing them. (optional, default: false)
	LogFilePath    string        `mapstructure:"log_file_path"`       // LogFilePath specifies the path to the log file to save the scanning results. (optional)
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
	if config.MaxConcurrent <= 0 {
		// Set default value based on available CPUs
		config.MaxConcurrent = runtime.NumCPU()
	}

	// You can add validation logic for other fields here

	return config, nil
}
