package main

import (
	"fmt"

	"flag"
)

func main() {
	var port string
	var verbose bool
	var valor int
	flag.StringVar(&port, "port", "7000", "Server port")
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.IntVar(&valor, "valor", 0, "some value")

	flag.Parse()

	if verbose {
		fmt.Println("Server is running on port", port)
		fmt.Println("Valor", valor)
	} else {
		fmt.Println(port)
	}
}
