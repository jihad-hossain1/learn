package main

import "fmt"

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

func main() {
	// stringInit()
	// stringDeclare()
	// runeString()
	stringIndexing()
}
