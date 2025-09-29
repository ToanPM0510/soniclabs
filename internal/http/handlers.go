package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"toanpm0510/soniclabs/internal/domain"
	"toanpm0510/soniclabs/internal/store/pg"
)

type Server struct {
	Log *zap.Logger
	DB  *pg.DB
}

func (s *Server) getCourses(c *gin.Context) {
	courses, err := pg.ListCourses(c, s.DB)
	if err != nil {
		s.Log.Error("list courses", zap.Error(err))
		WriteProblem(c, http.StatusInternalServerError, "Internal Server Error", "cannot list courses", "INT_001")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": courses})
}

type postCourseReq struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Difficulty  string  `json:"difficulty"`
}

func (s *Server) postCourse(c *gin.Context) {
	var req postCourseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteBadRequest(c, "invalid JSON body", "VAL_000")
		return
	}
	diff := domain.Difficulty(req.Difficulty)
	if err := domain.ValidateNewCourse(req.Title, diff); err != nil {
		WriteBadRequest(c, err.Error(), "VAL_001")
		return
	}
	course, err := pg.CreateCourse(c, s.DB, req.Title, req.Description, diff)
	if err != nil {
		s.Log.Error("create course", zap.Error(err))
		WriteProblem(c, http.StatusInternalServerError, "Internal Server Error", "cannot create course", "INT_002")
		return
	}
	c.JSON(http.StatusCreated, course)
}

type postEnrollReq struct {
	StudentEmail string `json:"student_email"`
	CourseID     string `json:"course_id"`
}

func (s *Server) postEnroll(c *gin.Context) {
	var req postEnrollReq
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteBadRequest(c, "invalid JSON body", "VAL_000")
		return
	}
	if err := domain.ValidateEmail(req.StudentEmail); err != nil {
		WriteBadRequest(c, err.Error(), "VAL_010")
		return
	}
	cid, err := uuid.Parse(req.CourseID)
	if err != nil {
		WriteBadRequest(c, "invalid course_id", "VAL_011")
		return
	}
	enr, err := pg.Enroll(c, s.DB, req.StudentEmail, cid)
	if err != nil {
		if err == pg.ErrDuplicateEnrollment {
			WriteProblem(c, http.StatusConflict, "Conflict", "student already enrolled in this course", "BUS_001")
			return
		}
		s.Log.Error("enroll", zap.Error(err))
		WriteProblem(c, http.StatusInternalServerError, "Internal Server Error", "cannot enroll", "INT_003")
		return
	}
	c.JSON(http.StatusCreated, enr)
}

func (s *Server) getStudentEnrollments(c *gin.Context) {
	email := c.Param("email")
	if err := domain.ValidateEmail(email); err != nil {
		WriteBadRequest(c, err.Error(), "VAL_010")
		return
	}
	list, err := pg.ListEnrollmentsByEmail(c, s.DB, email)
	if err != nil {
		s.Log.Error("list enrollments", zap.Error(err))
		WriteProblem(c, http.StatusInternalServerError, "Internal Server Error", "cannot list enrollments", "INT_004")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}
