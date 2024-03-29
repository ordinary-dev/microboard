package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ordinary-dev/microboard/db/users"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Create the first administrator.
// This function will return an error if at least one user already exists.
func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		count, err := users.GetAdminCount()
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

		admin, err := users.CreateAdmin(creds.Username, creds.Password)
		if err != nil {
			ctx.Error(err)
			return
		}

		token, err := users.GetAccessToken(admin.ID)
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
//
// POST /api/v0/users/token
func GetAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user User
		err := ctx.ShouldBindJSON(&user)
		if err != nil {
			ctx.Error(err)
			return
		}

		admin, err := users.VerifyPassword(user.Username, user.Password)
		if err != nil {
			// This is most likely a failed login attempt.
			ctx.Error(err)
			return
		}

		token, err := users.GetAccessToken(admin.ID)
		if err != nil {
			// This shouldn't happen.
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}
