package httpx

import (
	"net/http"
	"time"

	"toanpm0510/soniclabs/internal/http/middleware"
	"toanpm0510/soniclabs/internal/store/pg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(l *zap.Logger, db *pg.DB) *gin.Engine {
	r := gin.New()
	r.Use(
		middleware.RequestID(),
		middleware.Recover(l),
		middleware.CORS(),
		middleware.Timeout(10*time.Second),
	)

	// health endpoints
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.GET("/readyz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	s := &Server{Log: l, DB: db}

	api := r.Group("/api/v1")
	{
		api.GET("/courses", s.getCourses)
		api.POST("/courses", s.postCourse)

		api.POST("/enrollments", s.postEnroll)
		api.GET("/students/:email/enrollments", s.getStudentEnrollments)
	}

	return r
}
