package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func main() {
	fmt.Println("Servidor rodando na porta 5000")

	// http.Handle("/hello", hello)

	http.ListenAndServe(":5000", http.HandlerFunc(helloHandler))
}
