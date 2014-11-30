package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"
	"unicode/utf8"
)

type lexer struct {
	input string
	pos   Pos
	start Pos
	width Pos
	items chan item
	state stateFn
}

type item struct {
	typ itemType
	pos Pos
	val string
}

type Pos int

type itemType int

const (
	itemError itemType = iota
	itemText
	itemEOF
	itemSpace
	itemIdentifier
	itemClass
	itemDef
	itemEnd
)

// not sure why this is a var an not a const
var keywords = map[string]itemType{
	"class": itemClass,
	"def":   itemDef,
	"end":   itemEnd,
}

const eof = -1

const lineComment = "//"

type stateFn func(*lexer) stateFn

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	fmt.Printf("running - -")
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) next() rune {
		if int(l.pos) >= len(l.input) {
			l.width = 0
			return eof
		}
		r, w := utf8.DecodeRuneInString(l.input[l.pos:])
		l.width = Pos(w)
		l.pos += l.width
		return r
	}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func isAlphaNumeric(r rune) bool {
	return r == '-' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

////////////////////////////////////////
// The main lexer, the center of this statemachine
func lexText(l *lexer) stateFn {
	switch r := l.next(); {
	case isSpace(r):
		return lexSpace
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifier
	case strings.HasPrefix(l.input[l.pos:], lineComment):
		return lexLineComment
	}

	// inform the channel we're done
	l.emit(itemEOF)
	return nil
}

func lexLineComment(l *lexer) stateFn {
	l.pos += Pos(len(lineComment))
	i := strings.Index(l.input[l.pos:], "\n")
	// if -1 it's the last line of the program and we don't care
	l.pos += Pos(i + len("\n"))
	return lexText
}

func lexSpace(l *lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(itemSpace)
	return lexText
}

func lexIdentifier(l *lexer) stateFn {
	// yuck, this Loop is because of the nested for>switch
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			switch {
				case keywords[word] > 0:
				  l.emit(keywords[word])
			default:
				l.emit(itemIdentifier)
			}
			break Loop
		}
	}
	return lexText
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

////////////////////////////////////////
// usage

func main() {
	file := loadFile("test.rb")
	lexer := lex(file)

	for item := range lexer.items {
		fmt.Printf("%v, %v\n", item.val, item.typ)
	}

//	for _,token := range tokens {
//		fmt.Printf("%T, \t %v\n", token, token)
//	}
}

func loadFile(filename string) (string) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error %s", err)
	}

	return string(body)
}
