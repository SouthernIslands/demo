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

	ttl := flag.Int("ttl", 60, "Time to live")
	capacity := flag.Int("capacity", 3, "LRU capacity")
	flag.Parse()

	tmp := "inmemory"
	typ := &tmp
	log.Println("Now cache ttl is :", *ttl)
	log.Println("Now cache LRU capacity is :", *capacity)
	c := cache.New(*typ, *ttl, *capacity)
	//c := cache.New("inmemory")
	//go tcp.New(c).Listen()
	http.New(c).Listen()
}
