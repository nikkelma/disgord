// lru (least recently used) will overwrite cached items that have been used the least when the cache limit is reached.
package lru

import (
	"sync"

	"github.com/andersfylling/disgord/cache/interfaces"
	"github.com/andersfylling/snowflake/v2"
)

type Snowflake = snowflake.Snowflake

func NewCacheItem(content interface{}) *CacheItem {
	return &CacheItem{
		item: content,
	}
}

type CacheItem struct {
	item interface{}

	// allows for least recently used monitoring
	used uint64
}

func (i *CacheItem) Object() interface{} {
	return i.item
}

func (i *CacheItem) Set(v interface{}) {
	i.item = v
}

func (i *CacheItem) increment() {
	i.used++
}

func NewCacheList(size uint) *CacheList {
	return &CacheList{
		items: make(map[Snowflake]*CacheItem, size),
		limit: size,
	}
}

type CacheList struct {
	sync.RWMutex
	items map[Snowflake]*CacheItem
	limit uint // 0 == unlimited

	misses uint64 // opposite of cache hits
	hits   uint64
}

func (list *CacheList) size() uint {
	return uint(len(list.items))
}

func (list *CacheList) First() (item *CacheItem, key Snowflake) {
	for key, item = range list.items {
		return
	}

	return
}

// set adds a new item to the list or returns false if the item already exists
func (list *CacheList) Set(id Snowflake, newItemI interfaces.CacheableItem) {
	newItem := newItemI.(*CacheItem)
	if newItem.used == 0 {
		newItem.used++
	}
	if item, exists := list.items[id]; exists { // check if it points to a diff item
		if item.item != newItem.item || item != newItem {
			*item = *newItem
		}
		return
	} else {
		list.items[id] = newItem
	}

	if list.limit == 0 || list.size() <= list.limit {
		return
	}
	// if limit is reached, replace the content of the least recently used (lru)
	list.removeLRU(id)
}

func (list *CacheList) UpdateLifetime(item interfaces.CacheableItem) {}

func (list *CacheList) removeLRU(exception Snowflake) {
	lru, lruKey := list.First()
	for key, item := range list.items {
		if key != exception && item.used < lru.used {
			// TODO: create an lru map, for later?
			lru = item
			lruKey = key
		}
	}

	delete(list.items, lruKey)
}

// get an item from the list.
func (list *CacheList) Get(id Snowflake) (ret interfaces.CacheableItem, exists bool) {
	var item *CacheItem
	if item, exists = list.items[id]; exists {
		ret = item
		item.used++
		list.hits++
	} else {
		list.misses++
	}
	return
}

func (list *CacheList) CreateCacheableItem(content interface{}) interfaces.CacheableItem {
	return NewCacheItem(content)
}

func (list *CacheList) Efficiency() float64 {
	return float64(list.hits) / float64(list.misses+list.hits)
}

var _ interfaces.CacheAlger = (*CacheList)(nil)
var _ interfaces.CacheableItem = (*CacheItem)(nil)
