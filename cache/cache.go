package cache

import "container/list"

type Cache interface {
	Set(string, []byte) error
	//Set(string, interface {}) error
	Get(string) ([]byte, error)
	//Get(string) (interface {}, error)
	Del(string) error
	GetStat() Stat
	GetMap() map[string]*entry
	GetList() list.List
}
