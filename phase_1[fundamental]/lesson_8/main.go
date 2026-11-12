package main

import (
	"fmt"
	"slices"
)

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

func sortingMapMethod() {
	prices := map[string]int{
		"laptop":  3000,
		"mouse":   300,
		"desktop": 4000,
	}

	keys := make([]string, 0, len(prices))

	for key := range prices {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	for _, key := range keys {
		fmt.Println(key, prices[key])
	}

	fmt.Println("Total products:", len(keys))
}

func nestedMaps() {
	users := map[string]map[string]string{
		"u1": {
			"name": "jihad",
			"city": "dhaka",
		},
		"u2": {
			"name": "rahim",
			"city": "chittagong",
		},
	}

	fmt.Println(users["u1"]["name"])
	fmt.Println(users["u2"]["city"])

	prices := make(map[string]map[string]string)

	prices["p1"] = make(map[string]string)
	prices["p1"]["name"] = "laptop"
	prices["p1"]["price"] = "100"

	for key, value := range prices {
		fmt.Println(key, value)
	}
}

func differentValueWithMaps() {
	user := map[string]any{
		"name":     "jihad",
		"age":      31,
		"salary":   500000.00,
		"isActive": true,
	}

	fmt.Println(user)

	for key, value := range user {
		fmt.Println(key, value)
	}
}

func addProd(products map[string]float64, name string, price float64) {
	products[name] = price
	fmt.Println("Prod add sucs!")
}

func displayProd(products map[string]float64) {
	fmt.Println("\n---- Product Inventory ----\n")

	for name, price := range products {
		fmt.Printf("%s: %.2f TK\n", name, price)
	}

	fmt.Println("Total Products:", len(products))
}

func searchProd(products map[string]float64, name string) {
	price, ok := products[name]

	if !ok {
		fmt.Println("Prod not found!")
		return
	}

	fmt.Printf("Product: %s\nPrice: %.2f TK\n", name, price)
}

func updProd(products map[string]float64, name string, price float64) {
	_, ok := products[name]

	if !ok {
		fmt.Println("Product not found")
		return
	}

	products[name] = price
	fmt.Println("Product updated sucs")
}

func delProd(products map[string]float64, name string) {
	_, ok := products[name]

	if !ok {
		fmt.Println("Product not found!")
		return
	}

	delete(products, name)
	fmt.Println("Product deleted sucs")
}

func inventoryManage() {
	products := map[string]float64{
		"laptop":   75000,
		"mouse":    2300,
		"keyboard": 12000,
	}

	displayProd(products)
	addProd(products, "monitor", 20000)
	searchProd(products, "mouse")
	updProd(products, "laptop", 30000)
	delProd(products, "keyboard")
	displayProd(products)
}



func ex1(){
	fmt.Println("Create a map containing your:(name,age,city,profession) Print every key-value pair.")

	profile := map[string]any{
		"name": "Jihad Hossain",
		"age": 31,
		"city": "Dhaka",
		"profession": "Dev",
	}

	for key,value := range profile{
		fmt.Printf("%s : %s\n",key,value)
	}


}

func addStd(students map[string]map[string]any,name string,mark uint8){
	students[name] = map[string]any{
		"name": name,
		"mark": mark,
	}	
}


func calStdMrk(students map[string]map[string]any) float64{
	var total uint

	for _,student := range students {
		total += uint(student["mark"].(uint8))
	}

	if len(students) == 0 {
		return 0
	}

	return  float64(total) / float64((len(students)))
}

func displayStd(students map[string]map[string]any){
	for _,student := range students{
		fmt.Printf("name: %s, mark: %s",student["name"],student["mark"])
	}
}

func ex2(){
	fmt.Println("Create a map of student names and their marks. [1] Add five students. [2] Print all students and marks. [3] Calculate the average marks.")
	
	students := make(map[string]map[string]any)

	addStd(students,"abc",89)
	addStd(students,"efg",43)
	addStd(students,"def",80)
	addStd(students,"etc",55)
	addStd(students,"xyz",90)

	fmt.Println("Student Records:")
	displayStd(students)

	average := calStdMrk(students)

	fmt.Printf("\nAverage Marge: %.2f\n",average)
	
}

func practiceEx(){
	// ex1()
	ex2()
}

func main() {
	fmt.Println("hello!,Maps")
	// decAndInit()
	// mapWithMake()
	// decNilMap()
	// addElementsOnMap()
	// checkKey()
	// iterateOverMap()
	// sortingMapMethod()
	// nestedMaps()
	// differentValueWithMaps()
	// inventoryManage()
	practiceEx()
}
