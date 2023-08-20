package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/microboard/database"
	"net/http"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Create the first administrator.
// This function will return an error if at least one user already exists.
func CreateUser(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		count, err := db.GetAdminCount()
		if err != nil {
			ctx.Error(err)
			return
		}

		if count > 0 {
			ctx.Error(errors.New("at least 1 user already exists"))
			return
		}

		var creds User
		err = ctx.ShouldBindJSON(&creds)
		if err != nil {
			ctx.Error(err)
			return
		}

		admin, err := db.CreateAdmin(creds.Username, creds.Password)
		if err != nil {
			ctx.Error(err)
			return
		}

		token, err := db.GetAccessToken(admin.ID)
		if err != nil {
			// This shouldn't happen.
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}

// Get an access token for the user.
// At the moment, the token is valid for 24 hours.
// It should be saved in the "microboard-token" cookie.
//
// POST /api/v0/users/token
func GetAccessToken(db *database.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user User
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			ctx.Error(err)
			return
		}

		admin, err := db.VerifyPassword(user.Username, user.Password)
		if err != nil {
			// This is most likely a failed login attempt.
			ctx.Error(err)
			return
		}

		token, err := db.GetAccessToken(admin.ID)
		if err != nil {
			// This shouldn't happen.
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}
