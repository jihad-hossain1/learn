package main

import (
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w,"hello world")})

	fmt.Println("server is running on :8000")
	http.ListenAndServe(":8000",nil)
}
