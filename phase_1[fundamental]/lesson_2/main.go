package main

import "fmt"

func main() {
	var name string = "jihad"
	var age int = 30
	var price float64 = 99.99
	var isActive bool = true

	fmt.Println(name)
	fmt.Println(age)
	fmt.Println(price)
	fmt.Println(isActive)

	// Short variable declaration pattern
	prodName := "Apple"
	createdAt := "12-02-12"
	status := true
	prodPrice := 19.99

	fmt.Println(prodName, prodPrice, createdAt, status)

	// constant declaration
	const appName = "olo invoice"
	const taxRate = 15.0

	fmt.Println(appName)
	fmt.Println(taxRate)

	// zero values declare
	var stdName string
	var stdAge int
	var mark float64
	var stdStatus bool

	fmt.Println(stdName)
	fmt.Println(stdAge)
	fmt.Println(mark)
	fmt.Println(stdStatus)

	// practical example
	devPro()
}

func devPro() {
	name := "Jihad"
	profession := "Full Stack Developer"
	experience := 3
	salary := 50000.0
	isEmployed := true

	fmt.Println("Name:", name)
	fmt.Println("Profession:", profession)
	fmt.Println("Experience:", experience, "years")
	fmt.Println("Salary:", salary, "TK.")
	fmt.Println("Employed:", isEmployed)
}
