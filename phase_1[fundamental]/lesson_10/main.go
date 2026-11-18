package main

import (
	"fmt"
)

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

type Car struct {
	Brand      string
	Model      string
	IsBrandNew bool
	Price      float64
	Stock      int
}

func readAndUpdate() {
	car := Car{
		Brand:      "Toyota",
		Model:      "CHR 2026",
		IsBrandNew: true,
		Price:      32000000.00,
		Stock:      10,
	}

	// read fields
	fmt.Println("Brand:", car.Brand)
	fmt.Println("Model:", car.Model)

	// update fields
	car.Price = 38000000.00
	car.Stock = 15

	fmt.Println("Updated Car:", car)
}

type Book struct {
	Name    string
	Author  string
	Release string
	Price   float64
	Stock   int
}

func printBook(book Book) {
	fmt.Println("Name:", book.Name)
	fmt.Println("Price:", book.Price)
	fmt.Println("Author:", book.Author)
	fmt.Println("Stock:", book.Stock)
}

func updateBook(book *Book, price float64) {
	book.Price = price
}

func structWithFunc() {
	book := Book{
		Name:    "Hadith",
		Price:   1000,
		Author:  "Mohammad Ali",
		Release: "01/01/1998",
		Stock:   20,
	}

	printBook(book)

	updateBook(&book, 800)

	fmt.Println("\nAfter updating book:")
	printBook(book)

}

type Customer struct {
	Name   string
	Email  string
	Status bool
	Bio    CusBio
}

type CusBio struct {
	Age     uint8
	DOB     string
	Country string
	About   string
}

func nestedStruct() {
	customer := Customer{
		Name:   "Jihad",
		Email:  "jihad@gmail.com",
		Status: true,
		Bio: CusBio{
			Age:     31,
			DOB:     "27-12-1997",
			Country: "Bangladesh",
			About:   "I am an Muslim",
		},
	}

	fmt.Println("Name:", customer.Name)
	fmt.Println("Bio: age-> ", customer.Bio.Age)
}

func main() {
	fmt.Println("<Lesson 10>")

	// defineStruct()
	// readAndUpdate()
	// structWithFunc()
	nestedStruct()
}
