package httpserver

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/http-demo/src/engine"
	"github.com/SongZihuan/http-demo/src/flagparser"
	"net/http"
)

var HttpServer *http.Server = nil
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

func RunServer() error {
	fmt.Printf("http server start at %s\n", HttpAddress)
	err := HttpServer.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return ErrStop
	}
	return fmt.Errorf("http server error: %s", err)
}
