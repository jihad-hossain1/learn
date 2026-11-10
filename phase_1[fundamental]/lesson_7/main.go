package main

import "fmt"

func arrayDeclare() {
	var numbers [5]int

	numbers[0] = 10
	numbers[1] = 20
	numbers[2] = 30
	numbers[3] = 40
	numbers[4] = 50

	fmt.Println(numbers)
}

func arrayExm() {
	var ages [5]int
	var prices [10]float64
	var names [3]string
	var active [4]bool

	ages[0] = 1
	ages[1] = 2
	ages[2] = 3
	ages[3] = 4
	ages[4] = 5

	fmt.Println(ages)
	fmt.Println(prices)
	fmt.Println(names)
	fmt.Println(active)
}

func arrayInitialization() {
	// numbers := [3]int{10, 20, 30}
	numbers := [...]int{10, 20, 30}

	fmt.Println(numbers)
}

func main() {
	// arrayDeclare()
	// arrayExm()
	arrayInitialization()
}
