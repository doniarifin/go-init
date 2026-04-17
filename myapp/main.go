package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from myapp")
	})

	fmt.Println("Server running on :3000")
	http.ListenAndServe(":3000", nil)
}