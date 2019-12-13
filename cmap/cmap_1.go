package cmap

import "sync"

var SHARD_COUNT = 10

type ConcurrentMap []*ConcurrentMapShared

type ConcurrentMapShared struct {
	items map[string]interface{}
	sync.RWMutex
}

func New() ConcurrentMap {
	m := make(ConcurrentMap, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i ++{
		m[i] = &ConcurrentMapShared{
			items: make(map[string]interface{}),
		}
	}
	return m
}

func (m ConcurrentMap) GetShard(key string) *ConcurrentMapShared {
	return m[uint(fnv32(key))%uint(SHARD_COUNT)]
}

// FNV hash
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (m ConcurrentMap) Set(key string, value interface{}) {
	shard := m.GetShard(key) // 段定位找到分片
	shard.Lock()       // 分片上锁
	shard.items[key] = value // 分片操作
	shard.Unlock()       // 分片解锁
}

func (m ConcurrentMap) Get(key string) (interface{}, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (m ConcurrentMap) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++{
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (m ConcurrentMap) Keys() []string  {
	count := m.Count()
	ch := make(chan  string, count)

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(SHARD_COUNT)
		for _, shard := range m {
			// 每一个分片启动一个协程, 遍历key
			go func(shard *ConcurrentMapShared) {
				defer wg.Done()
				shard.RLock()
				// 每个分片中的key遍历后都写入统计用的channel
				for key := range shard.items {
					ch <- key
				}
				shard.RUnlock()
			}(shard)
		}
		wg.Wait()

		close(ch)
	}()

	keys := make([]string, count)
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}
