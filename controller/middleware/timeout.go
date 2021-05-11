package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

//
// timeout.go
// Copyright (C) 2021 forseason <me@forseason.vip>
//
// Distributed under terms of the MIT license.
//

func Timeout(timeout time.Duration) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		newctx, cancel := context.WithTimeout(ctx.Request.Context(), timeout)
		defer func() {
			if newctx.Err() == context.DeadlineExceeded {
				ctx.Writer.WriteHeader(504)
				ctx.Abort()
			}
			cancel()
		}()
		ctx.Request = ctx.Request.WithContext(newctx)
		ctx.Next()
	}
}
