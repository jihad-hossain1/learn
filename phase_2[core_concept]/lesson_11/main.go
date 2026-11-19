package main

import (
	"fmt"
	"strings"
)

func stringInit() {
	name := "Jihad"
	message := "Hello, Go Programming"

	fmt.Println(name)
	fmt.Println(message)

	name = "A" + name[1:]

	fmt.Println(name)
}

func stringDeclare() {
	name := "abc"

	var city string = "Dhaka"

	var country string

	const language = "Go"

	fmt.Println(name)
	fmt.Println(city)
	fmt.Println(country)
	fmt.Println(language)

	message := "Hello\nWorld"

	fmt.Println(message)

	message = `Hello\nWorld`

	fmt.Println(message)

	fmt.Println(len(message))

	text := "বাংলা"
	fmt.Println(len(text))
}

func runeString() {
	var letter rune = 'A'
	var bangla rune = 'ক'

	fmt.Println(letter)
	fmt.Println(bangla)

	fmt.Printf("%c\n", letter)
	fmt.Printf("%c\n", bangla)

	fmt.Printf("%T\n", bangla)
}

func stringIndexing() {
	text := "Go"

	fmt.Println(text[0])
	fmt.Println(text[1])

	fmt.Printf("%c\n", text[0])
	fmt.Printf("%c\n", text[1])

	unicodeText := "বাংলা"

	fmt.Println(unicodeText[0])
}

func stringWithUnicode(){
	text := "বাংলা"

	letters := []rune(text)

	fmt.Println(len(text))
	fmt.Println(len(letters))

	fmt.Printf("%c\n",letters[0])
}

func iterateOverString(){
	text := "Go"

	for i := 0; i < len(text); i++ {
		fmt.Printf("%c\n",text[i])
	}

	for idx,value := range text{
		fmt.Printf("%d %c\n", idx+1, value)
	} 

	textUnicode := "বাংলা"

	for idx,value := range textUnicode {
		fmt.Println(idx,string(value))
	}
}

func stringBuilder(){
	firstName := "Jihad"
	lastName := "Hossain"

	fullName := firstName + " " + lastName

	fmt.Println(fullName)

	var builder strings.Builder

	builder.WriteString("Hello")
	builder.WriteString(" ")
	builder.WriteString("Go")

	result := builder.String()

	fmt.Println(result)
}

func stringsBuiltInFunc(){
	name := "   Jihad Hossain   "

	fmt.Println(strings.TrimSpace(name))
	fmt.Println(strings.ToUpper(strings.TrimSpace(name)))
}





func main() {
	// stringInit()
	// stringDeclare()
	// runeString()
	// stringIndexing()
	// stringWithUnicode()
	// iterateOverString()
	// stringBuilder()
	stringsBuiltInFunc()
}
