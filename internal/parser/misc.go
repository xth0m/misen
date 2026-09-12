package parser

import (
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

