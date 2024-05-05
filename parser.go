package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Expr interface{}

type Program struct {
	body []Expr
}

type BoolExpr struct {
	value bool
}

type IdentExpr struct {
	value string
}

type AssignmentExpr struct {
	name  IdentExpr
	value Expr
}

type IntExpr struct {
	value int
}

type StringExpr struct {
	value string
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

type Block struct {
	body []Expr
}

type Tokens struct {
	curr   int
	tokens []Token
}

type IfExpr struct {
	true_expr  Expr
	true_block Block
	else_block Block
}

type FunctionDeclExpr struct {
	return_type Expr
	arguments   []ArgumentDeclExpr
	body        Block
}

type ArgumentDeclExpr struct {
	arg_type     Expr
	name         IdentExpr
	default_expr *Expr
}

type FuncCallExpr struct {
	name      IdentExpr
	arguments []Expr
}

type StructExpr struct {
	fields Block
}

type Import struct {
	path    []IdentExpr
	imports []IdentExpr
}

type UsingExpr struct {
	imports []Import
}

type PubExpr struct {
	expr Expr
}

type EmptyExpr struct{}

func (tks Tokens) peek(offset int) (token *Token) {
	if tks.curr+offset < len(tks.tokens) {
		return &tks.tokens[tks.curr+offset]
	} else {
		return nil
	}
}
func (tks *Tokens) consume(amount int) {
	tks.curr += amount
}

func keyword(token *Token) bool {
	return token.token_type == "keyword"
}

func parse(tks []Token) Program {
	program := Program{}
	tokens := Tokens{0, tks}
	for tokens.peek(0) != nil {
		expr := parse_expr(&tokens)
		program.body = append(program.body, expr)
	}
	return program
}

func parse_argument_decl(tokens *Tokens) ArgumentDeclExpr {
	if tokens.peek(0).token_type == "ident" && tokens.peek(1).token_type == "ident" {
		arg_type := tokens.peek(0)
		name := IdentExpr{tokens.peek(1).value}
		tokens.consume(2)
		if tokens.peek(0).token_type == "equals" {
			tokens.consume(1)
			expr := parse_expr(tokens)
			return ArgumentDeclExpr{arg_type, name, &expr}

		} else {
			return ArgumentDeclExpr{arg_type, name, nil}
		}
	}
	panic("Current token isnt an argument")
}

func parse_expr(tokens *Tokens) Expr {
	if tokens.peek(0).token_type == "ident" && tokens.peek(1).token_type == "equals" {
		name := tokens.peek(0).value
		tokens.consume(2)
		expr := parse_expr(tokens)
		return AssignmentExpr{IdentExpr{name}, expr}

	} else if tokens.peek(0).token_type == "ident" && ((tokens.peek(1).token_type == "ident" && tokens.peek(2).token_type == "open_paren") || tokens.peek(1).token_type == "open_paren") {
		var token_skip int
		var is_anon bool
		var paren_level string
		if tokens.peek(1).token_type == "ident" {
			token_skip = 2
			is_anon = false
			paren_level = tokens.peek(2).value
		} else if tokens.peek(1).token_type == "open_paren" {
			token_skip = 1
			is_anon = true
			paren_level = tokens.peek(1).value
		}
		for {
			if tokens.peek(token_skip).token_type == "close_paren" && tokens.peek(token_skip).value == paren_level {
				token_skip++
				break
			}
			token_skip++
		}
		if tokens.peek(token_skip).token_type == "colon" {
			if !is_anon {
				return AssignmentExpr{IdentExpr{tokens.peek(1).value}, parse_function_decl(tokens)}
			} else {
				return parse_function_decl(tokens)
			}

		} else {
			return parse_function_call(tokens)
		}
	}
	p := tokens.curr
	expr, err := parse_math_expr(tokens, 1)
	if err == nil {
		return expr
	} else {
		tokens.curr = p
	}
	if tokens.peek(0).token_type == "int_lit" {
		return parse_int_lit(tokens)
	} else if tokens.peek(0).token_type == "string_lit" {
		text := StringExpr{tokens.peek(0).value}
		tokens.consume(1)
		return text

	} else if tokens.peek(0).token_type == "eol" {
		expr := EmptyExpr{}
		tokens.consume(1)
		return expr

	} else if tokens.peek(0).token_type == "ident" {
		return parse_ident(tokens)
	} else if keyword(tokens.peek(0)) {
		switch tokens.peek(0).value {
		case "true":
			expr := BoolExpr{true}
			tokens.consume(1)
			return expr

		case "false":
			expr := BoolExpr{false}
			tokens.consume(1)
			return expr

		case "if":
			tokens.consume(1)
			true_expr := parse_expr(tokens)
			true_block := parse_block(tokens)
			var else_block Block
			if keyword(tokens.peek(0)) && tokens.peek(0).value == "else" {
				tokens.consume(1)
				else_block = parse_block(tokens)
			}
			return IfExpr{true_expr, true_block, else_block}
		case "data":
			// this doesnt handle the cases where the struct is anon
			tokens.consume(1)
			var struct_name string
			if tokens.peek(0).token_type == "ident" {
				struct_name = tokens.peek(0).value
				tokens.consume(1)
			}
			struct_block := parse_struct_block(tokens)
			if struct_name == "" {
				return StructExpr{struct_block}
			} else {
				return AssignmentExpr{IdentExpr{struct_name}, StructExpr{struct_block}}
			}
		case "using":
			tokens.consume(1)
			using_expr := UsingExpr{}
			for {
				if tokens.peek(0).token_type == "ident" {
					path := tokens.peek(0).value
					tokens.consume(1)

					path_list := strings.Split(path, ".")
					fmt.Println(path_list)

					var ident_path []IdentExpr
					for _, path_segment := range path_list {
						if path_segment != "" {
							ident_path = append(ident_path, IdentExpr{path_segment})
						}
					}
					var imports []IdentExpr
					if path_list[len(path_list)-1] == "" {
						if tokens.peek(0).token_type != "open_paren" {
							panic("uhoh")
						}
						tokens.consume(1)
						for {
							if tokens.peek(0).token_type == "ident" {
								imports = append(imports, IdentExpr{tokens.peek(0).value})
								tokens.consume(1)
							}
							if tokens.peek(0).token_type == "comma" {
								tokens.consume(1)
							}
							if tokens.peek(0).token_type == "close_paren" {
								tokens.consume(1)
								break
							}
						}
					} else {
						imports = []IdentExpr{{"*"}}
					}
					using_expr.imports = append(using_expr.imports, Import{ident_path, imports})
				}
				if tokens.peek(0).token_type == "comma" {
					tokens.consume(1)
				}
				if tokens.peek(0).token_type == "eol" {
					tokens.consume(1)
					break
				}
			}

			print_struct(using_expr)
			return using_expr
		case "pub":
			tokens.consume(1)
			return PubExpr{parse_expr(tokens)}
		case "match":
			tokens.consume(1)
			matchee := parse_ident(tokens)
			matchers := parse_match_block(tokens)
			var _default Block
			for i, matcher := range matchers {
				if reflect.TypeOf(matcher.matcher) == reflect.TypeOf(IdentExpr{}) {
					value := matcher.matcher.(IdentExpr).value
					if value == "_" {
						_default = matcher.block
						matchers = append(matchers[:i], matchers[i+1:]...)
						break
					}
				}
			}
			return MatchExpr{matchee, matchers, &_default}
		}
	}
	panic(fmt.Sprintln("Cannot convert token", tokens.peek(0)))
}

type MatchExpr struct {
	matchee       IdentExpr
	matchers      []Matcher
	default_match *Block
}

type Matcher struct {
	matcher Expr
	block   Block
}

func parse_match_block(tokens *Tokens) []Matcher {
	if tokens.peek(0).token_type == "colon" {
		if tokens.peek(1).token_type == "eol" {
			if tokens.peek(2).token_type == "indent" {
				indent_level := tokens.peek(2).value
				tokens.consume(3)

				var matchers []Matcher
				for {
					for tokens.peek(0).token_type == "eol" {
						tokens.consume(1)
					}
					if tokens.peek(0).token_type == "dedent" && tokens.peek(0).value == indent_level {
						tokens.consume(1)
						break
					}
					var matcher Matcher
					matcher.matcher = parse_expr(tokens)
					matcher.block = parse_block(tokens)
					matchers = append(matchers, matcher)
				}
				return matchers
			}
		} else {
			tokens.consume(1)
			var matcher Matcher
			matcher.matcher = parse_ident(tokens)
			matcher.block = parse_block(tokens)
			return []Matcher{matcher}
		}
	}
	panic("Block not found")
}

func print_struct(s any) {
	fmt.Printf("%#v\n", s)
}
func parse_block(tokens *Tokens) Block {
	if tokens.peek(0).token_type == "colon" {
		if tokens.peek(1).token_type == "eol" {
			if tokens.peek(2).token_type == "indent" {
				indent_level := tokens.peek(2).value
				tokens.consume(3)

				var block Block
				for {
					for tokens.peek(0).token_type == "eol" {
						tokens.consume(1)
					}
					if tokens.peek(0).token_type == "dedent" && tokens.peek(0).value == indent_level {
						tokens.consume(1)
						break
					}
					expr := parse_expr(tokens)
					block.body = append(block.body, expr)
				}
				return block
			}
		} else {
			tokens.consume(1)
			expr := parse_expr(tokens)
			return Block{[]Expr{expr}}
		}
	}
	panic("Block not found")
}

func parse_struct_block(tokens *Tokens) Block {
	if tokens.peek(0).token_type == "colon" {
		if tokens.peek(1).token_type == "eol" {
			if tokens.peek(2).token_type == "indent" {
				indent_level := tokens.peek(2).value
				tokens.consume(3)

				var block Block
				for {
					for tokens.peek(0).token_type == "eol" {
						tokens.consume(1)
					}
					if tokens.peek(0).token_type == "dedent" && tokens.peek(0).value == indent_level {
						tokens.consume(1)
						break
					}
					expr := parse_argument_decl(tokens)
					block.body = append(block.body, expr)
				}
				return block
			}
		} else {
			tokens.consume(1)
			expr := parse_expr(tokens)
			return Block{[]Expr{expr}}
		}
	}
	panic("Block not found")
}

func parse_int_lit(tokens *Tokens) IntExpr {
	num, err := strconv.Atoi(tokens.peek(0).value)
	if err != nil {
		panic("Cannot convert int literal to int")
	}
	expr := IntExpr{num}
	tokens.consume(1)
	return expr

}

func parse_ident(tokens *Tokens) IdentExpr {
	expr := IdentExpr{tokens.peek(0).value}
	tokens.consume(1)
	return expr

}

func parse_function_decl(tokens *Tokens) FunctionDeclExpr {
	return_type := tokens.peek(0).value
	func_decl := FunctionDeclExpr{return_type, []ArgumentDeclExpr{}, Block{}}
	tokens.consume(1)

	if tokens.peek(0).token_type == "ident" {
		tokens.consume(1)
	}
	if tokens.peek(0).token_type == "open_paren" {
		tokens.consume(1)
	} else {
		panic("uhoh")
	}
	for {
		for tokens.peek(0).token_type == "eol" {
			tokens.consume(1)
		}
		if tokens.peek(0).token_type == "comma" {
			tokens.consume(1)
		}
		if tokens.peek(0).token_type == "close_paren" {
			tokens.consume(1)
			break
		}
		func_decl.arguments = append(func_decl.arguments, parse_argument_decl(tokens))
	}
	func_decl.body = parse_block(tokens)
	return func_decl
}

type FuncCallArgumentExpr struct {
	name *string
	expr Expr
}

func parse_function_call_argument(tokens *Tokens) FuncCallArgumentExpr {
	if tokens.peek(0).token_type == "ident" && tokens.peek(1).token_type == "colon" {
		arg_name := tokens.peek(0).value
		tokens.consume(2)
		return FuncCallArgumentExpr{&arg_name, parse_expr(tokens)}
	} else {
		return FuncCallArgumentExpr{nil, parse_expr(tokens)}
	}
}

func parse_function_call(tokens *Tokens) FuncCallExpr {
	func_call := FuncCallExpr{IdentExpr{tokens.peek(0).value}, []Expr{}}
	paren_level := tokens.peek(1).value
	tokens.consume(2)
	for {
		for tokens.peek(0).token_type == "eol" {
			tokens.consume(1)
		}
		if tokens.peek(0).token_type == "comma" {
			tokens.consume(1)
		}
		if tokens.peek(0).token_type == "close_paren" && tokens.peek(0).value == paren_level {
			tokens.consume(1)
			break
		}
		func_call.arguments = append(func_call.arguments, parse_function_call_argument(tokens))
	}
	return func_call
}

func parse_math_atom(tokens *Tokens) (expr Expr, err error) {
	if tokens.peek(0) == nil {
		return nil, errors.New("unexpected EOF")
	}
	if tokens.peek(0).token_type == "open_paren" {
		tokens.consume(1)
		val, err := parse_math_expr(tokens, 0)
		if err != nil {
			return nil, err
		}
		if tokens.peek(0) == nil {
			return nil, errors.New("unexpected EOF")
		}
		if tokens.peek(0).token_type != "close_paren" {
			return nil, errors.New("unmatched '('")
		}
		tokens.consume(1)
		return val, nil
	} else if tokens.peek(0).token_type == "binop" {
		return nil, fmt.Errorf("expected an atom, not %v", tokens.peek(0))
	} else if tokens.peek(0).token_type == "int_lit" {
		return parse_int_lit(tokens), nil
	} else if tokens.peek(0).token_type == "ident" && tokens.peek(1).token_type == "open_paren" {
		return parse_function_call(tokens), nil
	} else if tokens.peek(0).token_type == "ident" {
		return parse_ident(tokens), nil
	}
	return nil, errors.New("not a known token")
}

func parse_math_expr(tokens *Tokens, min_prec int) (expr Expr, err error) {
	atom_lhs, err := parse_math_atom(tokens)
	if err != nil {
		return nil, err
	}
	precedences := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}
	for {
		curr := tokens.peek(0)
		if curr == nil || curr.token_type != "binop" || precedences[curr.value] < min_prec {
			break
		}
		if curr.token_type != "binop" {
			return nil, errors.New("current token isnt a binop")
		}
		op := curr.value
		precedence := precedences[curr.value]
		next_min_precedence := precedence + 1
		tokens.consume(1)
		atom_rhs, err := parse_math_expr(tokens, next_min_precedence)
		if err != nil {
			return nil, err
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
		}
	}
	return atom_lhs, nil
}
