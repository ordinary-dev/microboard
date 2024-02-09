package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/h2non/filetype"

	"github.com/ordinary-dev/microboard/config"
)

// Saves the uploaded file.
// Returns:
// 1. Relative path to it in the uploaded files directory.
// 2. Detected mime type.
func SaveBuffer(cfg *config.Config, buf []byte, boardCode string) (string, string, error) {
	// Get mime type
	fileKind, err := filetype.Match(buf)
	if err != nil || fileKind == filetype.Unknown {
		fileKind.MIME.Value = "application/octet-stream"
		fileKind.Extension = "bin"
	}

	// Calculate hash
	hasher := sha256.New()
	hasher.Write(buf)
	hashsum := hasher.Sum(nil)

	// Get file type
	filetype := "other"
	if strings.HasPrefix(fileKind.MIME.Value, "image") {
		filetype = "image"
	} else if strings.HasPrefix(fileKind.MIME.Value, "video") {
		filetype = "video"
	} else if strings.HasPrefix(fileKind.MIME.Value, "audio") {
		filetype = "audio"
	}

	// Generate file path
	filepath := path.Join(
		boardCode,
		filetype,
		fmt.Sprintf("%x.%v", hashsum, fileKind.Extension),
	)
	fullFilePath := path.Join(cfg.UploadDir, filepath)

	if _, err := os.Stat(fullFilePath); errors.Is(err, os.ErrNotExist) {
		if err = os.WriteFile(fullFilePath, buf, 0644); err != nil {
			return "", "", err
		}
	}

	return filepath, fileKind.MIME.Value, nil
}
