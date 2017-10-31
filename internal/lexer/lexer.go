package lexer

import (
	"strings"

	"fmt"
)

type TokenType string

const (
	TTArg        TokenType = "Arg"
	TTOpenPar    TokenType = "OpenPar"
	TTClosePar   TokenType = "ClosePar"
	TTOpenSq     TokenType = "OpenSq"
	TTCloseSq    TokenType = "CloseSq"
	TTChoice     TokenType = "Choice"
	TTOptions    TokenType = "Options"
	TTRep        TokenType = "Rep"
	TTShortOpt   TokenType = "ShortOpt"
	TTLongOpt    TokenType = "LongOpt"
	TTOptSeq     TokenType = "OptSeq"
	TTOptValue   TokenType = "OptValue"
	TTDoubleDash TokenType = "DblDash"
)

type Token struct {
	Typ TokenType
	Val string
	Pos int
}

func (t *Token) String() string {
	return fmt.Sprintf("%s('%s')@%d", t.Typ, t.Val, t.Pos)
}

type ParseError struct {
	Input string
	Msg   string
	Pos   int
}

func (t *ParseError) ident() string {
	return strings.Map(func(c rune) rune {
		switch c {
		case '\t':
			return c
		default:
			return ' '
		}
	}, t.Input[:t.Pos])
}

func (t *ParseError) Error() string {
	return fmt.Sprintf("Parse error at position %d:\n%s\n%s^ %s",
		t.Pos, t.Input, t.ident(), t.Msg)
}

func Tokenize(usage string) ([]*Token, error) {
	pos := 0
	res := []*Token{}
	var (
		tk = func(t TokenType, v string) {
			res = append(res, &Token{t, v, pos})
		}

		tkp = func(t TokenType, v string, p int) {
			res = append(res, &Token{t, v, p})
		}

		err = func(msg string) *ParseError {
			return &ParseError{usage, msg, pos}
		}
	)
	eof := len(usage)
	for pos < eof {
		switch c := usage[pos]; c {
		case ' ':
			pos++
		case '\t':
			pos++
		case '[':
			tk(TTOpenSq, "[")
			pos++
		case ']':
			tk(TTCloseSq, "]")
			pos++
		case '(':
			tk(TTOpenPar, "(")
			pos++
		case ')':
			tk(TTClosePar, ")")
			pos++
		case '|':
			tk(TTChoice, "|")
			pos++
		case '.':
			start := pos
			pos++
			if pos >= eof || usage[pos] != '.' {
				return nil, err("Unexpected end of usage, was expecting '..'")
			}
			pos++
			if pos >= eof || usage[pos] != '.' {
				return nil, err("Unexpected end of usage, was expecting '.'")
			}
			tkp(TTRep, "...", start)
			pos++
		case '-':
			start := pos
			pos++
			if pos >= eof {
				return nil, err("Unexpected end of usage, was expecting an option name")
			}

			switch o := usage[pos]; {
			case isLetter(o):
				pos++
				for ; pos < eof; pos++ {
					ok := isLetter(usage[pos])
					if !ok {
						break
					}
				}
				typ := TTShortOpt
				opt := usage[start:pos]
				if pos-start > 2 {
					typ = TTOptSeq
					opt = opt[1:]
				}
				tkp(typ, opt, start)
				if pos < eof && usage[pos] == '-' {
					return nil, err("Invalid syntax")
				}
			case o == '-':
				pos++
				if pos == eof || usage[pos] == ' ' {
					tkp(TTDoubleDash, "--", start)
					continue
				}
				for pos0 := pos; pos < eof; pos++ {
					ok := isOkLongOpt(usage[pos], pos == pos0)
					if !ok {
						break
					}
				}
				opt := usage[start:pos]
				if len(opt) == 2 {
					return nil, err("Was expecting a long option name")
				}
				tkp(TTLongOpt, opt, start)
			}

		case '=':
			start := pos
			pos++
			if pos >= eof || usage[pos] != '<' {
				return nil, err("Unexpected end of usage, was expecting '=<'")
			}
			closed := false
			for ; pos < eof; pos++ {
				closed = usage[pos] == '>'
				if closed {
					break
				}
			}
			if !closed {
				return nil, err("Unclosed option value")
			}
			if pos-start == 2 {
				return nil, err("Was expecting an option value")
			}
			pos++
			value := usage[start:pos]

			tkp(TTOptValue, value, start)

		default:
			switch {
			case isUppercase(c):
				start := pos
				for pos = pos + 1; pos < eof; pos++ {
					if !isOkInArg(usage[pos]) {
						break
					}
				}
				s := usage[start:pos]
				typ := TTArg
				if s == "OPTIONS" {
					typ = TTOptions
				}
				tkp(typ, s, start)
			default:
				return nil, err("Unexpected input")
			}

		}
	}

	return res, nil
}

func isLowercase(c uint8) bool {
	return c >= 'a' && c <= 'z'
}

func isUppercase(c uint8) bool {
	return c >= 'A' && c <= 'Z'
}

func isOkInArg(c uint8) bool {
	return isUppercase(c) || isDigit(c) || c == '_'
}

func isLetter(c uint8) bool {
	return isLowercase(c) || isUppercase(c)
}

func isDigit(c uint8) bool {
	return c >= '0' && c <= '9'
}
func isOkLongOpt(c uint8, first bool) bool {
	return isLetter(c) || isDigit(c) || c == '_' || (!first && c == '-')
}
