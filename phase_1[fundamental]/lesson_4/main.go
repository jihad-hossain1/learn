package main

import (
	"fmt"
	"bufio"
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

	fmt.Println("<full line name space>")
	fullLine()
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
