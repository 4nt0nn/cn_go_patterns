package sharding

import (
	"crypto/sha1"
	"sync"
)

// Shard represents the data partition
type Shard struct {
	sync.RWMutex // Compose from sync.RWMutex
	m map[string]interface{} // m contains the shard's data
}

// ShardedMap is a *Shards slice of pointers to some number of Shard values.
// Each Shard includes a map[string]interface{} that contains that shard's data, and
// a composed sync.RWMutex so that it can be individually locked
type ShardedMap []*Shard 

// NewShardedMap is the constructor function for ShardedMap
func NewShardedMap(nshards int) ShardedMap {
	shards := make([]*Shard, nshards) // Initialize a *Shards slice

	for i := 0; i < nshards; i++ {
		shard := make(map[string]interface{})
		shards[i] = &Shard{m: shard}
	}

	return shards // A ShardedMap IS a *Shards slice!
}

func (m ShardedMap) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key)) // Use sum from "crypto/sha1"
	hash := int(checksum[17]) // Pick an arbitrary byte as the hash, downside of byte-sized value as hash value is that we can only handle up to 255 shards.
							  // If we need more, sprinkle some binary arithmetic on it: hash := int(sum[13]) << 8 | int(sum[17])
	
	return hash % len(m)
}

func (m ShardedMap) getShard(key string) *Shard {
	index := m.getShardIndex(key)
	return m[index]
}

// Get is an example of a util method attached to ShardedMap that allows users to get a
// value of from a shard based on key
func (m ShardedMap) Get(key string) interface{} {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

// Set is an example of a util method attached to ShardedMap that allows users to set a
// value to the shards map based on key
func (m ShardedMap) Set(key string, value interface{}) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = value
}

// Keys is used to establish locks on all of the tables concurrently.
func (m ShardedMap) Keys() []string {
	keys := make([]string, 0) // Create an empty keys slice

	mutex := sync.Mutex{} // Mutex for write safety to keys

	wg := sync.WaitGroup{} // Create a wait group and add a wait value for each slice
	wg.Add(len(m))
	
	for _, shard := range m { // Run a goroutine for each slice
		go func(s *Shard) {
			s.RLock() // Establish a read lock on s

			for key := range s.m { // Get the slice's keys
				mutex.Lock()
				keys = append(keys, key)
				mutex.Unlock()
			}

			s.RUnlock() // Release the read lock
			wg.Done() // Tell the WaitGroup it's done

		}(shard)
	}

	wg.Wait() // Block until all reads are done

	return keys
}


//// example use /////

// func main() {
// 	shardedMap := NewShardedMap(5)

// 	shardedMap.Set("alpha", 1)
// 	shardedMap.Set("beta", 2)
// 	shardedMap.Set("gamma", 3)

// 	fmt.Println(shardedMap.Get("alpha"))
// 	fmt.Println(shardedMap.Get("beta"))
// 	fmt.Println(shardedMap.Get("gamma"))

// 	keys := shardedMap.Keys()
// 	for _, k := range keys {
// 		fmt.Println(k)
// 	}
// }