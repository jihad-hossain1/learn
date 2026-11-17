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

func valueChange(value *int) {
	*value = 200
}

func usePointer() {
	age := 20

	var ptr *int

	ptr = &age

	valueChange(ptr)
	fmt.Println(age)
}

func swap(a *int, b *int) {
	temp := *a
	*a = *b
	*b = temp
}

func swapNumber(){
	a := 10
	b := 20

	swap(&a, &b)

	fmt.Println("a:", a, "b:", b)
}

type  Prod struct {
	Name string 
	Price float64
	Stock int 
}

func updateProduct(product *Prod,price float64, stock int){
	product.Price = price
	product.Stock += stock
}

func productPricing(){
	product := Prod{
		Name: "abc",
		Stock: 20,
		Price: 100,
	}
	
	fmt.Println("Before update product:",product)

	updateProduct(&product,200,-5)

	fmt.Println("After update product:",product)
}

func nilPointerHandle(){
	var ptr *int

	fmt.Println("pointer",ptr)

	if ptr == nil {
		fmt.Println("pointer is nil")
		return
	}

	age := 20

	ptr = &age 

	fmt.Println("pointer",ptr)
}

func practice() {
	// usePointer()
	// swapNumber()
	// productPricing()
	nilPointerHandle()

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
	// practical()
	practice()
}
