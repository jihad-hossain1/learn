package main

import "fmt"

func memoryAdd() {
	var age int = 30

	fmt.Println("age memAd", &age)
}

func decPointer() {
	var age int = 32

	var ptr *int

	ptr = &age

	fmt.Println(ptr)

	price := 200

	point := &price

	fmt.Println("Price:", price)
	fmt.Println("Pointer:", point)
	fmt.Println("value through pointer:", *point)
}

func modifyValueAPointer() {
	price := 500

	ptr := &price

	*ptr = 750

	fmt.Println("Update price:", price)
}

func updatePrice(price float64) {
	fmt.Println("inside up price before up", price)
	price = 1000
	fmt.Println("inside up price after up", price)
}

func upPriceResult() {
	price := 500.0

	updatePrice(price)

	fmt.Println(price)
}

func upPrWithPoint(price *float64) {
	*price = 1000
}

func upPtrResult() {
	price := 500.0

	upPrWithPoint(&price)

	fmt.Println(price)
}

func pointerWithFunc() {
	// upPriceResult()
	upPtrResult()
}

func nilPointer() {
	var ptr *int

	fmt.Println(ptr)

	if ptr == nil {
		fmt.Println("Pointer is nill")
		return
	}

	fmt.Println(*ptr)
}

func makePointerWithNew() {
	ptr := new(int)

	fmt.Println(*ptr)

	*ptr = 100

	fmt.Println(*ptr)
}

type Product struct {
	Name  string
	Price float64
}

func upPrd(product *Product) {
	product.Price = 1500
}

func pointerWithStruct() {
	product := Product{
		Name:  "Laptop",
		Price: 1200,
	}

	upPrd(&product)

	fmt.Println(product)
}

type Pro struct {
	Name  string
	Price float64
	Stock int
}

func upPrdInv(product *Pro, newPrice float64, quantity int) {
	product.Price = newPrice
	product.Stock += quantity
}

func practical() {
	product := Pro{
		Name:  "Abc",
		Price: 200,
		Stock: 20,
	}

	upPrdInv(&product, 300, 13)

	fmt.Println("Product:", product.Name)
	fmt.Println("Price:", product.Price)
	fmt.Println("Stock:", product.Stock)
}

func main() {
	fmt.Println("Hello, lesson 9")
	// memoryAdd()
	// decPointer()
	// modifyValueAPointer()
	// pointerWithFunc()
	// nilPointer()
	// makePointerWithNew()
	// pointerWithStruct()
	practical()
}
