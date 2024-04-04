package main

import "unicode"

func tokenize(data string) []Token {
	rune_data := []rune(data + "\n") 
	var tokens []Token
	var buffer string
	var skip = -1
	// Added whitespace after file content so that every token gets added to the tokens list
	for pos := 0; pos < len(rune_data); pos++ {
		char := rune_data[pos]
		if skip != -1 {
			if skip >= len(data)-1 {
				break
			}
			if pos < skip {
				continue
			} else {
				skip = -1
			}
		}

		if unicode.IsLetter(char) || char == '_' {
			p := pos
			for unicode.IsLetter(rune_data[p]) || rune_data[p] == '_' {
				buffer += string(rune_data[p])
				p++
			}
			tokens = append(tokens, Token{"ident", buffer})
			buffer = ""
			skip = p
		} else if char == '=' {
			tokens = append(tokens, Token{"equals", "="})
		} else if char == '"' {
			// normal pos would capture the first quote
			p := pos + 1
			for pos < len(data)-1 && rune_data[p] != '"' {
				buffer += string(rune_data[p])
				p++
			}
			tokens = append(tokens, Token{"text", buffer})
			buffer = ""
			// skip + 1 so that it doesnt skip to the closing quote and create another text
			skip = p + 1

		} else if unicode.IsNumber(char) {
			p := pos
			for unicode.IsNumber(rune_data[p]) || rune_data[p] == '_' {
				if rune_data[p] == '_' {
					p++
					continue
				}
				buffer += string(rune_data[p])
				p++
			}
			tokens = append(tokens, Token{"int", buffer})
			buffer = ""
			skip = p

		} else if char == '\n' {
			tokens = append(tokens, Token{"eol", "\\n"})
		} else if char == ' ' {
			continue
		} else if char == '#' {
			p := pos
			for rune_data[p] != '\n' {
				p++
			}
			skip = p
		} else if char == '+' {
			tokens = append(tokens, Token{"plus", "+"})
		} else if char == '-' {
			tokens = append(tokens, Token{"minus", "-"})
		} else if char == '*' {
			tokens = append(tokens, Token{"asterisk", "*"})
		} else if char == '/' {
			tokens = append(tokens, Token{"slash", "/"})
		} else {
			panic("Unexpected token: " + string(char))
		}

	}
	return tokens
}