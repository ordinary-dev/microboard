package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

type Thread struct {
	// The thread ID is only used in the url.
	// It is not displayed on the site.
	ID uint64 `json:"id"`

	BoardCode string `json:"boardCode" db:"board_code"`

	// Time when the last post in the thread was created.
	// The field is used to speed up sorting.
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`

	// This field is populated if the thread has been banned by an admin
	// but has not yet been deleted due to board limits.
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

// The same structure as a regular thread, but with additional fields from the post.
type ThreadWithFirstPost struct {
	ID            uint64     `json:"id"`
	BoardCode     string     `json:"boardCode" db:"board_code"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt     *time.Time `json:"deletedAt" db:"deleted_at"`
	PostID        uint64     `json:"postID" db:"post_id"`
	Body          string     `json:"body"`
	PostCreatedAt time.Time  `json:"postCreatedAt" db:"post_created_at"`
	Files         []File     `json:"file"`
}

// Start a new thread.
// Thread ID and post ID will be filled.
// The "CreatedAt" and "UpdatedAt" fields will be overwritten.
func (db *DB) CreateThread(thread *Thread, post *Post, files []File) error {
	tx, err := db.pool.Begin(context.Background())
	if err != nil {
		return err
	}

	// Rollback is safe to call.
	defer tx.Rollback(context.Background())

	now := time.Now()
	thread.UpdatedAt = now
	post.CreatedAt = now

	// Save thread
	query := `INSERT INTO threads (board_code, updated_at) VALUES (@boardCode, @updatedAt) RETURNING id`
	args := pgx.NamedArgs{
		"boardCode": thread.BoardCode,
		"updatedAt": now,
	}
	err = tx.QueryRow(context.Background(), query, args).Scan(&thread.ID)
	if err != nil {
		return err
	}
	post.ThreadID = thread.ID

	// Save post
	query = `INSERT INTO posts (thread_id, body, created_at) VALUES (@threadID, @body, @createdAt) RETURNING id`
	args = pgx.NamedArgs{
		"threadID":  thread.ID,
		"body":      post.Body,
		"createdAt": now,
	}
	err = tx.QueryRow(context.Background(), query, args).Scan(&post.ID)
	if err != nil {
		return err
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

// Get threads that haven't been deleted.
func (db *DB) GetThreads(boardCode string, limit, offset int) ([]ThreadWithFirstPost, error) {
	query := `SELECT * FROM (
        SELECT DISTINCT ON (threads.id) threads.id, threads.board_code, threads.updated_at, posts.body, posts.id AS post_id, posts.created_at AS post_created_at
        FROM threads

        INNER JOIN posts
        ON posts.thread_id = threads.id

        WHERE threads.board_code = @boardCode
	    AND threads.deleted_at IS NULL

        ORDER BY threads.id, posts.created_at
    ) threads_with_first_post
    ORDER BY threads_with_first_post.updated_at DESC
    LIMIT @limit
    OFFSET @offset`

	args := pgx.NamedArgs{
		"boardCode": boardCode,
		"limit":     limit,
		"offset":    offset,
	}
	rows, err := db.pool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[ThreadWithFirstPost])
	if err != nil {
		return nil, err
	}

	for idx := range threads {
		files, err := db.GetFilesByPostID(threads[idx].PostID)
		if err != nil {
			return nil, err
		}

		threads[idx].Files = files
	}

	return threads, nil
}

// Get thread by id.
func (db *DB) GetThread(threadID uint64) (*Thread, error) {
	query := `SELECT * FROM threads WHERE id = @threadID`
	args := pgx.NamedArgs{
		"threadID": threadID,
	}
	rows, err := db.pool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	thread, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[Thread])
	if err != nil {
		return nil, err
	}

	return &thread, nil
}

// Get thread count.
func (db *DB) GetThreadCount(boardCode string) (int64, error) {
	query := `SELECT COUNT(*) FROM threads WHERE board_code = @boardCode AND deleted_at IS NULL`
	args := pgx.NamedArgs{
		"boardCode": boardCode,
	}
	var count int64
	err := db.pool.QueryRow(context.Background(), query, args).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Hides the thread from public access.
func (db *DB) DeleteThread(threadID uint64) error {
	query := `UPDATE threads SET deleted_at = @deletedAt WHERE id = @threadID`
	args := pgx.NamedArgs{
		"deletedAt": time.Now(),
		"threadID":  threadID,
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
