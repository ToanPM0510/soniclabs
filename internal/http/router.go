package httpx

import (
	"net/http"

	"toanpm0510/soniclabs/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(l *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Recover(l), middleware.CORS(), middleware.Timeout(10_000_000_000))
	// TODO: swap to real gzip middleware
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.GET("/readyz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	// v1 routes placeholder
	api := r.Group("/api/v1")
	{
		api.GET("/courses", func(c *gin.Context) { c.JSON(200, gin.H{"data": []string{}}) })
	}
	return r
}
