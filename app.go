package main

import (
	"demo/cache"
	"demo/http"
	"flag"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Cache service initializing ....")

	ttl := flag.Int("ttl", 10, "Time to live (sec)")
	capacity := flag.Int("capacity", 40, "LRU capacity")
	flag.Parse()

	tmp := "inmemory"
	typ := &tmp
	log.Println("Now cache ttl is :", *ttl)
	log.Println("Now cache LRU capacity is :", *capacity)
	c := cache.New(*typ, *ttl, *capacity)

	http.New(c).Listen()
}
