package main

import (
	"fmt"
	"net/http"
)

func fetchSvc(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
}

func main() {
	http.HandleFunc("/", fetchSvc)
	http.ListenAndServe(":8080", nil)
}
