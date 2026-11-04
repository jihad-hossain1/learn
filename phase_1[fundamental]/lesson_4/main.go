package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main(){
	// fmt.Println("<lesson 4>")
	// startLesson()

	// fmt.Println("<Format specifier>")
	// formatSpecifier()

	// fmt.Println("<User Input>")
	// userInput()

	// fmt.Println("<user multiple input>")
	// userMultiInput()

	// fmt.Println("<full line name space>")
	// fullLine()

	// fmt.Println("<Simple calculator>")
	// simpleCalculator()

	fmt.Println("<Practice>")
	practice()
}

func startLesson(){
	fmt.Print("hello ")
	fmt.Print("go!\n")

	name := "jihad"
	age := 31

	fmt.Printf("My name is %s and i am %d years old.\n\n",name,age)
}

func formatSpecifier(){
	name := "Jihad"
	age := 31
	salary := 45000.567
	isDev := true 

	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Age: %d\n",age)
	fmt.Printf("Salary: %.2f TK\n", salary)
	fmt.Printf("Developer: %t\n", isDev)

	fmt.Printf("Type of name: %T\n",name)
}

func userInput(){
	var name string 

	fmt.Print("Enter your name: ")
	fmt.Scan(&name)

	fmt.Println("Hello,",name)
}

func userMultiInput(){
	var name string
	var age int 
	var salary float64

	fmt.Print("Enter your name: ")
	fmt.Scan(&name)

	fmt.Print("Enter your age: ")
	fmt.Scan(&age)

	fmt.Print("Enter your salary: ")
	fmt.Scan(&salary)

	fmt.Println("\n---- User Info ----")
	fmt.Printf("Name: %s\n",name)
	fmt.Printf("Age: %d\n", age)
	fmt.Printf("Salary: %.2f TK\n", salary)
}

func fullLine(){
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your full name: ")

	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	name = strings.TrimSpace(name)

	fmt.Println("Hello,", name)
}

func simpleCalculator(){
	var a,b float64
	var operator string 

	fmt.Print("Enter first number: ")
	fmt.Scan(&a)

	fmt.Print("Enter operator (=,-,*,/): ")
	fmt.Scan(&operator)

	fmt.Print("Enter second number: ")
	fmt.Scan(&b)

	switch operator {
	case "+":
		fmt.Printf("Result: %.2f\n",a+b)
	case "-":
		fmt.Printf("Result: %.2f",a-b)
	case "*":
		fmt.Printf("Result: %.2f\n",a*b)
	case "/":
		if b == 0 {
			fmt.Println("Error: Division by zero!")
		} else {
		fmt.Printf("Result: %.2f\n",a/b)
		}

	default: 
		fmt.Println("Invalid operator!")
	}
}

func practice(){
	// fmt.Println("[1]. Take a user's name and age, then print a formatted greeting.")
	// p1()

	// fmt.Println("[2]. Take two numbers and display their sum.")
	// p2()

	fmt.Println("[3]. Take a product name, price, and quantity, then calculate the total price.")
	p3()
}

func p1(){
	var name string
	var age uint16 

	fmt.Print("Enter Your name: ")
	fmt.Scan(&name)

	fmt.Print("Enter Your Age: ")
	fmt.Scan(&age)

	fmt.Printf("Hello %s & age %d, thank you are joining with us!\n", name,age)
}

func p2(){
	var a,b uint 

	fmt.Print("Enter First Number: ")
	fmt.Scan(&a)

	if a == 0 {
		fmt.Println("Error: You can't add zero!")
		return
	}

	fmt.Print("Enter Second Number: ")
	fmt.Scan(&b)

	if b == 0 {
		fmt.Println("Error: You can't add zero!")
		return
	}

	fmt.Println("Result:", a+b)
}

func p3(){
	reader := bufio.NewReader(os.Stdin)

	var price float64
	var qty uint

	fmt.Print("Enter product name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	name = strings.TrimSpace(name)

	fmt.Print("Enter product price: ")
	fmt.Scan(&price)

	fmt.Print("Enter product quantity: ")
	fmt.Scan(&qty)

	total := price * float64(qty)

	fmt.Print("\n")

	fmt.Printf("Product name '%s', price/per %.2f, qty %d\n",name,price,qty)
	fmt.Println("Total Price:",total)


}