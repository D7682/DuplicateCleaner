package folder

import (
	"DuplicateCleaner/internal/config"
	"DuplicateCleaner/logger"
	"DuplicateCleaner/pkg/fileutil"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Folder represents a folder and its related operations.
type Folder struct {
	Path   string
	Files  []fileutil.File // Slice to hold File instances
	Logger *logger.Logger  // Logger instance
}

// NewFolder creates a new Folder instance with the given path and logger.
func NewFolder(path string) *Folder {
	return &Folder{
		Path:   path,
		Logger: logger.NewLogger(),
	}
}

// ScanWithConfig scans the folder concurrently based on the provided config.
func (f *Folder) ScanWithConfig(config *config.Config) error {
	// Implementation for scanning
	f.Logger.Color.Printf("Scanning folder: %s\n", f.Path)

	// Clear existing files
	f.Files = nil

	// Channel to receive file paths from goroutines
	filePaths := make(chan string, config.MaxConcurrentScans)

	// WaitGroup to wait for all worker goroutines to complete
	var wg sync.WaitGroup

	// WaitGroup for walking goroutine
	var walkWG sync.WaitGroup

	// Mutex to synchronize access to the f.Files slice
	var mu sync.Mutex

	// Start a goroutine to close the channel when all workers are done
	go func() {
		// Wait for all worker goroutines to complete
		wg.Wait()
		close(filePaths)
	}()

	// Walk the folder and send file paths to the channel
	walkWG.Add(1)
	go func() {
		defer walkWG.Done()
		err := filepath.Walk(f.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				// Increment WaitGroup before starting goroutine
				wg.Add(1)
				go func() {
					// Defer the decrement of WaitGroup after processing
					defer wg.Done()

					file := fileutil.File{Path: path}
					if err := file.CalculateMD5Checksum(); err != nil {
						f.Logger.Color.Printf("Error calculating checksum for file %s: %v\n", path, err)
					}

					// Lock the mutex before modifying the shared slice
					mu.Lock()
					f.Files = append(f.Files, file)
					mu.Unlock()
				}()
			}
			return nil
		})
		if err != nil {
			// Handle the error, if needed
		}
	}()

	// Wait for the walking goroutine to complete
	walkWG.Wait()

	return nil
}

// ScanWithoutConcurrency scans the folder without concurrency based on the provided config.
func (f *Folder) ScanWithoutConcurrency(config *config.Config) error {
	// Implementation for scanning
	f.Logger.Color.Printf("Scanning folder: %s\n", f.Path)

	// Clear existing files
	f.Files = nil

	// Walk the folder and process files
	err := filepath.Walk(f.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file := fileutil.File{Path: path}
			if err := file.CalculateMD5Checksum(); err != nil {
				f.Logger.Color.Printf("Error calculating checksum for file %s: %v\n", path, err)
			}

			// Modify the shared slice without using concurrency
			f.Files = append(f.Files, file)
		}
		return nil
	})

	if err != nil {
		// Handle the error, if needed
		return err
	}

	return nil
}

// BackupToFolder backs up the folder to the specified backup folder.
func (f *Folder) BackupToFolder(backupFolder string) {
	// Implementation for backup
	fmt.Printf("Backing up folder %s to %s\n", f.Path, backupFolder)
	// Your backup logic here
}
