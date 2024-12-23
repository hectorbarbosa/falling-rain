package stream

import (
	"log/slog"
	"math/rand"
	"time"
)

// stream struct contains info about one stream
type Stream struct {
	length int
	speed  int
	head   int
	tail   int
	quitCh chan bool
	logger *slog.Logger
}

type StreamShift struct {
	yHead      int
	yTail      int
	headRune   rune
	secondRune rune
}

var halfWidthKana = []rune{
	'｡', '｢', '｣', '､', '･', 'ｦ', 'ｧ', 'ｨ', 'ｩ', 'ｪ', 'ｫ', 'ｬ', 'ｭ', 'ｮ', 'ｯ',
	'ｰ', 'ｱ', 'ｲ', 'ｳ', 'ｴ', 'ｵ', 'ｶ', 'ｷ', 'ｸ', 'ｹ', 'ｺ', 'ｻ', 'ｼ', 'ｽ', 'ｾ', 'ｿ',
	'ﾀ', 'ﾁ', 'ﾂ', 'ﾃ', 'ﾄ', 'ﾅ', 'ﾆ', 'ﾇ', 'ﾈ', 'ﾉ', 'ﾊ', 'ﾋ', 'ﾌ', 'ﾍ', 'ﾎ', 'ﾏ',
	'ﾐ', 'ﾑ', 'ﾒ', 'ﾓ', 'ﾔ', 'ﾕ', 'ﾖ', 'ﾗ', 'ﾘ', 'ﾙ', 'ﾚ', 'ﾛ', 'ﾜ', 'ﾝ', 'ﾞ', 'ﾟ',
}

var alphaNumerics = []rune{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
}

var allTheCharacters = append(halfWidthKana, alphaNumerics...)

func NewStream(logger *slog.Logger) Stream {
	quitCh := make(chan bool)
	length := 10 + rand.Intn(8)
	tail := -1 * length
	speed := 50 + rand.Intn(90)

	return Stream{
		length: length,
		speed:  speed,
		head:   0,
		tail:   tail,
		quitCh: quitCh,
		logger: logger,
	}
}

func (s *Stream) Run(
	barHeight int,
	quit chan struct{},
	newHead chan StreamShift,
	newStream chan bool,
) {

	headRune := allTheCharacters[rand.Intn(len(allTheCharacters))]
	var secondRune rune

	s.logger.Debug("first stream run")

STREAM:
	for {
		select {
		case <-quit:
			s.logger.Debug("stream stopped")
			break STREAM
		case <-time.After(time.Duration(s.speed) * time.Millisecond):
			// the stream is out of screen, exiting ...
			if s.tail >= barHeight {
				break STREAM
			}

			streamShift := StreamShift{s.head, s.tail, headRune, secondRune}

			newHead <- streamShift

			if s.head < barHeight+s.length {
				s.head++
				secondRune = headRune
				headRune = allTheCharacters[rand.Intn(len(allTheCharacters))]
				s.tail++
				if s.tail == barHeight/2 {
					newStream <- true
				}
			}
		}
	}

}
