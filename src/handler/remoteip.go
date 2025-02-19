package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandlerRemoteIP(c *gin.Context) {
	str := c.RemoteIP()
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(str)))
	_, _ = c.Writer.WriteString(str)
	c.Status(http.StatusOK)
}
