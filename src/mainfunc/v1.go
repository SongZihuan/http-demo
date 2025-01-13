package mainfunc

import (
	"errors"
	"flag"
	"fmt"
	resource "github.com/SongZihuan/Http-Demo"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	RequestsForwarded       = "Forwarded"
	RequestsXForwardedFor   = "X-Forwarded-For"
	RequestsXForwardedHost  = "X-Forwarded-Host"
	RequestsXForwardedProto = "X-Forwarded-Proto"
	RequestsXMessage        = "X-Message"
	RequestsXVia            = "Via"
)

var address string
var engine *gin.Engine = nil
var server *http.Server = nil

func MainV1() (exitcode int) {
	func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println("option parse error")
				exitcode = 1
				return
			}
		}()

		flag.StringVar(&address, "address", ":3366", "http server address")
		flag.StringVar(&address, "a", ":3366", "http server address")
		flag.Parse()
	}()

	gin.SetMode(gin.ReleaseMode)

	engine = gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/", Handler)
	engine.GET("/ip", HandlerRemoteIP)
	engine.GET("/client/ip", HandlerClientIP)
	engine.GET("/timestamp", HandlerTimestamp)
	engine.GET("/datetime", HandlerDatetime)
	engine.GET("/hello", HandlerHelloWorld)
	engine.GET("/empty", HandlerEmpty)
	engine.NoRoute(HandlerMethodNotFound)
	engine.NoMethod(HandlerMethodNotAllowed)

	server = &http.Server{
		Addr:    address,
		Handler: engine,
	}

	var sigchan = make(chan os.Signal)
	var stopchan = make(chan error)

	err := initSignal(sigchan)
	if err != nil {
		fmt.Printf("Listen signal fail: %s\n", err.Error())
		return 1
	}

	go func() {
		fmt.Printf("server start at %s\n", address)
		err := server.ListenAndServe()
		if err != nil {
			stopchan <- err
		}
	}()

	select {
	case <-sigchan:
		fmt.Printf("Server closed: safe\n")
		return 0
	case err := <-stopchan:
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server closed: safe\n")
			return 0
		}
		fmt.Printf("Server error closed: %s\n", err.Error())
		return 1
	}
}

func initSignal(sigchan chan os.Signal) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("init signal error: %v", r)
		}
	}()

	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	return nil
}

func Handler(c *gin.Context) {
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

func HandlerRemoteIP(c *gin.Context) {
	str := c.RemoteIP()
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(str)))
	_, _ = c.Writer.WriteString(str)
	c.Status(http.StatusOK)
}

func HandlerClientIP(c *gin.Context) {
	str := c.ClientIP()
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(str)))
	_, _ = c.Writer.WriteString(str)
	c.Status(http.StatusOK)
}

func HandlerTimestamp(c *gin.Context) {
	str := fmt.Sprintf("%d", time.Now().Unix())
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(str)))
	_, _ = c.Writer.WriteString(str)
	c.Status(http.StatusOK)
}

func HandlerDatetime(c *gin.Context) {
	str := time.Now().Format("2006-01-15 15:04:05")
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(str)))
	_, _ = c.Writer.WriteString(str)
	c.Status(http.StatusOK)
}

func HandlerHelloWorld(c *gin.Context) {
	str := "Hello, world!"
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(str)))
	_, _ = c.Writer.WriteString(str)
	c.Status(http.StatusOK)
}

func HandlerEmpty(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Content-Length", "0")
	c.Status(http.StatusNoContent)
}

func HandlerMethodNotFound(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotFound)
}

func HandlerMethodNotAllowed(c *gin.Context) {
	c.AbortWithStatus(http.StatusMethodNotAllowed)
}
