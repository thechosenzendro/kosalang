package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("./examples/main.kosa")
	if err != nil {
		panic(err)
	}
	fmt.Println("Source code")
	fmt.Println(string(data))

	tokens := tokenize(string(data))
	fmt.Println("Tokens")
	fmt.Println(tokens)

	ast := parse(tokens)
	fmt.Println("AST")
	fmt.Printf("%#v\n", ast)
}

type Token struct {
	token_type string
	value      string
}
