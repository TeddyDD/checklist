package checklist

import (
	"fmt"
	"unicode/utf8"
)

//go:generate stringer -type=tokenType
type tokenType int

const (
	tokenEOF tokenType = iota
	tokenHypen
	tokenText
	tokenCommand
	tokenError
)

const (
	hypen   rune = '-'
	newline rune = '\n'
	grave   rune = '`'
)

const eof = -1

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	switch t.typ {
	case tokenEOF:
		return "Token<EOF>"
	case tokenError:
		return fmt.Sprintf("Token<error>: %s", t.val)
	default:
		return fmt.Sprintf("Token-%s<%s>", t.typ, t.val)
	}
}

type stateFn func(*lexer) stateFn

type lexer struct {
	pos    int
	start  int
	width  int
	input  string
	tokens chan token
}

func (l *lexer) run() {
	for state := lexTopLevel; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) nextToken() token {
	return <-l.tokens
}

func (l *lexer) next() (r rune) {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) emit(t tokenType) {
	val := l.input[l.start:l.pos]
	tok := token{
		typ: t,
		val: val,
	}
	l.start = l.pos
	l.tokens <- tok
}

func (l *lexer) errorf(format string, vals ...interface{}) {
	tok := token{
		typ: tokenError,
		val: fmt.Sprintf(format, vals...),
	}
	l.tokens <- tok
}

func newLexer(input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make(chan token, 2),
	}
	go l.run()
	return l
}

func lexTopLevel(l *lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case '-':
			return lexHypen
		case eof:
			l.emit(tokenEOF)
			return nil
		default:
			l.ignore()
		}
	}
}

func lexHypen(l *lexer) stateFn {
	l.emit(tokenHypen)
	return lexText
}

func lexText(l *lexer) stateFn {
	for {
		r := l.next()
		switch r {
		case eof:
			if l.pos > l.start {
				l.backup()
				l.emit(tokenText)
				l.emit(tokenEOF)
			}
			return nil

		case newline:
			p := l.peek()
			if p == hypen {
				if l.pos > l.start {
					l.emit(tokenText)
				}
				return lexTopLevel
			}
		case grave:
			l.backup()
			if l.pos > l.start {
				l.emit(tokenText)
			}
			return lexCommand
		}
	}
}

func lexCommand(l *lexer) stateFn {
	l.next() // get first grave quote
	for {
		switch r := l.next(); r {
		case grave:
			l.emit(tokenCommand)
			return lexTopLevel
		case eof:
			l.errorf("unexpected EOF in command: %s", l.input[l.start:])
			return nil
		}
	}
	panic("Lexer invalid state")
}
