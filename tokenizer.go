package main

import (
	"fmt"
	"unicode"
)

type CharStream struct {
	index int
	chars []rune
}

func (cs *CharStream) peek(offset int) *rune {
	index := cs.index + offset
	if index < len(cs.chars) {
		return &cs.chars[index]
	} else {
		return nil
	}
}

func (cs *CharStream) consume(amount int) {
	cs.index += amount
}

type Token struct {
	tokenType string
	value     string
}

type TokenStream struct {
	index  int
	tokens []Token
}

func (ts *TokenStream) peek(offset int) *Token {
	index := ts.index + offset
	if index < len(ts.tokens) {
		return &ts.tokens[index]
	} else {
		return nil
	}
}

func (ts *TokenStream) consume(amount int) {
	ts.index += amount
}

func tokenize(sourceCode string) *TokenStream {
	cs := &CharStream{index: 0, chars: []rune(sourceCode)}
	ts := &TokenStream{index: 0, tokens: []Token{}}
	for {
		if cs.peek(0) != nil {
			lex(cs, ts)
		} else {
			break
		}
	}
	return ts
}

func lex(cs *CharStream, ts *TokenStream) {
	char := *cs.peek(0)

	if unicode.IsNumber(char) || char == '-' {
		buffer := []rune{}
		isFloat := false
		if char == '-' {
			buffer = append(buffer, char)
			cs.consume(1)
		}
		//prevents the -.1 syntax
		if *cs.peek(0) == '.' {
			panic("Missing whole part of the float literal")
		}
		for unicode.IsNumber(*cs.peek(0)) || *cs.peek(0) == '_' || *cs.peek(0) == '.' {
			if *cs.peek(0) == '.' {
				if isFloat {
					panic("There can only be one dot in float literals")
				} else {
					buffer = append(buffer, '.')
					isFloat = true
					cs.consume(1)
				}
			}
			if *cs.peek(0) == '_' {
				cs.consume(1)
			} else if unicode.IsNumber(*cs.peek(0)) {
				buffer = append(buffer, *cs.peek(0))
				cs.consume(1)
			}
		}
		if isFloat {
			ts.tokens = append(ts.tokens, Token{"floatLit", string(buffer)})
		} else {
			ts.tokens = append(ts.tokens, Token{"intLit", string(buffer)})
		}
	} else if unicode.IsLetter(char) || char == '_' {
		buffer := ""
		for unicode.IsLetter(*cs.peek(0)) || *cs.peek(0) == '_' || unicode.IsNumber(*cs.peek(0)) {
			buffer += string(*cs.peek(0))
			cs.consume(1)
		}
		ts.tokens = append(ts.tokens, Token{"ident", buffer})
	} else if char == '#' {
		cs.consume(1)
		for *cs.peek(0) != '\n' {
			cs.consume(1)
		}
		cs.consume(1)
	} else if char == ' ' {
		cs.consume(1)
	} else if char == '=' {
		ts.tokens = append(ts.tokens, Token{"equals", "="})
		cs.consume(1)
	} else if char == '(' {
		ts.tokens = append(ts.tokens, Token{"openParen", "("})
		cs.consume(1)
	} else if char == ')' {
		ts.tokens = append(ts.tokens, Token{"closeParen", ")"})
		cs.consume(1)
	} else if char == '+' || char == '-' || char == '*' || char == '/' {
		ts.tokens = append(ts.tokens, Token{"binOp", string(char)})
		cs.consume(1)
	} else if char == '\n' {
		ts.tokens = append(ts.tokens, Token{"eol", "\\n"})
		cs.consume(1)
	} else {
		panic(fmt.Sprintf("Unexpected: %s", string(char)))
	}
}
