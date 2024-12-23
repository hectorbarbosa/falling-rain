package main

import (
	"flag"
	"fmt"
	"os"
	"rain/logging"
	"rain/stream"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	var logLevel int

	// parse flags, default flag value is Info
	flag.IntVar(&logLevel, "loglevel", 0, "Logging level: -4 Debug, default 0")
	flag.Parse()

	// Create logger
	logger, err := logging.GetLogger(logLevel)
	if err != nil {
		fmt.Printf("getting logger: %s\n", err)
		os.Exit(1)
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		logger.Error("%+v", "NewScreen", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		logger.Error("%+v", "Init", err)
		os.Exit(1)
	}

	// check screen size
	width, height := screen.Size()

	logger.Info("---------------------------")
	logger.Info("Screen size: %d, %d", "width", width, "height", height)

	// clear screen
	screen.HideCursor()
	screen.Clear()
	time.Sleep(time.Millisecond * 500)

	// flusher flushes the termbox every x milliseconds
	curFPS := 60
	fpsSleepTime := time.Duration(1000000/curFPS) * time.Microsecond
	logger.Info("fps sleep time", "st", fpsSleepTime.String())
	go func() {
		for {
			time.Sleep(fpsSleepTime)
			screen.Show()
		}
	}()

	// make chan for termbox events and run poller to send events on chan
	eventChan := make(chan tcell.Event)
	go func() {
		for {
			event := screen.PollEvent()
			eventChan <- event
		}
	}()

	done := make(chan struct{})

	var wg sync.WaitGroup

	// make n=width vartical bars
	for i := 0; i < width; i++ {
		wg.Add(1)
		s := stream.NewBar(screen, logger, height, i)
		go s.Run(done, &wg)
	}
	// s := stream.NewBar(screen, logger, height, 1)
	// go s.Run(done)

EVENTS:
	for {
		e := <-eventChan
		switch event := e.(type) {
		case *tcell.EventKey:
			switch event.Key() {

			case tcell.KeyCtrlC:
				logger.Info("Ctrl+c pressed, exiting")
				close(done)
				break EVENTS

			case tcell.KeyRune:
				switch event.Rune() {
				case 'q':
					logger.Info("q pressed, exiting")
					close(done)
					break EVENTS
				case 'c':
					logger.Info("c pressed, clearing screen")
					screen.Clear()
				}
			}
		}
	}

	wg.Wait()
	screen.Fini()
}
