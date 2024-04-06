package main

import (
	"fmt"
	"strconv"
)

type Expr interface{}

type Program struct {
	body []Expr
}

type BoolExpr struct {
	expr_type string
	value     bool
}

type IdentExpr struct {
	expr_type string
	value     string
}

type AssignmentExpr struct {
	expr_type string
	name      IdentExpr
	value     Expr
}

type IntExpr struct {
	expr_type string
	value     int
}

type StringExpr struct {
	expr_type string
	value     string
}

type AdditionExpr struct {
	lhs Expr
	rhs Expr
}

type SubtractionExpr struct {
	lhs Expr
	rhs Expr
}

type MultiplicationExpr struct {
	lhs Expr
	rhs Expr
}

type DivisionExpr struct {
	lhs Expr
	rhs Expr
}

type Tokens struct {
	curr   int
	tokens []Token
}

type IfExpr struct {
	expr_type string
	true_expr Expr
}

type EmptyExpr struct{}

func (tks Tokens) curr_tok() *Token {
	if tks.curr < len(tks.tokens) {
		return &tks.tokens[tks.curr]
	} else {
		return nil
	}
}
func (tks *Tokens) advance() {
	tks.curr++
}
func unexpected_eof(thing any) {
	if thing == (*Token)(nil) {
		panic("Unexpected EOF")
	}
}
func parse_atom(tokens *Tokens) Expr {
	tok := tokens.curr_tok()
	unexpected_eof(tok)
	if tok.token_type == "open_paren" {
		tokens.advance()
		val, _ := parse_expr(tokens, 1)

		tok = tokens.curr_tok()
		unexpected_eof(tok)

		if tok.token_type != "close_paren" {
			panic("Unmatched '('")
		}
		tokens.advance()
		return val
	} else if tok.token_type == "binop" {
		panic(fmt.Sprintf("Expected an atom, not %v\n", tok))
	} else if tok.token_type == "int" {
		tokens.advance()
		return tok.value
	}

	panic("UH OH")
}
func parse_expr(tokens *Tokens, min_prec int) (Expr, int) {
	atom_lhs := parse_atom(tokens)
	precedences := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}
	for true {
		cur := tokens.curr_tok()

		if cur == (*Token)(nil) || cur.token_type != "binop" || precedences[cur.value] < min_prec {
			break
		}
		if cur.token_type != "binop" {
			panic("Current token isnt a binop")
		}
		op := cur.value
		prec := precedences[cur.value]
		next_min_prec := prec + 1
		tokens.advance()
		atom_rhs, _ := parse_expr(tokens, next_min_prec)
		switch op {
		case "+":
			atom_lhs = AdditionExpr{atom_lhs, atom_rhs}

		case "-":
			atom_lhs = SubtractionExpr{atom_lhs, atom_rhs}

		case "*":
			atom_lhs = MultiplicationExpr{atom_lhs, atom_rhs}

		case "/":
			atom_lhs = DivisionExpr{atom_lhs, atom_rhs}
		default:
			panic("No known binop")
		}
	}
	return atom_lhs, tokens.curr
}

func parse(tokens []Token) Program {
	program := Program{}
	skip := 0

	for skip < len(tokens) {
		x := tokens[skip:]
		node, s := expr(x, skip)
		if s > skip {
			skip = s
		} else {
			skip = skip + s
		}
		program.body = append(program.body, node)
	}

	return program
}

func get(tokens []Token, i int) Token {
	if i < len(tokens) {
		return tokens[i]
	} else {
		return Token{"", ""}
	}
}

func expr(tokens []Token, skip int) (Expr, int) {
	if get(tokens, 0).token_type == "ident" && get(tokens, 1).token_type == "equals" {
		node, s := expr(tokens[2:], skip)
		return AssignmentExpr{"assignment", IdentExpr{"identifier", tokens[0].value}, node}, 3 + s + skip

	} else if get(tokens, 1).token_type == "binop" {
		p := 0
		for true {
			if p < len(tokens) && tokens[p].token_type != "eol" {
				p++
			} else {
				break
			}

		}
		t := tokens[:p]
		tks := Tokens{0, t}
		e, skip := parse_expr(&tks, 1)
		return e, skip - 1

	} else if get(tokens, 0).token_type == "int" {
		num, err := strconv.Atoi(tokens[0].value)
		if err != nil {
			panic("Cannot convert int to int")
		}
		return IntExpr{"number", num}, 1

	} else if get(tokens, 0).token_type == "text" {
		return StringExpr{"string", tokens[0].value}, 1
	} else if get(tokens, 0).token_type == "eol" {
		return EmptyExpr{}, skip + 1
	} else if get(tokens, 0).token_type == "ident" {
		return IdentExpr{"identifier", tokens[0].value}, 1
	} else if get(tokens, 0).token_type == "keyword" && (get(tokens, 0).value == "true" || get(tokens, 0).value == "false") {
		b := tokens[0].value
		switch b {
		case "true":
			return BoolExpr{"bool", true}, 1
		case "false":
			return BoolExpr{"bool", false}, 1
		}
	} else if get(tokens, 0).token_type == "keyword" && get(tokens, 0).value == "if" {
		// finish this pls
		true_expr, s := expr(tokens[1:], skip)
		return IfExpr{"if_expr", true_expr}, skip + s
	}
	panic(fmt.Sprintf("Cannot convert token %v\n", tokens[0]))
}
