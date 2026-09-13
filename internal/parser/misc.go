package parser

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func (it *ParserIter) next() rune {
	if it.pos >= len(it.Source) {
		return 0
	}

	r, width := utf8.DecodeRuneInString(it.Source[it.pos:])
	it.pos += width

	if r == '\n' {
		it.loc += 1
		it.col = 1
		it.currLineStartIdx = it.pos
	} else {
		it.col += 1
	}

	return r
}

func (it *ParserIter) peek() rune {
	if it.pos >= len(it.Source) {
		return 0
	}

	r, _ := utf8.DecodeRuneInString(it.Source[it.pos:])

	return r
}

func (it *ParserIter) match(ch rune) bool {
	if it.peek() == ch {
		it.next()
		return true
	}

	return false
}

func (it *ParserIter) matchStr(s string) bool {
	if strings.HasPrefix(it.Source[it.pos:], s) {
		for range s {
			it.next()
		}

		return true
	}

	return false
}

func (it *ParserIter) skipWhitespace() {
	for {
		switch it.peek() {
		case ' ', '\t', '\n':
			it.next()

		default:
			return
		}
	}
}

func (it *ParserIter) consume(ch rune) bool {
	tempIt := *it

	tempIt.skipWhitespace()
	if tempIt.match(ch) {
		*it = tempIt
		return true
	}

	return false
}

func (it *ParserIter) until(del rune) (string, bool) {
	tempIt := *it

	for {
		ch := tempIt.peek()

		if ch == 0 || (del != '\n' && ch == '\n') {
			return "", false
		}

		if ch == del {
			break
		}

		tempIt.next()
	}

	slice := tempIt.Source[it.pos:tempIt.pos]
	*it = tempIt

	return slice, true
}

func (it *ParserIter) lexNum() (float64, bool) {
	tempIt := *it

	for tempIt.peek() >= '0' && tempIt.peek() <= '9' {
		tempIt.next()
	}

	if tempIt.pos != it.pos && (tempIt.match('.') || tempIt.match('/')) {

		for tempIt.peek() >= '0' && tempIt.peek() <= '9' {
			tempIt.next()
		}
	}

	if tempIt.pos == it.pos {
		return -1, false
	}

	n, _ := strconv.ParseFloat(it.Source[it.pos:tempIt.pos], 64)

	*it = tempIt

	return n, true
}
