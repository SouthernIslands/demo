package main

import (
	"demo/cache"
	"demo/http"
	"fmt"
)

func main() {
	fmt.Println("Hello Go")

	c := cache.New("inmemory")
	//go tcp.New(c).Listen()
	http.New(c).Listen()
}
