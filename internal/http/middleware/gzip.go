package middleware

import "github.com/gin-gonic/gin"

func Gzip() gin.HandlerFunc { // placeholder: dùng gin’s built-in gzip in prod
	return func(c *gin.Context) { c.Next() }
}
