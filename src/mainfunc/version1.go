package mainfunc

import (
	"errors"
	"fmt"
	resource "github.com/SongZihuan/http-demo"
	"github.com/SongZihuan/http-demo/src/engine"
	"github.com/SongZihuan/http-demo/src/flagparser"
	"github.com/SongZihuan/http-demo/src/httpserver"
	"github.com/SongZihuan/http-demo/src/httpsslserver"
	"github.com/SongZihuan/http-demo/src/signalchan"
)

func MainV1() (exitcode int) {
	defer func() {
		if recover() != nil {
			exitcode = 1
			return
		}
	}()

	fmt.Printf("")
	err := flagparser.InitFlagParser()
	if err != nil {
		fmt.Printf("init flag fail: %s\n", err.Error())
		return 1
	}

	if flagparser.Version {
		fmt.Printf("Version: %s\n", resource.Version)
		return 0
	}

	if flagparser.DryRun {
		flagparser.Print()
		return 0
	}

	err = engine.InitEngine()
	if err != nil {
		fmt.Printf("init engine fail: %s\n", err.Error())
		return 1
	}

	err = signalchan.InitSignal()
	if err != nil {
		fmt.Printf("init signal fail: %s\n", err.Error())
		return 1
	}
	defer signalchan.CloseSignal()

	var httpchan = make(chan error)
	var httpsslchan = make(chan error)
	defer func() {
		close(httpchan)
		close(httpsslchan)
	}()

	if flagparser.HttpAddress != "" {
		err = httpserver.InitHttpServer()
		if err != nil {
			fmt.Printf("init http server fail: %s\n", err.Error())
			return 1
		}

		go func() {
			httpchan <- httpserver.RunServer()
		}()
	}

	if flagparser.HttpsAddress != "" {
		err = httpsslserver.InitHttpSSLServer()
		if err != nil {
			fmt.Printf("init https server fail: %s\n", err.Error())
			return 1
		}

		go func() {
			httpsslchan <- httpsslserver.RunServer()
		}()
	}

	select {
	case <-signalchan.SignalChan:
		fmt.Printf("Server closed: safe\n")
		return 0
	case err := <-httpchan:
		if errors.Is(err, httpserver.ErrStop) {
			fmt.Printf("Http Server closed: safe\n")
			return 0
		}
		fmt.Printf("Http Server error closed: %s\n", err.Error())
		return 1
	case err := <-httpsslchan:
		if errors.Is(err, httpsslserver.ErrStop) {
			fmt.Printf("Https Server closed: safe\n")
			return 0
		}
		fmt.Printf("Https Server error closed: %s\n", err.Error())
		return 1
	}
	// 后续不可达
}
