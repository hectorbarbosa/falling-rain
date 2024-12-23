package stream

import (
	"log/slog"
	"sync"

	"github.com/gdamore/tcell"
	"golang.org/x/exp/rand"
)

// vertical bar containing streams
type StreamBar struct {
	screen    tcell.Screen
	logger    *slog.Logger
	column    int
	newStream chan bool
	newHead   chan StreamShift
	streams   map[int]Stream
	height    int
}

func NewBar(screen tcell.Screen, logger *slog.Logger, height int, i int) StreamBar {
	newStream := make(chan bool, 1)
	streams := make(map[int]Stream)
	newHead := make(chan StreamShift, 1)

	return StreamBar{
		screen:    screen,
		logger:    logger,
		column:    i,
		newStream: newStream,
		newHead:   newHead,
		streams:   streams,
		height:    height,
	}
}

func (bar *StreamBar) Run(done chan struct{}, wg *sync.WaitGroup) {
	bar.logger.Debug("StreamBar Run")
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	blackStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorBlack)
	bar.screen.SetStyle(defStyle)

	midStyleA := blackStyle.Foreground(tcell.ColorGreen)
	midStyleB := blackStyle.Foreground(tcell.ColorLime)
	headStyleA := blackStyle.Foreground(tcell.ColorSilver)
	headStyleB := blackStyle.Foreground(tcell.ColorWhite)

	var s StreamShift

	firstStream := NewStream(bar.logger)
	go firstStream.Run(bar.height, done, bar.newHead, bar.newStream)

STREAMS:
	for {
		select {
		case <-done:
			break STREAMS
		case <-bar.newStream:
			stream := NewStream(bar.logger)
			bar.logger.Debug("New stream")
			go stream.Run(bar.height, done, bar.newHead, bar.newStream)
		case s = <-bar.newHead:
			// Making about a third of the green characters bright-lime...
			if rand.Intn(100) > 33 {
				bar.screen.SetContent(bar.column, s.yHead-1, s.secondRune, nil, midStyleA)
			} else {
				bar.screen.SetContent(bar.column, s.yHead-1, s.secondRune, nil, midStyleB)
			}

			// turning a fifth of the heads from gray to white
			if rand.Intn(100) > 20 {
				bar.screen.SetContent(bar.column, s.yHead, s.headRune, nil, headStyleA)
			} else {
				bar.screen.SetContent(bar.column, s.yHead, s.headRune, nil, headStyleB)
			}

			bar.screen.SetContent(bar.column, s.yTail, ' ', nil, blackStyle)

		}
	}

	wg.Done()
}
