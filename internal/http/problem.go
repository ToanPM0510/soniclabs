package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Code     string `json:"code,omitempty"`
}

func WriteProblem(c *gin.Context, status int, title, detail, code string) {
	c.AbortWithStatusJSON(status, Problem{
		Type:   "about:blank",
		Title:  title,
		Detail: detail,
		Code:   code,
	})
}

func WriteBadRequest(c *gin.Context, detail string, code string) {
	WriteProblem(c, http.StatusBadRequest, "Bad Request", detail, code)
}
