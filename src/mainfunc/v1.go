package mainfunc

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/Http-Demo/src/engine"
	"github.com/SongZihuan/Http-Demo/src/flagparser"
	"github.com/SongZihuan/Http-Demo/src/httpserver"
	"github.com/SongZihuan/Http-Demo/src/httpsslserver"
	"github.com/SongZihuan/Http-Demo/src/signalchan"
)

func MainV1() (exitcode int) {
	err := flagparser.InitFlag()
	if err != nil {
		fmt.Printf("init flag fail: %s\n", err.Error())
		return 1
	}

	if flagparser.DryRun {
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

	err = httpserver.InitHttpServer()
	if err != nil {
		fmt.Printf("init http server fail: %s\n", err.Error())
		return 1
	}

	var httpchan = make(chan error)
	var httpsslchan = make(chan error)

	go func() {
		httpchan <- httpserver.RunServer()
	}()

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
		if errors.Is(err, httpserver.ErrStop) {
			fmt.Printf("Https Server closed: safe\n")
			return 0
		}
		fmt.Printf("Https Server error closed: %s\n", err.Error())
		return 1
	}
}
