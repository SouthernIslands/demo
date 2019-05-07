package main

import (
	"demo/cache"
	"demo/http"
	"fmt"
)

func main() {
	fmt.Println("Hello Go")

	c := cache.New("inmemory")
	http.New(c).Listen()
}
