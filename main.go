package main

import (
	"fmt"
	"net/http"
)

func handlder(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Neat Downloaded Running!")
}
func main() {
	http.HandleFunc("/", handlder)
	http.ListenAndServe(":8000", nil)
}
