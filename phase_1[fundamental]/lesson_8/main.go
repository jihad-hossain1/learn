package main

import "fmt"

func decAndInit() {
	user := map[string]string{
		"name": "Jihad",
		"city": "Dhaka",
		"role": "Developer",
	}

	fmt.Println(user)
}

func mapWithMake() {
	user := make(map[string]string)

	user["name"] = "Jihad"
	user["city"] = "Dhaka"
	user["role"] = "Dev"

	fmt.Println(user)
}

func decNilMap() {
	var user map[string]string

	fmt.Println(user == nil)

	// user["name"] = "jihad"

	// fmt.Println(user)
}

func addElementsOnMap() {
	products := make(map[string]int)

	products["Laptop"] = 75000
	products["Mouse"] = 1200
	products["Keyboard"] = 2500

	fmt.Println(products)

	// access value
	fmt.Println("Laptop:", products["Laptop"])
}

func checkKey() {
	prices := map[string]int{
		"Laptop": 75000,
		"Mouse":  1200,
	}

	price, ok := prices["Laptop"]
	if ok {
		fmt.Println("Price:", price)
	} else {
		fmt.Println("Product not found")
	}

	// key remove from map
	delete(prices, "Mouse")
	fmt.Println(prices)
}

func iterateOverMap() {
	user := map[string]string{
		"name": "jihad",
		"city": "dhaka",
		"role": "dev",
	}

	for key, value := range user {
		fmt.Println(key, ":", value)
	}
}

func main() {
	fmt.Println("hello!,Maps")
	// decAndInit()
	// mapWithMake()
	// decNilMap()
	// addElementsOnMap()
	// checkKey()
	iterateOverMap()
}
