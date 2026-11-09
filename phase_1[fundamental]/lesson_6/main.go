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

	fmt.Println(variadicFunc(10, 20))
	fmt.Println(variadicFunc(10, 20, 30))
	fmt.Println(variadicFunc(10, 20, 30, 40, 50))

}
