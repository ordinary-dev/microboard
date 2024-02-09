package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Board struct {
	// Board code, no longer than 2-4 characters.
	// For example: 'r'.
	Code string `json:"code" form:"code"`

	// The name of the board, a couple of words long.
	// For example: 'random'.
	Name string `json:"name" form:"boardName"`

	// Information about the board, a couple of sentences long.
	Description string `json:"description" form:"description"`

	// The number of pages after which threads will be deleted.
	// 0 - no limit.
	// Typical values: 0, 10.
	PageLimit int16 `json:"pageLimit" db:"page_limit" form:"pageLimit"`

	// The number of posts after which the thread will not rise higher in the list.
	// 0 - no limit.
	// Typical values: 500, 1000.
	BumpLimit int16 `json:"bumpLimit" db:"bump_limit" form:"bumpLimit"`

	// Is the board hidden from the public list?
	Unlisted bool `json:"unlisted" form:"unlisted"`
}

// Save board to the database.
// Fill in the ID.
func (db *DB) CreateBoard(board *Board) error {
	query := `INSERT INTO boards(code, name, description, page_limit, bump_limit, unlisted)
        VALUES (@code, @name, @description, @pageLimit, @bumpLimit, @unlisted)`
	args := pgx.NamedArgs{
		"code":        board.Code,
		"name":        board.Name,
		"description": board.Description,
		"pageLimit":   board.PageLimit,
		"bumpLimit":   board.BumpLimit,
		"unlisted":    board.Unlisted,
	}

	_, err := db.pool.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}

	return nil
}

// Get all boards, even unlisted.
func (db *DB) GetBoards() ([]Board, error) {
	query := `SELECT * FROM boards ORDER BY code`

	rows, err := db.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards, err := pgx.CollectRows(rows, pgx.RowToStructByName[Board])
	if err != nil {
		return nil, err
	}

	return boards, nil
}

// Get board info.
func (db *DB) GetBoard(boardCode string) (*Board, error) {
	query := `SELECT * FROM boards WHERE code = @code`
	args := pgx.NamedArgs{
		"code": boardCode,
	}

	rows, err := db.pool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	board, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Board])
	if err != nil {
		return nil, err
	}

	return &board, nil
}

// Save updated board information to the database.
func (db *DB) UpdateBoard(board *Board) error {
	query := `UPDATE boards
        SET name = @name, description = @description, page_limit = @pageLimit, bump_limit = @bumpLimit, unlisted = @unlisted
        WHERE code = @code`
	args := pgx.NamedArgs{
		"code":        board.Code,
		"name":        board.Name,
		"description": board.Description,
		"pageLimit":   board.PageLimit,
		"bumpLimit":   board.BumpLimit,
		"unlisted":    board.Unlisted,
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
