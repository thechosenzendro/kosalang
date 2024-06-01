package main

import "fmt"

func analyze(program Program) Program {
	blockConstants := map[string]string{}

	for _, expr := range program.body {
		switch expr := expr.(type) {
		case AssignmentExpr:

			_, bad := blockConstants[expr.name.value]
			if bad {
				panic(fmt.Sprintf("Cannot reassign \"%s\"", expr.name.value))
			}
			blockConstants[expr.name.value] = ""
			switch value := expr.value.(type) {
			case IntExpr, FloatExpr, AdditionExpr, SubtractionExpr, MultiplicationExpr, DivisionExpr:
			default:
				panic(fmt.Sprintf("Cannot assign %T to \"%s\"", value, expr.name.value))
			}
		}

	}

	return program
}
