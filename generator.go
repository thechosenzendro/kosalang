package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func generate(program Program) *ir.Module {
	module := ir.NewModule()

	for _, expr := range program.body {
		switch expr := expr.(type) {
		case AssignmentExpr:
			var value constant.Constant
			switch assignmentValue := expr.value.(type) {
			case IntExpr:
				value = convertIntExpr(assignmentValue)

			case FloatExpr:
				value = convertFloatExpr(assignmentValue)

			case AdditionExpr, SubtractionExpr, MultiplicationExpr, DivisionExpr:
				value = convertArithmeticExpr(module, assignmentValue)
			default:
				panic(fmt.Sprintf("%T cannot be converted into IR", assignmentValue))
			}
			module.NewGlobalDef(expr.name.value, value)
		}
	}
	main := module.NewFunc("main", types.I32)
	entry := main.NewBlock("")
	entry.NewRet(constant.NewInt(types.I32, 0))
	return module
}

func convertArithmeticExpr(module *ir.Module, expr Expr) constant.Constant {
	switch expr := expr.(type) {
	case AdditionExpr:

		var lhs value.Value
		lhs = convertArithmeticExpr(module, expr.lhs)

		switch _lhs := lhs.(type) {
		case *ir.Global:
			lhs = ir.NewLoad(_lhs.Type(), _lhs)
		default:
			lhs = _lhs
		}

		var rhs value.Value
		rhs = convertArithmeticExpr(module, expr.rhs)

		switch _rhs := rhs.(type) {
		case *ir.Global:
			rhs = ir.NewLoad(_rhs.Type(), _rhs)
		default:
			rhs = _rhs
		}
		name := fmt.Sprintf("tmp%d", rand.Int())
		variable := module.NewGlobalDef(name, ir.NewAdd(lhs, rhs))
		return variable
	case SubtractionExpr:
		lhs := convertArithmeticExpr(module, expr.lhs)
		rhs := convertArithmeticExpr(module, expr.rhs)
		name := fmt.Sprintf("tmp%d", rand.Int())
		variable := module.NewGlobalDef(name, constant.NewSub(lhs, rhs))
		return variable
	case MultiplicationExpr:
		lhs := convertArithmeticExpr(module, expr.lhs)
		rhs := convertArithmeticExpr(module, expr.rhs)
		name := fmt.Sprintf("tmp%d", rand.Int())
		variable := module.NewGlobalDef(name, constant.NewMul(lhs, rhs))
		return variable
	case DivisionExpr:
		lhs := convertArithmeticExpr(module, expr.lhs)
		rhs := convertArithmeticExpr(module, expr.rhs)
		name := fmt.Sprintf("tmp%d", rand.Int())
		variable := module.NewGlobalDef(name, constant.NewSDiv(lhs, rhs))
		return variable
	default:
		return convertToConstant(expr)
	}
}

func convertToConstant(expr Expr) constant.Constant {
	switch expr := expr.(type) {
	case IntExpr:
		return convertIntExpr(expr)
	case FloatExpr:
		return convertFloatExpr(expr)
	default:
		panic(fmt.Sprintf("Cannot convert %T to a constant", expr))
	}

}

func convertFloatExpr(expr FloatExpr) *constant.Float {
	return constant.NewFloat(types.Float, expr.value)
}

func convertIntExpr(expr IntExpr) *constant.Int {
	return constant.NewInt(types.I64, int64(expr.value))
}

func hello() {
	m := ir.NewModule()

	globalG := m.NewGlobalDef("g", constant.NewInt(types.I32, 2))

	funcAdd := m.NewFunc("add", types.I32,
		ir.NewParam("x", types.I32),
		ir.NewParam("y", types.I32),
	)
	ab := funcAdd.NewBlock("")
	ab.NewRet(ab.NewAdd(funcAdd.Params[0], ir.NewLoad(types.I32, globalG)))

	funcMain := m.NewFunc(
		"main",
		types.I32,
	) // omit parameters
	mb := funcMain.NewBlock("") // llir/llvm would give correct default name for block without name
	mb.NewRet(mb.NewCall(funcAdd, constant.NewInt(types.I32, 1)))

	println(m.String())
}
