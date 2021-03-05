package scanner

import (
	"bytes"
	"io"
)

type Scanner struct {
	r   *reader
	buf bytes.Buffer
}


type BufScanner struct {
	s   *Scanner
	i   int // buffer index
	n   int // buffer size
	buf [3]TokenInfo
}

type reader struct {
	r   io.RuneScanner
	i   int // buffer index
	n   int // buffer char count
	pos Pos // last read rune position
	buf [3]struct {
		ch  rune
		pos Pos
	}
	eof bool // true if reader has ever seen eof.
}
