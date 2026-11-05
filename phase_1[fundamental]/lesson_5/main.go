package main

import "fmt"

func main (){
	fmt.Println("<Lesson 5>")

	fmt.Println("if statement")
	// ifState()
	// loginValidation()
	// shortState()
	// switchState()
	// multiState()
	switchWithExpression()
}	

func ifState(){
	age := 12

	if age >= 18 {
		fmt.Println("you are adult")
	} else {
		fmt.Println("not eligible to vote.")
	}

	marks := 60 

	if marks >= 80 {
		fmt.Println("Grade: A+")
	} else if marks >= 70 {
		fmt.Println("Grade: A")
	} else if marks >= 60 {
		fmt.Println("Grade: A-")
	} else if marks >= 50 {
		fmt.Println("Grade: B")
	} else {
		fmt.Println("Grade: Fail")
	}
}

func loginValidation(){
	username := "admin"
	password := "1234565"

	if username == "admin" && password == "123456" {
		fmt.Println("Login success")
	} else {
		fmt.Println("Invalid credentials")
	}
}

func shortState(){
	if age := 30; age >= 18 {
		fmt.Println("Adult")
	} else {
		fmt.Println("Minor")
	}
}

func switchState(){
	day := 6

	switch day {
	case 1:
		fmt.Println("Sun")
	case 2:
		fmt.Println("Mon")
	case 3:
		fmt.Println("Tue")
	case 4:
		fmt.Println("Wed")
	default:
		fmt.Println("Invalid Day")
	}
}

func multiState(){
	day := "Sat"

	switch day {
	case "Fri", "Sat":
		fmt.Println("Weekend!")
	default:
		fmt.Println("Working Day")
	}
}

func switchWithExpression(){
	age := 75 

	switch {
	case age < 13:
		fmt.Println("Child")
	case age < 20:
		fmt.Println("Teenager")
	case age < 60:
		fmt.Println("Adult")
	default:
		fmt.Println("Senior")
	}
}