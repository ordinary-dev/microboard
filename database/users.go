package database

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/argon2"

	"github.com/ordinary-dev/microboard/config"
)

const (
	ARGON2_TIME    uint32 = 1
	ARGON2_MEMORY  uint32 = 64 * 1024
	ARGON2_THREADS uint8  = 4
	ARGON2_KEYLEN  uint32 = 32
)

var (
	ErrWrongPassword           = errors.New("wrong password")
	ErrTokenWasNotFound        = errors.New("token was not found")
	ErrTokenHasExpired         = errors.New("token has expired")
	ErrTokenIsEmpty            = errors.New("token is empty")
	ErrInvalidTokenOwner       = errors.New("invalid token owner")
	ErrEmptyUsernameOrPassword = errors.New("empty username or password")
	ErrAtLeastOneUserExists    = errors.New("at least one user exists")
)

// A user with additional privileges.
type Admin struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
}

type adminWithHash struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Salt     []byte `json:"-"`
	Hash     []byte `json:"-"`
}

// Access token for admins.
type AccessToken struct {
	Value     string
	AdminID   int32     `db:"admin_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (db *DB) CreateDefaultUser(cfg *config.Config) error {
	if cfg.DefaultUsername == "" || cfg.DefaultPassword == "" {
		return ErrEmptyUsernameOrPassword
	}

	count, err := db.GetAdminCount()
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrAtLeastOneUserExists
	}

	if _, err := db.CreateAdmin(cfg.DefaultUsername, cfg.DefaultPassword); err != nil {
		return err
	}

	return nil
}

// Save new admin user in the database.
// Before calling this function, you need to make sure that the creator has the necessary rights.
func (db *DB) CreateAdmin(username string, password string) (*Admin, error) {
	salt, err := getRandomBytes(16)
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey([]byte(password), salt, ARGON2_TIME, ARGON2_MEMORY, ARGON2_THREADS, ARGON2_KEYLEN)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO admins (username, salt, hash) VALUES (@username, @salt, @hash) RETURNING id`
	args := pgx.NamedArgs{
		"username": username,
		"salt":     salt,
		"hash":     hash,
	}

	admin := Admin{
		Username: username,
	}

	err = db.Pool.QueryRow(context.Background(), query, args).Scan(&admin.ID)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (db *DB) VerifyPassword(username string, password string) (*Admin, error) {
	admin, err := db.getAdminByUsername(username)
	if err != nil {
		return nil, err
	}

	otherHash := argon2.IDKey([]byte(password), admin.Salt, ARGON2_TIME, ARGON2_MEMORY, ARGON2_THREADS, ARGON2_KEYLEN)

	if subtle.ConstantTimeCompare(admin.Hash, otherHash) != 1 {
		return nil, ErrWrongPassword
	}

	adminWithoutHash := Admin{
		ID:       admin.ID,
		Username: admin.Username,
	}

	return &adminWithoutHash, nil
}

func (db *DB) getAdminByUsername(username string) (*adminWithHash, error) {
	query := `SELECT id, salt, hash FROM admins WHERE username = @username`
	args := pgx.NamedArgs{
		"username": username,
	}

	admin := adminWithHash{
		Username: username,
	}

	err := db.Pool.QueryRow(context.Background(), query, args).Scan(&admin.ID, &admin.Salt, &admin.Hash)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// Get an access token for the user.
// At the moment, the token is valid for 7 days.
// It should be saved in the "microboard-token" cookie.
func (db *DB) GetAccessToken(adminID int32) (*AccessToken, error) {
	query := `SELECT value, created_at FROM access_tokens WHERE admin_id = @adminID`
	args := pgx.NamedArgs{
		"adminID": adminID,
	}

	rows, err := db.Pool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokenExists := false

	if rows.Next() {
		accessToken := AccessToken{
			AdminID: adminID,
		}
		err = rows.Scan(&accessToken.Value, &accessToken.CreatedAt)
		if err != nil {
			return nil, err
		}

		if err := accessToken.IsValid(); err == nil {
			return &accessToken, nil
		} else {
			tokenExists = true
		}
	}

	// Token doesn't exist or it's not valid.
	// Create a new one.

	rawToken, err := getRandomBytes(32)
	if err != nil {
		return nil, err
	}
	token := base64.StdEncoding.EncodeToString(rawToken)

	if tokenExists {
		query = `UPDATE access_tokens SET value = @value, created_at = @createdAt WHERE admin_id = @adminID`
	} else {
		query = `INSERT INTO access_tokens(value, admin_id, created_at) VALUES (@value, @adminID, @createdAt)`
	}

	now := time.Now()
	args = pgx.NamedArgs{
		"value":     token,
		"adminID":   adminID,
		"createdAt": now,
	}

	_, err = db.Pool.Exec(context.Background(), query, args)
	if err != nil {
		return nil, err
	}

	accessToken := AccessToken{
		Value:     token,
		AdminID:   adminID,
		CreatedAt: now,
	}

	return &accessToken, nil
}

func (db *DB) ValidateAccessToken(tokenValue string) (*AccessToken, error) {
	query := `SELECT admin_id, created_at FROM access_tokens WHERE value = @value`
	args := pgx.NamedArgs{
		"value": tokenValue,
	}

	rows, err := db.Pool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		accessToken := AccessToken{
			Value: tokenValue,
		}

		err = rows.Scan(&accessToken.AdminID, &accessToken.CreatedAt)
		if err != nil {
			return nil, err
		}

		if err := accessToken.IsValid(); err != nil {
			return nil, err
		}

		return &accessToken, nil
	}

	return nil, ErrTokenWasNotFound
}

func (db *DB) GetAdminCount() (int64, error) {
	query := `SELECT COUNT(*) FROM admins`
	var count int64
	err := db.Pool.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func getRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (token *AccessToken) IsValid() error {
	if time.Since(token.CreatedAt).Hours() >= 24*7 {
		return ErrTokenHasExpired
	}

	if token.Value == "" {
		return ErrTokenIsEmpty
	}

	if token.AdminID < 1 {
		return ErrInvalidTokenOwner
	}

	return nil
}
