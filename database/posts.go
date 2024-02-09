package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Post struct {
	ID uint64 `json:"id"`

	ThreadID uint64 `json:"threadID" db:"thread_id"`

	// Text with markup.
	Body string `json:"text"`

	// Post creation time.
	// When a post is created, the value is copied into the 'Thread.UpdatedAt' field
	// if the number of posts has not exceeded the limit.
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type PostWithFiles struct {
	ID        uint64     `json:"id"`
	ThreadID  uint64     `json:"threadID" db:"thread_id"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
	Files     []File     `json:"files"`
}

func (db *DB) CreatePost(post *Post, files []File) error {
	tx, err := db.pool.Begin(context.Background())
	if err != nil {
		return err
	}

	// Rollback is safe to call.
	defer tx.Rollback(context.Background())

	query := `INSERT INTO posts (thread_id, body, created_at) VALUES (@threadID, @body, @createdAt) RETURNING id`
	args := pgx.NamedArgs{
		"threadID":  post.ThreadID,
		"body":      post.Body,
		"createdAt": post.CreatedAt,
	}

	err = tx.QueryRow(context.Background(), query, args).Scan(&post.ID)
	if err != nil {
		return err
	}

	// Get information about the thread and board.
	query = `SELECT COUNT(*), boards.bump_limit
        FROM posts
        INNER JOIN threads
        ON posts.thread_id = threads.id
        INNER JOIN boards
        ON threads.board_code = boards.code
        WHERE posts.thread_id = @threadID
        GROUP BY threads.id, boards.bump_limit`
	var postCount int64
	var bumpLimit int64
	err = tx.QueryRow(context.Background(), query, args).Scan(&postCount, &bumpLimit)
	if err != nil {
		return err
	}

	// Move the thread higher in the list.
	if postCount < bumpLimit || bumpLimit == 0 {
		query = `UPDATE threads SET updated_at = @createdAt WHERE id = @threadID`
		_, err = tx.Exec(context.Background(), query, args)
		if err != nil {
			return err
		}
	}

	// Save files
	query = `INSERT INTO files (post_id, filepath, name, size, mimetype) VALUES (@postID, @filepath, @name, @size, @mimetype) RETURNING id`
	for idx, file := range files {
		args := pgx.NamedArgs{
			"postID":   post.ID,
			"filepath": file.Filepath,
			"name":     file.Name,
			"size":     file.Size,
			"mimetype": file.MimeType,
		}
		err := tx.QueryRow(context.Background(), query, args).Scan(&files[idx].ID)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetPostsFromThread(threadID uint64) ([]PostWithFiles, error) {
	query := `SELECT * FROM posts WHERE thread_id = @threadID`
	args := pgx.NamedArgs{
		"threadID": threadID,
	}

	rows, err := db.pool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[PostWithFiles])
	if err != nil {
		return nil, err
	}

	for idx := range posts {
		files, err := db.GetFilesByPostID(posts[idx].ID)
		if err != nil {
			return nil, err
		}
		posts[idx].Files = files
	}

	return posts, nil
}

func (db *DB) GetPostsWithMissingPreviews() ([]Post, error) {
	query := `SELECT DISTINCT ON (posts.id) posts.* FROM posts
        INNER JOIN files
        ON files.post_id = posts.id
        WHERE files.preview IS NULL`

	rows, err := db.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Post])
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// Marks a post as deleted (but does not remove it from the database)
func (db *DB) DeletePost(postID uint64) error {
	query := `UPDATE posts SET deleted_at = @deletedAt WHERE id = @postID`
	args := pgx.NamedArgs{
		"postID":    postID,
		"deletedAt": time.Now(),
	}

	cmdTag, err := db.pool.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() != 1 {
		return ErrNoRowsWereAffected
	}

	return nil
}

// Get 3 latest posts from the thread, not including the first one
func (db *DB) GetLatestPostsFromThread(threadID uint64) ([]PostWithFiles, error) {
	query := `SELECT sub.*
        FROM (
            SELECT * FROM posts
            WHERE thread_id = @threadID
            AND deleted_at IS NULL
            AND id NOT IN (
                SELECT id FROM posts
                WHERE thread_id = @threadID
                ORDER BY created_at ASC
                LIMIT 1
            )
            ORDER BY created_at DESC
            LIMIT 3
        ) sub
        ORDER BY created_at ASC`
	args := pgx.NamedArgs{
		"threadID": threadID,
	}

	rows, err := db.pool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[PostWithFiles])
	if err != nil {
		return nil, err
	}

	for idx := range posts {
		files, err := db.GetFilesByPostID(posts[idx].ID)
		if err != nil {
			return nil, err
		}

		posts[idx].Files = files
	}

	return posts, nil
}
