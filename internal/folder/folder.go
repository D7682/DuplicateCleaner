// internal/folder/folder.go

package folder

import (
	"DuplicateCleaner/internal/config"
	"DuplicateCleaner/internal/file"
	"DuplicateCleaner/logger"
	"fmt"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// Folder represents a folder and its related operations.
type Folder struct {
	Path  string       // Path to the folder
	Files []*file.File // Slice to hold File instances

	config.Config                // Embedding config.Config
	Logger        *logger.Logger // Embedding logger.Logger
	mu            sync.Mutex
}

// NewFolder creates a new Folder instance with the given path, config, and logger.
func NewFolder(path string, cfg *config.Config) *Folder {
	return &Folder{
		Path:   path,
		Config: *cfg,
		Logger: logger.NewLogger(cfg.LogFilePath),
	}
}

func (f *Folder) getTotalFiles() (int64, error) {
	var totalFiles int64

	err := filepath.Walk(f.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalFiles++
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return totalFiles, nil
}

func (f *Folder) ScanWithConfig() error {
	f.Logger.CyanPrintf("Scanning folder: %s\n", f.Path)

	var (
		doneCh     = make(chan struct{})
		start      = time.Now()
		wg         sync.WaitGroup
		mu         sync.Mutex
		counter    int64 // Atomic counter for progress bar
		totalFiles int64 // Total file count
	)

	// Clear existing files
	f.Files = nil

	// Get the total file count
	totalFiles, err := f.getTotalFiles()
	if err != nil {
		return err
	}

	// Create a progress bar with the initial file count (zero in this case)
	bar := progressbar.NewOptions64(totalFiles)

	// File processing loop
	processFile := func(filePath string) {
		defer wg.Done()
		newFile, err := file.NewFile(filePath)
		if err != nil {
			f.Logger.PrettyError(err)
			return
		}

		// Append the file to f.Files
		atomic.AddInt64(&counter, 1)
		mu.Lock()
		f.Files = append(f.Files, newFile)
		bar.Set64(atomic.LoadInt64(&counter))
		progressString := bar.String()
		mu.Unlock()

		// Print progress bar string
		fmt.Printf("\r%s", progressString)
	}

	// Concurrent file processing
	err = filepath.Walk(f.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		wg.Add(1)
		go processFile(path)
		return nil
	})
	if err != nil {
		return err
	}

	// Wait for all worker goroutines to finish
	wg.Wait()

	// Close the channel to signal completion
	close(doneCh)

	// Display results after all workers are done
	f.DisplayResults("\nScan with Concurrency")
	elapsedConcurrent := time.Since(start)
	f.Logger.PastelBluePrintf("[Duplicate Cleaner] %s Time Elapsed (Concurrent): %.2fs\n", time.Now().Format("2006/01/02 15:04:05"), elapsedConcurrent.Seconds())

	return nil
}

// countFilesInFolder counts the total number of files in a folder.
func countFilesInFolder(folderPath string) (int, error) {
	var count int
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

func (f *Folder) ScanWithoutConcurrency() error {
	f.Logger.CyanPrintf("Scanning folder: %s\n", f.Path)

	var (
		start      = time.Now()
		totalFiles int64 // Total file count
	)

	// Clear existing files
	f.Files = nil

	// Get the total file count
	totalFiles, err := f.getTotalFiles()
	if err != nil {
		return err
	}

	// Create a progress bar with the initial file count (zero in this case)
	bar := progressbar.NewOptions64(totalFiles)

	// File processing loop
	err = filepath.Walk(f.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		newFile, err := file.NewFile(path)
		if err != nil {
			f.Logger.PrettyError(err)
			return nil // Continue with other files even if one fails
		}

		// Append the file to f.Files
		f.Files = append(f.Files, newFile)

		// Update progress bar
		bar.Add(1)

		// Print progress bar string
		fmt.Printf("\r%s", bar.String())

		return nil
	})
	if err != nil {
		return err
	}

	// Display results after all files are processed
	f.DisplayResults("\nScan without Concurrency")
	elapsed := time.Since(start)
	f.Logger.PastelBluePrintf("[Duplicate Cleaner] %s Time Elapsed: %.2fs\n", time.Now().Format("2006/01/02 15:04:05"), elapsed.Seconds())

	return nil
}

// BackupToFolder backs up the folder to the specified backup folder.
func (f *Folder) BackupToFolder() {
	f.Logger.CyanPrintf("Backing up folder %s to %s\n", f.Path, f.BkpFolder)
	// Your backup logic here
}

// DisplayResults displays the results of scanning using the logger with pastel colors.
func (f *Folder) DisplayResults(title string) {
	f.Logger.PastelPrintf(color.FgHiWhite, "%s Results:\n", title)
	f.Logger.PastelPrintf(color.FgHiCyan, "Source Folder: %v\n", f.SrcFolder)
	f.Logger.PastelPrintf(color.FgHiYellow, "Backup Folder: %v\n", f.BkpFolder)
	f.Logger.PastelPrintf(color.FgHiMagenta, "Max Scan Depth: %v\n", f.MaxScanDepth)

	// Display file count using pastel green color
	fileCount := len(f.Files)
	f.Logger.PastelPrintf(color.FgHiGreen, "File Count: %d\n", fileCount)

	// Display additional information with pastel colors
	cyan := color.New(color.FgHiCyan).SprintFunc()
	magenta := color.New(color.FgHiMagenta).SprintFunc()

	// Calculate and display the number of duplicate files
	duplicateCount := f.calculateDuplicateCount()
	f.Logger.PastelPrintf(color.FgHiWhite, "[%s] %s %s: %s\n",
		cyan("Duplicate Cleaner"), magenta(time.Now().Format("2006/01/02 15:04:05")),
		cyan(title), magenta(fmt.Sprintf("Found %d duplicate files", duplicateCount)))
}

// calculateDuplicateCount calculates the total number of duplicate files.
func (f *Folder) calculateDuplicateCount() int {
	// Use a map to store unique MD5 checksums
	uniqueChecksums := make(map[string]bool)

	// Iterate through files and count duplicates
	duplicateCount := 0
	for _, fil := range f.Files {
		if _, found := uniqueChecksums[fil.MD5Checksum]; found {
			// Found a duplicate
			duplicateCount++
		} else {
			// Add the checksum to the map
			uniqueChecksums[fil.MD5Checksum] = true
		}
	}

	return duplicateCount
}
