package cache

import (
	"container/heap"
	"container/list"
	"fmt"
	"sync"
	"time"
)

type entry struct {
	element list.Element
	key     string
	value   []byte
	expire  time.Time
	index   int
}

type inMemoryCache struct {
	table    map[string]*entry
	pq       priorityQueue
	lrulist  list.List
	mutex    sync.RWMutex
	ttl      time.Duration
	capacity int
	Stat
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	tmp, exist := c.table[k]
	//fmt.Printf("%s is exist %s",tmp,exist)
	//exist?
	if exist {
		c.removeEntry(tmp)
		//c.del(k, tmp.value)
	} else {
		if c.lrulist.Len() == c.capacity {
			e := c.leastUsedEntry()
			c.removeEntry(e)
		}
		tmp = &entry{}
		tmp.element = list.Element{}
		//tmp.index
	}

	tmp.key = k
	tmp.value = v
	tmp.expire = time.Now().Add(c.ttl)
	c.insertEntry(tmp)
	//c.add(k, v)
	return nil
}

func (c *inMemoryCache) Get(key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	tmp := c.table[key]
	if tmp == nil {
		fmt.Printf("key %s not found", key)
	}

	//if tmp.expire.Before(time.Now())
	c.touchEntry(c.table[key])
	return c.table[key].value, nil
	//tmp := c.c[k].created
	//c.c[k].created = tmp //cannot modify value of struct in map
	//fmt.Println(c.c[k].created)
	//return c.c[k].v, nil
}

func (c *inMemoryCache) Del(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tmp, exist := c.table[key]
	if exist {
		delete(c.table, key)
		c.removeEntry(tmp)
		//c.del(key, tmp.value)
	}
	return nil
}

func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

func (c *inMemoryCache) GetMap() map[string]*entry {
	return c.table
}

func (c *inMemoryCache) GetList() list.List {
	return c.lrulist
}

func (c *inMemoryCache) Init(capacity int, ttl time.Duration) {
	c.table = make(map[string]*entry, capacity)
	c.pq = make([]*entry, 0, capacity)
	c.lrulist.Init()
	heap.Init(&c.pq)
	c.ttl = ttl
	c.capacity = capacity
	c.mutex = sync.RWMutex{}
	c.Stat = Stat{}
	//entry1 := c.table["key1"]
	//entry1.key = "key1"
	//entry1.value = []byte{1,2,3}
	//entry1.expire = time.Now().Add(ttl)
	////c.lrulist.PushFront(&entry1.element)
	//c.table[entry1.key] = entry1

	//freelist?
}

func (c *inMemoryCache) removeEntry(e *entry) {
	if e.index != -1 {
		heap.Remove(&c.pq, e.index)
	}
	c.lrulist.Remove(&e.element)
	delete(c.table, e.key)
	fmt.Println("delete key :", e.key)
	e.key = "" //need?
	e.value = nil
}

func (c *inMemoryCache) insertEntry(e *entry) {
	if !e.expire.IsZero() {
		heap.Push(&c.pq, e)
	}
	c.lrulist.PushFront(&e.element)
	c.table[e.key] = e
}

func (c *inMemoryCache) touchEntry(e *entry) {
	c.lrulist.MoveToFront(&e.element)
}

func (c *inMemoryCache) leastUsedEntry() *entry {
	return c.lrulist.Back().Value.(*entry) // value?
}

func (c *inMemoryCache) expiredEntry(now time.Time) *entry {
	if len(c.pq) == 0 {
		return nil
	}

	if e := c.pq[0]; e.expire.Before(now) {
		return e
	}
	return nil
}

func newInMemoryCache(capacity, ttl int) *inMemoryCache {
	c := &inMemoryCache{}
	c.Init(capacity, time.Second*time.Duration(ttl))
	if ttl > 0 {
		go c.expire()
	}
	return c
}

//func (c *inMemoryCache) expireHelper(){
//	c.mutex.Lock()
//	defer c.mutex.Unlock()
//	e := c.expiredEntry(time.Now())
//	if e == nil {
//		return
//	}
//	c.removeEntry(e)
//}

func (c *inMemoryCache) expire() {
	//c.mutex.Lock()
	//defer c.mutex.Unlock()

	for {
		time.Sleep(time.Duration(5) * time.Second)
		i := 0
		for {
			c.mutex.Lock()
			e := c.expiredEntry(time.Now())
			c.mutex.Unlock()
			if e == nil {
				break
			}
			c.removeEntry(e)
			i += 1
		}
		fmt.Printf("%d entries expired", i)
	}
	//for {
	//	time.Sleep(c.ttl)
	//	c.mutex.RLock()
	//	//defer c.mutex.RUnlock()
	//	for{
	//
	//		c.mutex.RUnlock()
	//		if v.created.Add(c.ttl).Before(time.Now()) {
	//			fmt.Println("delete key :", k)
	//			c.Del(k)
	//		}
	//		c.mutex.RLock()
	//	}
	//
	//	c.mutex.RUnlock()
	//}
}
