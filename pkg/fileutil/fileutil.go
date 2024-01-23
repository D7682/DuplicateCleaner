package fileutil

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// File represents a file with its path and MD5 checksum.
type File struct {
	Path        string
	MD5Checksum string
}

// CalculateMD5Checksum calculates the MD5 checksum for the file.
func (f *File) CalculateMD5Checksum() error {
	content, err := os.ReadFile(f.Path)
	if err != nil {
		return err
	}

	md5Checksum := md5.Sum(content)
	f.MD5Checksum = hex.EncodeToString(md5Checksum[:])

	return nil
}
