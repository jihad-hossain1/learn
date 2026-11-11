package main

import "fmt"

func decAndInit(){
	user := map[string]string{
		"name": "Jihad",
		"city": "Dhaka",
		"role": "Developer",
	}

	fmt.Println(user)
}

func main(){
	fmt.Println("hello!,Maps")
	decAndInit()
}






