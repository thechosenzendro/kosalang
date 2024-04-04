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
	fmt.Println(string(data))
	
	tokens := tokenize(string(data))
	fmt.Println(tokens)

	ast := parse(tokens)
	fmt.Println(ast)
}

type Token struct {
	token_type string
	value      string
}