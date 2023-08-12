package database

import (
	"context"
	"github.com/jackc/pgx/v5"
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
}

func (db *DB) CreateFile(file *File) error {
	query := `INSERT INTO files (post_id, filepath, name, size, mimetype)
        VALUES (@postID, @filepath, @name, @size, @mimetype)
        RETURNING id`
	args := pgx.NamedArgs{
		"postID":   file.PostID,
		"filepath": file.Filepath,
		"name":     file.Name,
		"size":     file.Size,
		"mimetype": file.MimeType,
	}

	err := db.pool.QueryRow(context.Background(), query, args).Scan(&file.ID)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetFilesByPostID(postID uint64) ([]File, error) {
	query := `SELECT * FROM files WHERE post_id = @postID`
	args := pgx.NamedArgs{
		"postID": postID,
	}

	rows, err := db.pool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files, err := pgx.CollectRows(rows, pgx.RowToStructByName[File])
	if err != nil {
		return nil, err
	}

	return files, nil
}