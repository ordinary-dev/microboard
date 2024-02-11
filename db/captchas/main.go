package captchas

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ordinary-dev/microboard/db"
)

var (
	ErrWrongAnswer = errors.New("captcha is invalid")
	alphabet       = []rune("abcdefghjkmnpqrstuvwxyz123456789")
)

type Captcha struct {
	ID        uuid.UUID `db:"id" json:"id" form:"id"`
	Answer    string    `db:"answer" json:"answer" form:"answer"`
	CreatedAt time.Time `db:"created_at" json:"createdAt" form:"createdAt"`
}

func CreateCaptcha(ctx context.Context) (*Captcha, error) {
	answer := GenerateRandomString(6)

	query := `INSERT INTO captchas(answer) VALUES ($1) RETURNING id`
	captcha := Captcha{Answer: answer, CreatedAt: time.Now()}
	err := db.DB.QueryRow(ctx, query, answer).Scan(&captcha.ID)

	return &captcha, err
}

func GetAnswer(ctx context.Context, captchaID uuid.UUID) (answer string, err error) {
	query := `SELECT answer FROM captchas WHERE id = $1`
	err = db.DB.QueryRow(ctx, query, captchaID).Scan(&answer)
	return
}

// Returns `true` if the answer is valid.
// Captcha will be deleted and cannot be reused.
func ValidateCaptcha(ctx context.Context, captchaID uuid.UUID, answer string) (bool, error) {
	answer = strings.ToLower(answer)

	query := `
        DELETE FROM captchas
        WHERE id = $1 AND answer = $2
    `
	cmdTag, err := db.DB.Exec(ctx, query, captchaID, answer)

	return cmdTag.RowsAffected() == 1, err
}

func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}
