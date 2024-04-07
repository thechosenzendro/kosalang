package main

import (
	"errors"
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
	expr_type  string
	true_expr  Expr
	true_block []Expr
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
func parse_atom(tokens *Tokens) (expr Expr, err error) {
	tok := tokens.curr_tok()
	if tok == (*Token)(nil) {
		return nil, errors.New("Unexpected EOF")
	}

	if tok.token_type == "open_paren" {
		tokens.advance()
		val, _, err := parse_expr(tokens, 1)
		if err != nil {
			return nil, errors.New("HelloWorld")
		}

		tok = tokens.curr_tok()
		if tok == (*Token)(nil) {
			return nil, errors.New("Unexpected EOF")
		}
		if tok.token_type != "close_paren" {
			return nil, errors.New("Unmatched '('")
		}
		tokens.advance()
		return val, nil

	} else if tok.token_type == "binop" {
		return nil, errors.New(fmt.Sprintf("Expected an atom, not %v\n", tok))
	} else if tok.token_type == "int" {
		tokens.advance()
		return tok.value, nil
	} else if tok.token_type == "ident" {
		tokens.advance()
		return IdentExpr{"identifier", tok.value}, nil
	}

	return nil, errors.New("uhoh")
}
func parse_expr(tokens *Tokens, min_prec int) (expr Expr, skip int, err error) {
	atom_lhs, err := parse_atom(tokens)
	if err != nil {
		return nil, -1, err
	}
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
			return nil, -1, errors.New("Current token isnt a binop")
		}
		op := cur.value
		prec := precedences[cur.value]
		next_min_prec := prec + 1
		tokens.advance()
		atom_rhs, _, err := parse_expr(tokens, next_min_prec)
		if err != nil {
			return nil, -1, err
		}
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
			return nil, -1, errors.New("No known binop")
		}
	}
	return atom_lhs, tokens.curr, nil
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

	}

	p := 0
	//TODO: it creates a problem
	for true {
		if p < len(tokens) && tokens[p].token_type != "eol" {
			p++
		} else {
			break
		}

	}
	t := tokens[:p]
	tks := Tokens{0, t}
	e, skip, err := parse_expr(&tks, 1)
	if err == nil {
		return e, skip
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
		p := 1
		for tokens[p].token_type != "colon" {
			p++
		}
		l := tokens[1:p]
		fmt.Println(l)
		true_expr, s := expr(l, skip)
		s = s + 3
		var true_block []Expr
		if tokens[s].token_type == "indent" {
			n := tokens[s].value
			p = s + 1
			for true {
				if tokens[p].token_type == "dedent" && tokens[p].value == n {
					break
				}
				p++
			}
			fmt.Println(tokens[p])
			block_tokens := tokens[s+1 : p]
			true_block = parse(block_tokens).body
		} else {
			panic("No block found")
		}
		return IfExpr{"if_expr", true_expr, true_block}, skip + s + p + 2
	}
	panic(fmt.Sprintf("Cannot convert token %v\n", tokens[0]))
}
