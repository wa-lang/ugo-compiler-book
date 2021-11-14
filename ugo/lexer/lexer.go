package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/chai2010/ugo/token"
)

const eof = 0

type Option struct {
	SkipComment    bool
	DontInsertSemi bool
}

func Lex(name, input string, opt Option) []token.Token {
	l := &lexer{
		name:  name,
		input: input,
		opt:   opt,
	}
	l.run()

	if len(l.items) == 0 {
		l.items = append(l.items, token.Token{Type: token.EOF})
	}

	if l.items[len(l.items)-1].Type != token.EOF {
		l.items = append(l.items, token.Token{Type: token.EOF})
	}

	// return multi ';'
	items := l.items[:1]
	for _, x := range l.items[1:] {
		if x.Type == token.SEMICOLON {
			if items[len(items)-1].Type == token.SEMICOLON {
				continue
			}
		}
		items = append(items, x)
	}

	l.items = items
	return l.items
}

// lexer holds the state of the scanner.
type lexer struct {
	opt   Option
	name  string        // used only for error reports.
	input string        // the string being scanned.
	start int           // start position of this item.
	pos   int           // current position in the input.
	width int           // width of last rune read from input.
	items []token.Token // channel of scanned items.
}

// next returns the next rune in the input.
func (l *lexer) next() (rune int) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	rune, l.width = int(r), size
	l.pos += l.width
	return
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() int {
	rune := l.next()
	l.backup()
	return rune
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(typ token.TokenType) {
	tok := token.Token{
		Type:    typ,
		Literal: l.input[l.start:l.pos],
		Pos:     token.Pos(l.start + 1),
	}

	if typ == token.IDENT {
		tok.Type = token.Lookup(tok.Literal)
	}

	l.items = append(l.items, tok)
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune
// if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, rune(l.next())) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, rune(l.next())) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on. Doing it this way
// means we don't have to worry about peek double counting.
func (l *lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.pos], "\n")
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) {
	l.items = append(l.items, token.Token{
		Type:    token.ILLEGAL,
		Literal: fmt.Sprintf(format, args...),
		Pos:     token.Pos(l.start + 1),
	})
}
