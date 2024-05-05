package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TextStream struct {
	curr int
	data []rune
}

func (ts TextStream) peek(offset int) (token *rune) {
	if ts.curr+offset < len(ts.data) {
		return &ts.data[ts.curr+offset]
	} else {
		return nil
	}
}
func (ts *TextStream) consume(amount int) {
	ts.curr += amount
}

var tokens []Token
var buffer string
var indent_level = 0
var indent_levels = []int{0}

var paren_levels = 0
var curly_bracket_levels = 0

func tokenize(raw_text string) []Token {
	text := []rune(strings.Replace(strings.Replace(raw_text+"\n", "\r\n", "\n", -1), "\t", "    ", -1))
	ts := TextStream{0, text}
	for ts.peek(0) != nil {
		tokens = lex_chars(&ts, tokens)
	}
	return tokens
}

func lex_chars(ts *TextStream, tokens []Token) []Token {
	char := *ts.peek(0)
	if unicode.IsLetter(char) || char == '_' {
		for unicode.IsLetter(*ts.peek(0)) || *ts.peek(0) == '_' || *ts.peek(0) == '.' {
			buffer += string(*ts.peek(0))
			ts.consume(1)
		}
		if buffer == "if" ||
			buffer == "true" ||
			buffer == "false" ||
			buffer == "else" ||
			buffer == "and" ||
			buffer == "or" ||
			buffer == "data" ||
			buffer == "pub" ||
			buffer == "using" ||
			buffer == "match" {
			tokens = append(tokens, Token{"keyword", buffer})
		} else if buffer == "f" && *ts.peek(0) == '"' {
			ts.consume(1)
			for *ts.peek(0) != '"' {
				if *ts.peek(0) == '{' {
					ts.consume(1)
					// finish
				} else {
					buffer += string(*ts.peek(0))
					ts.consume(1)
				}
			}
		} else {
			tokens = append(tokens, Token{"ident", buffer})
		}
		buffer = ""
	} else if char == '=' {
		ts.consume(1)
		tokens = append(tokens, Token{"equals", "="})
	} else if char == '"' {
		ts.consume(1)
		for *ts.peek(0) != '"' {
			fmt.Println(string(*ts.peek(0)))
			buffer += string(*ts.peek(0))
			ts.consume(1)
		}
		ts.consume(1)
		tokens = append(tokens, Token{"string_lit", buffer})
		buffer = ""
	} else if unicode.IsNumber(*ts.peek(0)) {
		is_float := false
		for unicode.IsNumber(*ts.peek(0)) || *ts.peek(0) == '_' || *ts.peek(0) == '.' {
			if *ts.peek(0) == '.' {
				if is_float {
					panic("Cannot have two dots in float literal.")
				} else {
					is_float = true
				}
				ts.consume(1)
			}
			if *ts.peek(0) == '_' {
				ts.consume(1)
				continue
			}
			buffer += string(*ts.peek(0))
			ts.consume(1)
		}
		if !is_float {
			tokens = append(tokens, Token{"int_lit", buffer})
		} else {
			tokens = append(tokens, Token{"float_lit", buffer})
		}
		buffer = ""
	} else if char == '\n' {
		tokens = append(tokens, Token{"eol", "\\n"})
		curr_indent := 0
		ts.consume(1)
		if ts.curr+1 >= len(ts.data) {
			return tokens
		}
		for *ts.peek(0) == ' ' {
			curr_indent++
			ts.consume(1)
		}
		if curr_indent > indent_level {
			tokens = append(tokens, Token{"indent", strconv.Itoa(curr_indent)})
			indent_levels = append(indent_levels, curr_indent)
			indent_level = curr_indent
			return tokens
		}
		for curr_indent < indent_level {
			if len(indent_levels) != 1 {
				indent_levels = indent_levels[:len(indent_levels)-1]
			}
			if indent_level == 0 {
				break
			}
			tokens = append(tokens, Token{"dedent", strconv.Itoa(indent_level)})
			indent_level = indent_levels[len(indent_levels)-1]
			if indent_level < curr_indent {
				panic("Wut?")
			}
		}
	} else if char == ' ' {
		ts.consume(1)
		return tokens
	} else if char == '#' {
		for *ts.peek(0) != '\n' {
			ts.consume(1)
		}
	} else if char == '+' {
		ts.consume(1)
		tokens = append(tokens, Token{"binop", "+"})
	} else if char == '-' {
		ts.consume(1)
		tokens = append(tokens, Token{"binop", "-"})
	} else if char == '*' {
		ts.consume(1)
		tokens = append(tokens, Token{"binop", "*"})
	} else if char == '/' {
		ts.consume(1)
		tokens = append(tokens, Token{"binop", "/"})
	} else if char == '(' {
		paren_levels++
		ts.consume(1)
		tokens = append(tokens, Token{"open_paren", strconv.Itoa(paren_levels)})
	} else if char == ')' {
		tokens = append(tokens, Token{"close_paren", strconv.Itoa(paren_levels)})
		ts.consume(1)
		paren_levels--
	} else if char == ':' {
		ts.consume(1)
		tokens = append(tokens, Token{"colon", ":"})
	} else if char == ',' {
		ts.consume(1)
		tokens = append(tokens, Token{"comma", ","})
	} else if char == '{' {
		ts.consume(1)
		curly_bracket_levels++
		tokens = append(tokens, Token{"open_curly_bracket", strconv.Itoa(curly_bracket_levels)})
	} else if char == '}' {
		ts.consume(1)
		tokens = append(tokens, Token{"close_curly_bracket", strconv.Itoa(curly_bracket_levels)})
		curly_bracket_levels--
	} else {
		panic(fmt.Sprintf("Unexpected token: %s", string(char)))
	}
	return tokens
}
