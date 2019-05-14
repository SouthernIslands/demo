package cache

import (
	"fmt"
	"sync"
	"time"
)

type value struct {
	v       []byte
	created time.Time
}

type inMemoryCache struct {
	c     map[string]value
	mutex sync.RWMutex
	Stat
	ttl time.Duration
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	tmp, exist := c.c[k]
	//exist?
	if exist {
		c.del(k, tmp.v)
	}
	c.c[k] = value{v, time.Now()}
	c.add(k, v)
	return nil
}

func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	tmp := c.c[k].created
	c.c[k].created = tmp //cannot modify value of struct in map
	fmt.Println(c.c[k].created)
	return c.c[k].v, nil
}

func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.c[k]
	if exist {
		delete(c.c, k)
		c.del(k, v.v)
	}
	return nil
}

func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

func (c *inMemoryCache) GetMap() map[string]value {
	return c.c
}

func newInMemoryCache(ttl int) *inMemoryCache {
	c := &inMemoryCache{make(map[string]value), sync.RWMutex{}, Stat{}, time.Duration(ttl) * time.Second}
	if ttl > 0 {
		go c.expire()
	}
	return c
}

func (c *inMemoryCache) expire() {
	for {
		time.Sleep(c.ttl)
		c.mutex.RLock()
		//defer c.mutex.RUnlock()
		for k, v := range c.c {
			c.mutex.RUnlock()
			if v.created.Add(c.ttl).Before(time.Now()) {
				fmt.Println("delete key :", k)
				c.Del(k)
			}
			c.mutex.RLock()
		}

		c.mutex.RUnlock()
	}
}
