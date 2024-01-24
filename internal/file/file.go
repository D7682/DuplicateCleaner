// internal/file/file.go
package file

import (
	"DuplicateCleaner/logger"
	"crypto/md5"
	"encoding/hex"
	"github.com/fatih/color"
	"os"
)

// File represents a file with additional information.
type File struct {
	Path        string
	MD5Checksum string
	// Add more fields as needed
}

// NewFile creates a new File instance.
func NewFile(path string) (*File, error) {
	md5Checksum, err := generateMD5Checksum(path)
	if err != nil {
		return nil, err
	}

	return &File{
		Path:        path,
		MD5Checksum: md5Checksum,
	}, nil
}

// generateMD5Checksum generates the MD5 checksum for a file.
func generateMD5Checksum(path string) (string, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	hash.Write(fileContent)
	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

// LogInfo logs file information with pastel colors.
func (f *File) LogInfo(log *logger.Logger) {
	log.PastelPrintf(color.FgHiCyan, "File: %s\n", f.Path)
	log.PastelPrintf(color.FgHiYellow, "MD5 Checksum: %s\n", f.MD5Checksum)
	log.PastelPrintf(color.FgHiWhite, "------------------------\n")
}

// LogInfoToFile logs file information to the log file.
func (f *File) LogInfoToFile(log *logger.Logger) {
	log.PastelPrintf(color.FgHiCyan, "File: %s\n", f.Path)
	log.PastelPrintf(color.FgHiYellow, "MD5 Checksum: %s\n", f.MD5Checksum)
	log.PastelPrintf(color.FgHiWhite, "------------------------\n")
	log.Printf("\n") // Add a newline for better separation in the log file
}
