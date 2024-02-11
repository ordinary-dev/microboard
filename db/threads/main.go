package threads

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/ordinary-dev/microboard/db"
	"github.com/ordinary-dev/microboard/db/files"
	"github.com/ordinary-dev/microboard/db/posts"
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
type ThreadWithFirstAndLastPosts struct {
	ID            uint64                `json:"id"`
	BoardCode     string                `json:"boardCode" db:"board_code"`
	UpdatedAt     time.Time             `json:"updatedAt" db:"updated_at"`
	DeletedAt     *time.Time            `json:"deletedAt" db:"deleted_at"`
	PostID        uint64                `json:"postID" db:"post_id"`
	Body          string                `json:"body"`
	PostCreatedAt time.Time             `json:"postCreatedAt" db:"post_created_at"`
	Files         []files.File          `json:"file"`
	LatestPosts   []posts.PostWithFiles `json:"lastPosts"`
}

// Start a new thread.
// Thread ID and post ID will be filled.
// The "CreatedAt" and "UpdatedAt" fields will be overwritten.
func CreateThread(thread *Thread, post *posts.Post, files []files.File) error {
	tx, err := db.DB.Begin(context.Background())
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
	return err
}

// Get threads that haven't been deleted.
func GetThreads(boardCode string, limit, offset int) ([]ThreadWithFirstAndLastPosts, error) {
	query := `SELECT * FROM (
        SELECT DISTINCT ON (threads.id)
            threads.id,
            threads.board_code,
            threads.updated_at,
            posts.body,
            posts.id AS post_id,
            posts.created_at AS post_created_at
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
	rows, _ := db.DB.Query(context.Background(), query, args)
	threads, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[ThreadWithFirstAndLastPosts])
	if err != nil {
		return nil, err
	}

	for idx := range threads {
		threads[idx].Files, err = files.GetFilesByPostID(threads[idx].PostID)
		if err != nil {
			return nil, err
		}

		threads[idx].LatestPosts, err = posts.GetLatestPostsFromThread(threads[idx].ID)
		if err != nil {
			return nil, err
		}
	}

	return threads, nil
}

// Get thread by id.
func GetThread(threadID uint64) (Thread, error) {
	query := `SELECT * FROM threads WHERE id = $1`
	rows, _ := db.DB.Query(context.Background(), query, threadID)
	return pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[Thread])
}

// Get thread count.
func GetThreadCount(boardCode string) (int64, error) {
	query := `SELECT COUNT(*) FROM threads WHERE board_code = $1 AND deleted_at IS NULL`
	var count int64
	err := db.DB.QueryRow(context.Background(), query, boardCode).Scan(&count)
	return count, err
}

// Hides the thread from public access.
func DeleteThread(threadID uint64) error {
	query := `UPDATE threads SET deleted_at = @deletedAt WHERE id = @threadID`
	args := pgx.NamedArgs{
		"deletedAt": time.Now(),
		"threadID":  threadID,
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
