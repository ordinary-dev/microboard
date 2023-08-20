package api

import (
	"github.com/gin-gonic/gin"
)

func ToPaginatedResult[T any](count int64, results []T) gin.H {
	return gin.H{
		"count":   count,
		"results": results,
	}
}
