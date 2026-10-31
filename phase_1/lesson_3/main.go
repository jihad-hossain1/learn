package main

import "fmt"

func main() {
	fmt.Println("<Arithmetic Operators>")
	// Arithmetic Operators
	arithmeticOperators()

	fmt.Println("<Comparison Operators>")
	// Comparison Operators
	comparisonOperators()

	fmt.Println("<Logical Operators>")
	// Logical Operators
	logicalOperators()

	fmt.Println("<Assignment Operators>")
	// Assignment Operators
	assignmentOperators()

	fmt.Println("<Increment & Decrement>")
	// Increment & Decrement
	incDec()

	fmt.Println("<Precedence Rules>")
	// Precedence Rules
	precedenceRule()
}

func arithmeticOperators() {
	a := 20
	b := 6

	fmt.Println(a+b, "Addition")
	fmt.Println(a-b, "Subtraction")
	fmt.Println(a*b, "Multiplication")
	fmt.Println(a/b, "Division")
	fmt.Println(a%b, "Remainder")

	fmt.Println(10 / 3)
	fmt.Println(10.0 / 3.0)
}

func comparisonOperators() {
	a := 10
	b := 20

	fmt.Println(a == b)
	fmt.Println(a != b)
	fmt.Println(a > b)
	fmt.Println(a < b)
	fmt.Println(a >= b)
	fmt.Println(a <= b)
}

func logicalOperators() {
	age := 25
	hasID := true

	fmt.Println(age >= 18 && hasID)
	fmt.Println(age < 18 || hasID)
	fmt.Println(!hasID)
}

func assignmentOperators() {
	x := 10

	x += 5
	fmt.Println(x)

	x -= 3
	fmt.Println(x)

	x *= 2
	fmt.Println(x)

	x /= 4
	fmt.Println(x)

	x %= 4
	fmt.Println(x)

}

func incDec() {
	count := 10

	count++
	fmt.Println(count)

	count--
	fmt.Println(count)
}

func precedenceRule() {
	result := 10 + 5*2
	fmt.Println(result)

	result = (10 + 5) * 2
	fmt.Println(result)
}
