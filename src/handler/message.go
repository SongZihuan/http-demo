package handler

import (
	"fmt"
	resource "github.com/SongZihuan/Http-Demo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func HandlerMessage(c *gin.Context) {
	nowTime := time.Now()

	var res strings.Builder

	res.WriteString(fmt.Sprintf("Hello, this is HTTP-DEMO %s\n", resource.Version))
	res.WriteString(fmt.Sprintf("Date: %s\n", nowTime.Format("2006-01-02 15:04:05")))
	res.WriteString(fmt.Sprintf("Timestamp(Unix Second): %d\n", nowTime.Unix()))
	res.WriteString(fmt.Sprintf("Host: %s\n", c.Request.Host))
	res.WriteString(fmt.Sprintf("Proto: %s\n", c.Request.Proto))
	if c.Request.TLS == nil {
		res.WriteString(fmt.Sprintf("Https/TLS: %s\n", "No"))
		res.WriteString(fmt.Sprintf("Http/TLS: %s\n", "Yes"))
		res.WriteString(fmt.Sprintf("Scheme: %s\n", "HTTP"))
	} else {
		res.WriteString(fmt.Sprintf("Https/TLS: %s\n", "Yes"))
		res.WriteString(fmt.Sprintf("Http/TLS: %s\n", "No"))
		res.WriteString(fmt.Sprintf("Scheme: %s\n", "HTTPS"))
	}
	res.WriteString(fmt.Sprintf("Path: %s\n", c.Request.URL.Path))
	res.WriteString(fmt.Sprintf("Query: %s\n", c.Request.URL.RawQuery))
	res.WriteString(fmt.Sprintf("ClientIP: %s\n", c.ClientIP()))
	res.WriteString(fmt.Sprintf("RemoteIP: %s\n", c.RemoteIP()))
	res.WriteString(fmt.Sprintf("Via: %s\n", c.Request.Header.Get(RequestsXVia)))
	res.WriteString(fmt.Sprintf("Forwarded: %s\n", c.Request.Header.Get(RequestsForwarded)))
	res.WriteString(fmt.Sprintf("X-Forwarded-For: %s\n", c.Request.Header.Get(RequestsXForwardedFor)))
	res.WriteString(fmt.Sprintf("X-Forwarded-Proto: %s\n", c.Request.Header.Get(RequestsXForwardedProto)))
	res.WriteString(fmt.Sprintf("X-Forwarded-Host: %s\n", c.Request.Header.Get(RequestsXForwardedHost)))
	res.WriteString(fmt.Sprintf("X-Forwarded-Host: %s\n", c.Request.Header.Get(RequestsXForwardedHost)))
	res.WriteString(fmt.Sprintf("X-Message: %s\n", strings.Join(c.Request.Header.Values(RequestsXMessage), " ")))

	str := res.String()
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(str)))
	_, _ = c.Writer.WriteString(str)
	c.Status(http.StatusOK)
}
