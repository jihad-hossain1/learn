package main

import (
	"fmt"
	"strconv"
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

func stringWithUnicode() {
	text := "বাংলা"

	letters := []rune(text)

	fmt.Println(len(text))
	fmt.Println(len(letters))

	fmt.Printf("%c\n", letters[0])
}

func iterateOverString() {
	text := "Go"

	for i := 0; i < len(text); i++ {
		fmt.Printf("%c\n", text[i])
	}

	for idx, value := range text {
		fmt.Printf("%d %c\n", idx+1, value)
	}

	textUnicode := "বাংলা"

	for idx, value := range textUnicode {
		fmt.Println(idx, string(value))
	}
}

func stringBuilder() {
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

func stringsBuiltInFunc() {
	name := "   Jihad Hossain   "

	fmt.Println(strings.TrimSpace(name))
	fmt.Println(strings.ToUpper(strings.TrimSpace(name)))
}

func splitAndJoin() {
	categories := "Electronics,Clothing,Food"

	result := strings.Split(categories, ",")

	for _, value := range result {
		fmt.Println(value)
	}

	fmt.Println(result)

	cats := []string{
		"a",
		"b",
		"c",
	}

	rslt := strings.Join(cats, " | ")

	fmt.Println(rslt)
}

func stringConversion() {
	age, err := strconv.Atoi("30")

	if err != nil {
		fmt.Println("Invalid number")
		return
	}

	fmt.Println(age)
	fmt.Printf("%T\n", age)

	birthYear := 1993

	result := strconv.Itoa(birthYear)

	fmt.Println(result)
	fmt.Printf("%T\n", result)

	price, err := strconv.ParseFloat("1020.653", 64)

	if err != nil {
		fmt.Println("Invalid price")
		return
	}

	fmt.Println(price)

	sprice := 1020.657

	sellPrice := strconv.FormatFloat(sprice, 'f', 2, 64)

	fmt.Println(sellPrice)

}

type Product struct {
	ID    int
	Name  string
	Price float64
}

func searchProducts(products []Product, keyword string) {
	keyword = strings.ToLower(strings.TrimSpace(keyword))

	for _, product := range products {
		name := strings.ToLower(product.Name)

		if strings.Contains(name, keyword) {
			fmt.Printf(
				"ID: %d | Name: %s | Price: %.2f\n",
				product.ID,
				product.Name,
				product.Price,
			)
		}
	}
}

func practicalProject() {
	products := []Product{
		{1, "laptop", 70000},
		{2, "Mouse", 4000},
		{3, "Keyboard", 2000},
		{4, "USB Driver", 1000},
		{5, "Monitor", 20000},
		{6, "Gaming Laptop", 300000},
		{7, "USB Cable", 200},
		{8, "HDMI Cable", 400},
		{9, "Floppy Driver", 300},
	}

	var search string

	fmt.Printf("Enter Your Query:")
	fmt.Scan(&search)

	searchProducts(products, search)
}

func practice() {
	var builder strings.Builder

	builder.WriteString("firstName")
	builder.WriteString(" ")
	builder.WriteString("lastName")

	fmt.Println(builder.Len())
	fmt.Printf("%s\n", builder.String())

	unicode := "বাংলাদেশ"

	convertRune := []rune(unicode)

	fmt.Println(len(unicode))
	fmt.Println(len(convertRune))
	fmt.Println(unicode)
}

func countCharacters(text string) int {
	makeItRune := []rune(text)
	fmt.Println("Bytes:", len(text))
	fmt.Println("Code Points:", len(makeItRune))

	return len(makeItRune)
}

func reverseString(text string) string {
	characters := []rune(text)

	// reverse using two pointers

	for left, right := 0, len(characters)-1; left < right; left, right = left+1, right-1 {
		characters[left], characters[right] = characters[right], characters[left]
	}

	return string(characters)
}

func main() {
	// stringInit()
	// stringDeclare()
	// runeString()
	// stringIndexing()
	// stringWithUnicode()
	// iterateOverString()
	// stringBuilder()
	// stringsBuiltInFunc()
	// splitAndJoin()
	// stringConversion()
	// practicalProject()
	// practice()
	// count := countCharacters("বাংলাদেশ,বাংলাদেশ,বাংলাদেশ,a")
	// fmt.Println("Total Char:", count)

	text := "1234567890987654321"

	reversed := reverseString(text)

	fmt.Println("Original:", text)
	fmt.Println("Reversed:", reversed)
}
