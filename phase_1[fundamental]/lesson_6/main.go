package main

import "fmt"

func displayUser(name string, age int) {
	fmt.Println("Name:", name, "Age:", age)
}

func synShort(a, b int) {
	fmt.Println("Sum:", a+b)
}

func returnValue(a int, b int) int {
	result := a + b
	return result
}

func calculate(a, b int) (int, int, int, int) {
	sum := a + b
	sub := a - b
	mul := a * b
	div := a / b

	return sum, sub, mul, div
}

func namedReturnValue(a, b int) (sum int, difference int) {
	sum = a + b
	difference = a - b
	return
}

func variadicFunc(numbers ...int) int {
	total := 0

	for _, number := range numbers {
		total += number
	}

	return total
}

func functionAsValue(a, b int) int {
	return a + b
}

func functionArgAnFunc(a, b int, operation func(int, int) int) int {
	return operation(a, b)
}

func sum(a, b int) int {
	return a + b
}

func multi(a, b int) int {
	return a * b
}

func subtract(a, b int) int {
	return a - b
}

func division(a, b int) int {
	return a / b
}

func anonymousFunc() {
	great := func(name string) {
		fmt.Println("hello,", name)
	}

	great("Jihad")
}

func calculateSubtotal(price float64, quantity int) float64 {
	return price * float64(quantity)
}

func calculateDiscount(subtotal, discountRate float64) float64 {
	return subtotal * discountRate / 100
}

func calculateTotal(subtotal, discount float64) float64 {
	return subtotal - discount
}

func generateInvoice(
	productName string,
	price float64,
	quantity int,
	discountRate float64,
) {
	subtotal := calculateSubtotal(price, quantity)
	discount := calculateDiscount(subtotal, discountRate)
	total := calculateTotal(subtotal, discount)

	fmt.Println("-------INVOICE-------")
	fmt.Println("Product:", productName)
	fmt.Printf("Price: %.2f TK\n", price)
	fmt.Println("Quantity:", quantity)
	fmt.Printf("Subtotal: %.2f TK\n", subtotal)
	fmt.Printf("Discount: %.2f TK\n", discount)
	fmt.Printf("Total: %.2f TK\n", total)
}

func main() {
	fmt.Println("Function")

	// displayUser("jihad",31)

	// synShort(10,20)

	// total := returnValue(30,70)
	// fmt.Println("Total:",total)

	// sum,sub,mul,div := calculate(10,5)
	// fmt.Printf("\nAddition:%d, Subtraction:%d, Multiplication:%d, Division:%d\n",sum,sub,mul,div)

	// sum,_,_,_ := calculate(10,2)
	// fmt.Println("sum:",sum)

	// total,diff := namedReturnValue(20,10)
	// fmt.Println("Sum:",total,"diff:",diff)

	// fmt.Println(variadicFunc(10, 20))
	// fmt.Println(variadicFunc(10, 20, 30))
	// fmt.Println(variadicFunc(10, 20, 30, 40, 50))

	// operation := functionAsValue
	// result := operation(2, 3)
	// fmt.Println(result)

	// fmt.Println(functionArgAnFunc(10, 20, sum))
	// fmt.Println(functionArgAnFunc(10, 20, multi))
	// fmt.Println(functionArgAnFunc(10, 20, subtract))
	// fmt.Println(functionArgAnFunc(10, 20, division))

	// anonymousFunc()

	// func() {
	// 	fmt.Println("IIFE")
	// }()

	generateInvoice("Mouse", 1900, 3, 10)
}
