package scanner

import (
	"bufio"
	"io"
)

func init() {

}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r: &reader{
			r: bufio.NewReaderSize(r,128),
		},
	}
}