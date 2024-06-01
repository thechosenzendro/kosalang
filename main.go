package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	//sourceCodePath := os.Args[1]
	sourceCodePath := "./main.kosa"
	outputPath := "./main.ll"
	data, err := os.ReadFile(sourceCodePath)
	if err != nil {
		panic(err)
	}
	sourceCode := strings.Replace(strings.Replace(string(data)+"\n", "\r\n", "\n", -1), "\t", "    ", -1)
	fmt.Println("Source code\n", sourceCode)

	tokens := tokenize(sourceCode)
	fmt.Printf("Tokens\n%v\n\n", tokens.tokens)

	ast := parse(tokens)
	fmt.Printf("AST\n%#v\n\n", ast)

	analyzedAST := analyze(ast)
	fmt.Printf("Analyzed AST\n%#v\n\n", analyzedAST)

	module := generate(analyzedAST)
	fmt.Printf("Generated code\n%v\n\n", module)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Write the module to the output file.
	if _, err := module.WriteTo(outputFile); err != nil {
		panic(err)
	}

}
