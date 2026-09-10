package main

import (
	"fmt"
	"net/http"

	"github.com/vanshsharma3777/WorkerBox/internal/api"
)

func main() {
	fmt.Println("Hit main.go ")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", api.Test)
	http.ListenAndServe(":8080", mux)

}
