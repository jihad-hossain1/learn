package main

import "fmt"

func main (){
	fmt.Println("<Lesson 5>")

	// fmt.Println("if statement")
	// ifState()
	// loginValidation()
	// shortState()
	// switchState()
	// multiState()
	// switchWithExpression()
	// forLoop()
	// forWhile()
	// infiniteLoop()
	// breakAndContinue()
	// nestedLoop()
	// invGradSys()
	practiceEx()
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

func forLoop(){
	for i := 1; i <= 5; i++ {
		fmt.Println(i)
	}
}

func forWhile(){
	i := 1

	for i <= 5 {
		fmt.Println(i)
		i++
	}
}

func infiniteLoop(){
	for {
		fmt.Println("Running...")
		break
		// continue
	}
}

func breakAndContinue(){
	for i := 1; i <= 10; i++ {

		if i >= 5 {
			break
		}

		fmt.Println(i)
	}

	for i := 1; i <=5; i ++ {
		if i == 3 {
			continue
		}

		fmt.Println(i)
	}
}

func nestedLoop(){
	for i := 1; i <=3; i++ {
		for j := 1; j <= 2; j++ {
			fmt.Printf("i=%d, j=%d\n",i,j)
		}
	}
}

func invGradSys(){
	var total float64

	fmt.Print("Enter invoice total (TK): ")
	fmt.Scan(&total)

	switch {
	case total <= 0:
		fmt.Println("Invoice invoice amount")
	case total < 1000:
		fmt.Println("Small Invoice")
	case total < 10000: 
		fmt.Println("Medium Invoice")
	default:
		fmt.Println("Large Invoice")
	}

	for i := 1; i <= 3; i ++ {
		fmt.Printf("Invoice processing step %d\n",i)
	}
}

func practiceEx(){
	fmt.Println("Exercise 1: Take an integer and check whether it is positive, negative, or zero.")
	// ex1()
	
	fmt.Println(`Exercise 2: Take a user's age and categorize them as:
	- Child (0–12)
	- Teenager (13–19)
	- Adult (20–59)
	- Senior (60+)`,)
	ex2()

}

func ex1(){
	var value int 

	fmt.Print("Choose a value: ")
	fmt.Scan(&value)

	if value <= 0 {
		fmt.Println("negative")
	} else {
		fmt.Println("positive")
	}
}

func ex2(){
	var age int

	fmt.Print("Enter your age: ")
	fmt.Scan(&age)

	switch {
	case age < 0:
		fmt.Println("Not allow") 
	case age <= 12:
		fmt.Println("Child")
	case age >= 19:
		fmt.Println("Teenager")
	case age >= 59:
		fmt.Println("Adult")
	default:
		fmt.Println("Senior")
	}
}