package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TanishkBansode/cli-chat/router"
)

func main() {
	r := router.Router()
	fmt.Println("Server is getting live...")
	log.Fatal(http.ListenAndServe(":4000", r))
}
