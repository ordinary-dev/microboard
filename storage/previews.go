package storage

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/sirupsen/logrus"

	"github.com/ordinary-dev/microboard/config"
	"github.com/ordinary-dev/microboard/db/files"
	"github.com/ordinary-dev/microboard/db/posts"
)

// Image size in pixels.
const PREVIEW_SIZE = 100

func GenerateMissingPreviews(cfg *config.Config) {
	posts, err := posts.GetPostsWithMissingPreviews()
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, post := range posts {
		go GeneratePreviewsForPost(cfg, post.ID)
	}
}

// Requests all files from the post without previews and tries to create them.
func GeneratePreviewsForPost(cfg *config.Config, postID uint64) {
	logrus.Debugf("Generating preview for post #%v", postID)

	files, err := files.GetFilesWithoutPreview(postID)
	if err != nil {
		logrus.Error(err)
		return
	}

	waitGroup := sync.WaitGroup{}

	for idx := range files {
		waitGroup.Add(1)
		go processFile(cfg, &files[idx], &waitGroup)
	}

	waitGroup.Wait()
}

// Generate a preview for a single file.
func processFile(cfg *config.Config, file *files.File, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	var err error

	originalFilepath := path.Join(cfg.UploadDir, file.Filepath)

	// Generate file path
	filepathWithoutExt := strings.TrimSuffix(file.Filepath, filepath.Ext(file.Filepath))
	filepathWithExt := filepathWithoutExt + ".webp"
	previewPath := path.Join(cfg.PreviewDir, filepathWithExt)

	// Generate preview
	if _, err := os.Stat(previewPath); errors.Is(err, os.ErrNotExist) {
		// The preview has not been generated yet.
		if strings.HasPrefix(file.MimeType, "image") {
			err = processImageFile(originalFilepath, previewPath)
		} else if strings.HasPrefix(file.MimeType, "video") {
			err = processVideoFile(originalFilepath, previewPath)
		} else {
			logrus.Debugf("Unsupported mimetype: %v", file.MimeType)
		}

		if err != nil {
			logrus.Error(err)
			return
		}
	}

	file.Preview = &filepathWithExt
	err = files.UpdateFilePreview(file)
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Debugf("Preview successfully generated")
}

func processImageFile(originalFilepath, previewPath string) error {
	image, err := vips.NewImageFromFile(originalFilepath)
	if err != nil {
		return err
	}

	// Rotate the picture upright and reset EXIF orientation tag.
	if err = image.AutoRotate(); err != nil {
		return err
	}

	// Resize the picture
	if err = image.Thumbnail(PREVIEW_SIZE, PREVIEW_SIZE, vips.InterestingAttention); err != nil {
		return err
	}

	// Set export parameters
	ep := vips.NewWebpExportParams()
	ep.Quality = 70
	ep.StripMetadata = true
	imageBytes, _, err := image.ExportWebp(ep)
	if err != nil {
		return err
	}

	// Save preview
	preview, err := os.Create(previewPath)
	if err != nil {
		return err
	}
	defer preview.Close()
	if _, err = preview.Write(imageBytes); err != nil {
		return err
	}

	return nil
}

func processVideoFile(originalFilepath, previewPath string) error {
	// Run ffmpeg
	cmd := exec.Command(
		"ffmpeg",
		"-i", originalFilepath,
		"-frames:v", "1",
		"-y",
		"-vf", fmt.Sprintf("thumbnail,scale=%v:%v:force_original_aspect_ratio=increase", PREVIEW_SIZE, PREVIEW_SIZE),
		previewPath,
	)

	out, err := cmd.Output()
	if err != nil {
		logrus.Debug(string(out))
		return fmt.Errorf("%v: %v", cmd.String(), err)
	}

	return nil
}
