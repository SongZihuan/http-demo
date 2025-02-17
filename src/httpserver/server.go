package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/SongZihuan/http-demo/src/engine"
	"github.com/SongZihuan/http-demo/src/flagparser"
	"github.com/pires/go-proxyproto"
	"net"
	"net/http"
	"time"
)

var HttpServer *http.Server = nil
var HttpListener net.Listener = nil
var HttpAddress string

var ErrStop = fmt.Errorf("http server error")

func InitHttpServer() error {
	HttpAddress = flagparser.HttpAddress

	HttpServer = &http.Server{
		Addr:    HttpAddress,
		Handler: engine.Engine,
	}

	return nil
}

func RunServer() (err error) {
	tcpListener, err := net.Listen("tcp", HttpServer.Addr)
	if err != nil {
		return err
	}

	proxyListener := &proxyproto.Listener{
		Listener:          tcpListener,
		ReadHeaderTimeout: 10 * time.Second,
	}

	HttpListener = proxyListener

	defer func() {
		_ = HttpListener.Close()

		HttpServer = nil
		HttpListener = nil
	}()

	fmt.Printf("http server start at %s\n", HttpAddress)
	err = HttpServer.Serve(HttpListener)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return ErrStop
	}
	return fmt.Errorf("http server error: %s", err)
}

func StopServer() (err error) {
	if HttpServer == nil {
		return nil
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	err = HttpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	HttpServer = nil
	HttpListener = nil

	return nil
}
