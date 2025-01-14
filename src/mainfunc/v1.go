package mainfunc

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/Http-Demo/src/engine"
	"github.com/SongZihuan/Http-Demo/src/flagparser"
	"github.com/SongZihuan/Http-Demo/src/httpserver"
	"github.com/SongZihuan/Http-Demo/src/signalchan"
)

func MainV1() (exitcode int) {
	err := flagparser.InitFlag()
	if err != nil {
		fmt.Printf("init flag fail: %s\n", err.Error())
		return 1
	}

	err = engine.InitEngine()
	if err != nil {
		fmt.Printf("init engine fail: %s\n", err.Error())
		return 1
	}

	err = httpserver.InitHttpServer()
	if err != nil {
		return 1
	}

	err = signalchan.InitSignal()
	if err != nil {
		fmt.Printf("init http server fail: %s\n", err.Error())
		return 1
	}

	var httpchan = make(chan error)
	go func() {
		httpchan <- httpserver.RunServer()
	}()

	select {
	case <-signalchan.SignalChan:
		fmt.Printf("Server closed: safe\n")
		return 0
	case err := <-httpchan:
		if errors.Is(err, httpserver.ErrStop) {
			fmt.Printf("Server closed: safe\n")
			return 0
		}
		fmt.Printf("Server error closed: %s\n", err.Error())
		return 1
	}
}
