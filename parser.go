package main

import "strconv"

type Expr interface{}

type Program struct {
	body []Expr
}

type IdentExpr struct {
	expr_type string
	value string
}

type AssignmentExpr struct {
	expr_type string
	name  IdentExpr
	value Expr
}

type IntExpr struct {
	expr_type string
	value int
}

type StringExpr struct {
	expr_type string
	value string	
}

func parse(tokens []Token) Program {
	program := Program{}
	skip := 0

	for skip < len(tokens) {
		node, s := expr(tokens[skip:])
		skip = s
		program.body = append(program.body, node)
	}

	return program
}

func expr(tokens []Token) (Expr, int) {
	if tokens[0].token_type == "ident" && tokens[1].token_type == "equals" {
		node, s := expr(tokens[2:])
		return AssignmentExpr{"assignment", IdentExpr{"identifier", tokens[0].value}, node}, 2 + s

	} else if tokens[0].token_type == "int" {
		num, err := strconv.Atoi(tokens[0].value)
		if err != nil {
			panic("Cannot convert int to int")
		}
		return IntExpr{"number", num}, 1

	} else if tokens[0].token_type == "text" {
		return StringExpr{"string", tokens[0].value}, 1
	}
	panic("No known expr.")
} 