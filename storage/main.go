package storage

import (
	"os"
	"path"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/db/boards"
)

// Create the necessary directories for storing user files.
func CreateDirs(cfg *config.Config) error {
	// uploads
	if err := os.MkdirAll(cfg.UploadDir, os.ModePerm); err != nil {
		return err
	}

	// uploads/{board.Code}/{image,video,audio,other}/
	filetypes := []string{"image", "video", "audio", "other"}
	boards, err := boards.GetBoards()
	if err != nil {
		return err
	}

	for _, filetype := range filetypes {
		for _, board := range boards {
			if err := os.MkdirAll(path.Join(cfg.UploadDir, board.Code, filetype), os.ModePerm); err != nil {
				return err
			}
		}
	}

	// previews
	if err := os.MkdirAll(cfg.PreviewDir, os.ModePerm); err != nil {
		return err
	}

	// previews/{board.Code}/{image,video,audio,other}/
	for _, filetype := range filetypes {
		for _, board := range boards {
			if err := os.MkdirAll(path.Join(cfg.PreviewDir, board.Code, filetype), os.ModePerm); err != nil {
				return err
			}
		}
	}

	return nil
}
