package main

import (
	"fmt"

	"example.com/go-lang/calculator"
)

func main() {
	fmt.Println("<lesson 12>")

	result := calculator.Add(10, 20)
	fmt.Println("Result:", result)
}
