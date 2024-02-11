package files

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/ordinary-dev/microboard/db"
)

type File struct {
	ID uint64 `json:"id"`

	PostID uint64 `json:"postID" db:"post_id"`

	Filepath string `json:"filepath"`

	// The original name of the uploaded file with extension.
	Name string `json:"name"`

	// File size in bytes.
	Size uint64 `json:"size"`

	// Guessed file type (e.g. image/png)
	MimeType string `json:"mimetype"`

	// The name of the preview file, which is located in cfg.UploadDir/previews/<preview>
	Preview *string `json:"preview"`
}

func CreateFile(file *File) error {
	query := `
        INSERT INTO files (post_id, filepath, name, size, mimetype, preview)
        VALUES (@postID, @filepath, @name, @size, @mimetype, @preview)
        RETURNING id`
	args := pgx.NamedArgs{
		"postID":   file.PostID,
		"filepath": file.Filepath,
		"name":     file.Name,
		"size":     file.Size,
		"mimetype": file.MimeType,
		"preview":  file.Preview,
	}

	return db.DB.QueryRow(context.Background(), query, args).Scan(&file.ID)
}

func GetFilesByPostID(postID uint64) ([]File, error) {
	query := `SELECT * FROM files WHERE post_id = $1`

	rows, _ := db.DB.Query(context.Background(), query, postID)
	return pgx.CollectRows(rows, pgx.RowToStructByName[File])
}

func GetFilesWithoutPreview(postID uint64) ([]File, error) {
	query := `SELECT * FROM files WHERE preview IS NULL AND post_id = $1`

	rows, _ := db.DB.Query(context.Background(), query, postID)
	return pgx.CollectRows(rows, pgx.RowToStructByName[File])
}

func UpdateFilePreview(file *File) error {
	query := `UPDATE files SET preview = @preview WHERE id = @id`
	args := pgx.NamedArgs{
		"preview": file.Preview,
		"id":      file.ID,
	}

	cmdTag, err := db.DB.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() != 1 {
		return pgx.ErrNoRows
	}

	return nil
}
