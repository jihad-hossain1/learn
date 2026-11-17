package main

import "fmt"

type Product struct {
	Name  string
	Price float64
	Stock int
}

func defineStruct() {
	product := Product{
		Name:  "Abc",
		Price: 100,
		Stock: 10,
	}

	product2 := Product{
		"Abc 2",
		100,
		20,
	}

	fmt.Println(product)
	fmt.Println(product2)

	var newProduct Product

	newProduct.Name = "Monitor"
	newProduct.Price = 15000
	newProduct.Stock = 12

	fmt.Println("declare first, assign later", newProduct)

	usingPointerProduct := &Product{
		Name:  "Laptop",
		Price: 45000,
		Stock: 20,
	}

	fmt.Println("using a pointer", usingPointerProduct)
}

func main() {
	fmt.Println("<Lesson 10>")

	defineStruct()
}
