package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home Handler")
}

func olaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ola Handler")
}

func main() {
	fmt.Println("Servidor rodando na porta 5000")
	mux := http.NewServeMux()
	mux.HandleFunc("/ola/", olaHandler)
	mux.HandleFunc("/ola/mundo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Olá, mundo Handler")
	})
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World Handler")
	})
	mux.HandleFunc("/", homeHandler)

	mux.HandleFunc("ola.pessoas.com/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Este é o meu site de olá")
	})

	http.ListenAndServe(":5000", mux)
}
