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

func accessingAndUpdate() {
	fruits := [3]string{"Apple", "Banana", "Orange"}

	fmt.Println(fruits[0])
	fmt.Println(fruits[1])

	fruits[1] = "Mango"

	fmt.Println(fruits)

	// array length
	length := len(fruits)
	fmt.Println(length)

}

func iterateOverArray() {
	numbers := [5]int{10, 20, 30, 40, 50}

	// traditional for loop
	for i := 0; i < len(numbers); i++ {
		fmt.Println(numbers[i])
	}

	// using range
	for index, value := range numbers {
		fmt.Println("Index:", index, "Value:", value)
	}

	// without index
	for _, value := range numbers {
		fmt.Println(value)
	}
}

func slices() {
	names := []string{"Jihad", "Rahim", "Karim"}

	for _, value := range names {
		fmt.Println(value)
	}

	// using make to create slice
	numbers := make([]int, 5)
	fmt.Println(numbers)
}

func sliceWithMake() {
	slice := make([]int, 3, 5)

	fmt.Println("Slice:", slice)
	fmt.Println("Length:", len(slice))
	fmt.Println("Capacity:", cap(slice))

	slice = append(slice, 10)
	slice = append(slice, 20)
	slice = append(slice, 30)

	for _, value := range slice {
		fmt.Println(value)
	}
}

func sliceOperation() {
	numbers := []int{10, 20, 30, 40, 50}

	fmt.Println(numbers[1:4])
	fmt.Println(numbers[1:3])
	fmt.Println(numbers[2:])
	fmt.Println(numbers[:])
}

func copySlice() {
	source := []int{10, 20, 30, 40, 50}

	destination := make([]int, len(source))

	copied := copy(destination, source)

	destination = append(destination, 60)

	fmt.Println("Source:", source)
	fmt.Println("Destination:", destination)
	fmt.Println("Copied elements:", copied)

	fmt.Println("=====")

	for index, value := range source {
		fmt.Printf("src: ind-%d, val-%d\n", index, value)
	}

	fmt.Println("=====")

	for index, value := range destination {
		fmt.Printf("dest: ind-%d, val-%d\n", index, value)
	}

}

func removeElementFromSlice() {
	numbers := []int{1, 2, 3, 4, 5, 6}

	index := 2

	numbers = append(numbers[:index], numbers[index+1:]...)

	fmt.Println(numbers)
}

func displayProducts(products []string) {
	fmt.Println("\n--- Product Inventory ---")

	for index, product := range products {
		fmt.Printf("%d. %s\n", index+1, product)
	}

	fmt.Println("Total products:", len(products))
}

func addProduct(products []string, product string) []string {
	products = append(products, product)
	return products
}

func practice() {
	products := []string{"Keyboard", "Mouse", "Monitor"}

	displayProducts(products)

	products = addProduct(products, "Laptop")

	displayProducts(products)
}

func main() {
	// arrayDeclare()
	// arrayExm()
	// arrayInitialization()
	// accessingAndUpdate()
	// iterateOverArray()
	// slices()
	// sliceWithMake()
	// sliceOperation()
	// copySlice()
	// removeElementFromSlice()
	practice()
}
