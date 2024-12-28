package main

import (
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := router.GenerateRouter()
	fmt.Println("API running on port 3333")
	log.Fatal(http.ListenAndServe(":3333", r))
}