package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func HandlerTimestamp(c *gin.Context) {
	str := fmt.Sprintf("%d", time.Now().Unix())
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(str)))
	_, _ = c.Writer.WriteString(str)
	c.Status(http.StatusOK)
}
