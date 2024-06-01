package main

import (
	"fmt"
	"strconv"
)

type Expr interface{}

type Program struct {
	body []Expr
}

type IntExpr struct {
	value int
}

type FloatExpr struct {
	value float64
}

type IdentExpr struct {
	value string
}

type AssignmentExpr struct {
	name  IdentExpr
	value Expr
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
type EOF struct{}

func parse(ts *TokenStream) Program {
	program := Program{}
	for ts.peek(0) != nil {
		expr := parseExpr(ts)
		program.body = append(program.body, expr)
	}
	return program
}

func parseExpr(ts *TokenStream) Expr {
	if ts.peek(0).tokenType == "ident" && ts.peek(1).tokenType == "equals" {
		name := ts.peek(0).value
		ts.consume(2)
		value := parseExpr(ts)
		return AssignmentExpr{IdentExpr{name}, value}
	} else if isArithmeticExpr(ts) {
		return parseArithmeticExpr(ts, 0)
	} else if ts.peek(0).tokenType == "ident" {
		value := ts.peek(0).value
		ts.consume(1)
		return IdentExpr{value}
	} else if ts.peek(0).tokenType == "intLit" {
		return parseIntLit(ts)
	} else if ts.peek(0).tokenType == "floatLit" {
		return parseFloatLit(ts)
	} else if ts.peek(0).tokenType == "eol" {
		ts.consume(1)
		if ts.peek(0) != nil {
			return parseExpr(ts)
		} else {
			return EOF{}
		}
	}
	panic(fmt.Sprintf("Unexpected token: %v\n", ts.peek(0)))
}

func isArithmeticExpr(ts *TokenStream) bool {
	offset := 0
	if ts.peek(0).tokenType == "openParen" {
		offset = 1
	}
	return (ts.peek(offset).tokenType == "intLit" || ts.peek(offset).tokenType == "floatLit") && ts.peek(offset+1).tokenType == "binOp"
}

func parseArithmeticExpr(ts *TokenStream, minPrec int) Expr {
	atomLHS := parseAtom(ts)
	precedences := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}
	for {
		curr := ts.peek(0)
		if curr == nil || curr.tokenType != "binOp" || precedences[curr.value] < minPrec {
			break
		}
		if curr.tokenType != "binOp" {
			panic(fmt.Sprintf("Token %#v is not a binary operator", curr))
		}
		operator := curr.value
		precedence := precedences[operator]
		minPrec = precedence + 1
		ts.consume(1)
		atomRHS := parseArithmeticExpr(ts, minPrec)
		switch operator {
		case "+":
			atomLHS = AdditionExpr{atomLHS, atomRHS}
		case "-":
			atomLHS = SubtractionExpr{atomLHS, atomRHS}
		case "*":
			atomLHS = MultiplicationExpr{atomLHS, atomRHS}
		case "/":
			atomLHS = DivisionExpr{atomLHS, atomRHS}
		}
	}
	return atomLHS
}
func parseAtom(ts *TokenStream) Expr {
	if ts.peek(0) == nil {
		panic("Unexpected EOF")
	}
	if ts.peek(0).tokenType == "openParen" {
		ts.consume(1)
		expr := parseArithmeticExpr(ts, 0)
		if ts.peek(0) == nil {
			panic("Unexpected EOF")
		}
		if ts.peek(0).tokenType != "closeParen" {
			panic("Unmatched \"(\"")
		}
		ts.consume(1)
		return expr
	} else if ts.peek(0).tokenType == "intLit" {
		return parseIntLit(ts)
	} else if ts.peek(0).tokenType == "floatLit" {
		return parseFloatLit(ts)
	} else {
		panic(fmt.Sprintf("Expected an atom, not a \"%s\"", ts.peek(0).value))
	}
}

func parseFloatLit(ts *TokenStream) FloatExpr {
	value, err := strconv.ParseFloat(ts.peek(0).value, 64)
	if err != nil {
		panic(fmt.Sprintf("Cannot convert \"%s\" to float", ts.peek(0).value))
	}
	ts.consume(1)
	return FloatExpr{value}
}

func parseIntLit(ts *TokenStream) IntExpr {
	value, err := strconv.Atoi(ts.peek(0).value)
	if err != nil {
		panic(fmt.Sprintf("Cannot convert \"%s\" to int", ts.peek(0).value))
	}
	ts.consume(1)
	return IntExpr{value}
}
