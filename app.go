package main

import (
	"demo/cache"
	"demo/http"
	"flag"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello Go")

	ttl := flag.Int("ttl", 5, "Time to live")
	flag.Parse()

	tmp := "inmemory"
	typ := &tmp
	log.Println("Now cache ttl is :", *ttl)
	c := cache.New(*typ, *ttl)
	//c := cache.New("inmemory")
	//go tcp.New(c).Listen()
	http.New(c).Listen()
}
