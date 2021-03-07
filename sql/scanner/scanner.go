package scanner

import (
	"bytes"
	"errors"
	"io"
	"unicode/utf8"
)

type Scanner struct {
	r   *reader
	buf bytes.Buffer
}

/*read */
/*rune 代表四个字节，byte 代表一个字节*/
/*一种是 uint8 类型，或者叫 byte 型，代表了 ASCII 码的一个字符。*/
/*一种是 rune 类型，代表一个 UTF-8 字符，当需要处理中文、日文或者其他复合字符时，则需要用到 rune 类型。rune 类型等价于 int32 类型*/
func (s *Scanner) read() (ch rune, pos Pos) {
	ch, pos = s.r.read()
	if ch != eof {
		/*讲ch字符写入缓存中*/
		s.buf.WriteRune(ch)
	}
	return
}


func (s *Scanner) unread() {
	if ch, _ := s.r.curr(); ch != eof {
		s.buf.Truncate(s.buf.Len() - utf8.RuneLen(ch))
	}
	s.r.unread()
}

func (s *Scanner) unbuffer() string {
	str := s.buf.String()
	s.buf.Reset()

	return str
}

/*Scan returns the next token and position from the underlying reader.*/
/*TODO*/
func (s *Scanner) Scan() TokenInfo {
	ch0, pos := s.read()

	if isWhitespace(ch0){
		return	s.scanWhitespace()
	}else if isLetter(ch0) || ch0 == '_'{
		s.unread()
		return s.scanIdent(true)
	}else if isDigit(ch0) {
		return s.scanIdent()
	}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() TokenInfo {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	ch, pos := s.r.curr()
	_, _ = buf.WriteRune(ch)

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		ch, _ = s.read()
		if ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return TokenInfo{WS, pos, buf.String(), s.unbuffer()}
}

func (s *Scanner) scanIdent(lookup bool) TokenInfo {
	// Save the starting position of the identifier.
	_, pos := s.read()
	s.unread()

	var buf bytes.Buffer
	for {
		if ch, _ := s.read(); ch == eof {
			break
		} else if ch == '`' {
			ti0 := s.scanString()
			if ti0.Tok == BADSTRING || ti0.Tok == BADESCAPE {
				return ti0
			}
			return TokenInfo{IDENT, pos, ti0.Lit, ti0.Raw}
		} else if isIdentChar(ch) {
			s.unread()
			bi := ScanBareIdent(s.r.r)
			buf.WriteString(bi)
			s.buf.WriteString(bi)
		} else {
			s.unread()
			break
		}
	}
	lit := buf.String()

	// If the literal matches a keyword then return that keyword.
	if lookup {
		if tok := Lookup(lit); tok != IDENT {
			return TokenInfo{tok, pos, "", s.unbuffer()}
		}
	}
	return TokenInfo{IDENT, pos, lit, s.unbuffer()}
}

// scanString consumes a contiguous string of non-quote characters.
// Quote characters can be consumed if they're first escaped with a backslash.
func (s *Scanner) scanString() TokenInfo {
	s.unread()
	_, pos := s.r.curr()

	lit, err := ScanString(s)

	if err == errBadString {
		return TokenInfo{BADSTRING, pos, lit, s.unbuffer()}
	} else if err == errBadEscape {
		_, pos = s.r.curr()
		return TokenInfo{BADESCAPE, pos, lit, s.unbuffer()}
	}
	return TokenInfo{STRING, pos, lit, s.unbuffer()}
}

/*BufScanner represents a wrapper for scanner to add a buffer.*/
/*It provides a fixed-length circular buffer that can be unread.*/
/*循环缓冲区，无法读取*/
type BufScanner struct {
	s   *Scanner
	i   int // buffer index
	n   int // buffer size
	buf [3]TokenInfo
}

/*NewBufScanner returns a new buffered scanner for a reader*/
func NewBufScanner(r io.Reader) *BufScanner {
	return &BufScanner{
		s: NewScanner(r),
	}
}

/*Scan reads the next token from the scanner.*/
func (s *BufScanner) Scan() TokenInfo {
	return s.scanFunc(s)
}

func (s *BufScanner) scanFunc(scan func() TokenInfo) TokenInfo {

}

/*reader represents a buffered rune reader used by the scanner*/
/*It provides a fixed-length circular buffer that can be unread*/
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

/*TokenInfo holds information about a token*/
type TokenInfo struct {
	Tok Token
	Pos Pos
	Lit string
	Raw string
}

/*eof is a marker code point to signify that the reader can't read any more.*/
var eof = rune(0)

/*isWhitespace returns true if the rune is a space, tab, or newline.*/
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}
// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// isIdentChar returns true if the rune can be used in an unquoted identifier.
func isIdentChar(ch rune) bool { return isLetter(ch) || isDigit(ch) || ch == '_' }

// isIdentFirstChar returns true if the rune can be used as the first char in an unquoted identifer.
func isIdentFirstChar(ch rune) bool { return isLetter(ch) || ch == '_' }

/*read reads the next rune from the reader*/
func (r reader) read() (ch rune, pos Pos) {
	/*If we have unread characters then read them off the buffer first.*/
	/*如果我们有未读的字符，我们先从里面读取它*/
	if r.n > 0 {
		r.n--
		return r.curr()
	}
	/*Read next rune from underlying reader.*/
	/*阅读下一个字符*/
	/*Any error (including io.EOF) should return as EOF.*/
	/*ReadRune读取单个utf-8编码的字符，返回该字符和它的字节长度。如果没有有效的字符，会返回错误。*/
	ch, _, err := r.r.ReadRune()
	if err != nil {
		ch = eof
	} else if ch == '\r' { // 换行
		if ch, _, err := r.r.ReadRune(); err != nil {
			// nop
		} else if ch != '\n' { // 回车
			/*UnreadRune方法让下一次调用ReadRune时返回之前调用ReadRune时返回的同一个utf-8字符。
			连续调用两次UnreadRune方法而中间没有调用ReadRune时，可能会导致错误。
			*/
			_ = r.r.UnreadRune()
		}
		ch = '\n'
	}

	/*save character and position to the buffer */
	r.i = (r.i + 1) % len(r.buf)
	buf := &r.buf[r.i]
	buf.ch, buf.pos = ch, r.pos

	// Update position.
	// Only count EOF once.
	if ch == '\n' {
		r.pos.Line++
		r.pos.Char = 0
	} else if !r.eof {
		r.pos.Char++
	}

	// Mark the reader as EOF.
	// This is used so we don't double count EOF characters.
	if ch == eof {
		r.eof = true
	}

	return r.curr()

}

// unread pushes the previously read rune back onto the buffer.
func (r *reader) unread() {
	r.n++
}

/*curr returns the last read character and position*/
/*TODO*/
func (r reader) curr() (ch rune, pos Pos) {
	/*这里没看明白*/
	i := (r.i - r.n + len(r.buf)) % len(r.buf)
	buf := &r.buf[i]
	return buf.ch, buf.pos
}


// ScanString reads a quoted string from a rune reader.
func ScanString(r io.RuneReader) (string, error) {
	ending, _, err := r.ReadRune()
	if err != nil {
		return "", errBadString
	}

	var buf bytes.Buffer
	for {
		ch0, _, err := r.ReadRune()
		if ch0 == ending {
			return buf.String(), nil
		} else if err != nil || ch0 == '\n' {
			return buf.String(), errBadString
		} else if ch0 == '\\' {
			// If the next character is an escape then write the escaped char.
			// If it's not a valid escape then return an error.
			ch1, _, _ := r.ReadRune()
			if ch1 == 'n' {
				_, _ = buf.WriteRune('\n')
			} else if ch1 == '\\' {
				_, _ = buf.WriteRune('\\')
			} else if ch1 == '"' {
				_, _ = buf.WriteRune('"')
			} else if ch1 == '`' {
				_, _ = buf.WriteRune('`')
			} else if ch1 == '\'' {
				_, _ = buf.WriteRune('\'')
			} else {
				return string(ch0) + string(ch1), errBadEscape
			}
		} else {
			_, _ = buf.WriteRune(ch0)
		}
	}
}

var errBadString = errors.New("bad string")
var errBadEscape = errors.New("bad escape")

// ScanBareIdent reads bare identifier from a rune reader.
func ScanBareIdent(r io.RuneScanner) string {
	// Read every ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	var buf bytes.Buffer
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			break
		} else if !isIdentChar(ch) {
			r.UnreadRune()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	return buf.String()
}