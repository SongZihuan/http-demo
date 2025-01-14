package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandlerEmpty(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", "0")
	c.Status(http.StatusNoContent)
}
