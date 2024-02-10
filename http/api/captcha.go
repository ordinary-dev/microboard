package api

import (
	"bytes"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ordinary-dev/microboard/captcha"
	dbcaptchas "github.com/ordinary-dev/microboard/database/captchas"
)

func ShowCaptcha(db *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		captchaIDText := ctx.Param("id")
		captchaID, err := uuid.Parse(captchaIDText)
		if err != nil {
			ctx.Error(err)
			return
		}

		answer, err := dbcaptchas.GetAnswer(ctx, db, captchaID)
		if err != nil {
			ctx.Error(err)
			return
		}

		img, err := captcha.GenerateCaptcha(answer)
		if err != nil {
			ctx.Error(err)
			return
		}

		var buf bytes.Buffer
		png.Encode(&buf, img)

		ctx.DataFromReader(http.StatusOK, int64(buf.Len()), "image/png", &buf, map[string]string{
			"content-disposition": "inline; filename=\"captcha.png\"",
		})
	}
}
