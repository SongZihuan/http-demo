package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func HandlerDatetime(c *gin.Context) {
	str := time.Now().Format("2006-01-15 15:04:05")
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(str)))
	_, _ = c.Writer.WriteString(str)
	c.Status(http.StatusOK)
}
