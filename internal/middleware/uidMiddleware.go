package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/srcgod/apigateway/internal/utils"
)

const (
	UidInt64Key = "uid_int64"
)

func UIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, exists := ctx.Get("uid")
		if !exists {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "uid not found in context"},
			)
			return
		}
		uidInt64 := utils.ConvertToInt64(uid)
		if uidInt64 == 0 {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "invalid uid format"},
			)
			return
		}
		ctx.Set(UidInt64Key, uidInt64)
		ctx.Next()
	}
}
