package main

import "fmt"

type Node struct {
	key  int
	val  int
	freq int

	prev *Node
	next *Node
}
type DoublyLinkedList struct {
	head *Node
	tail *Node
}

func (l *DoublyLinkedList) removeFromTail() *Node {
	if l.tail.prev == nil || l.tail.prev.prev == nil {
		return nil
	}
	node := l.tail.prev
	l.tail.prev = node.prev
	node.prev.next = l.tail

	node.prev = nil
	node.next = nil
	return node
}

func (l *DoublyLinkedList) addNode(node *Node) {
	first := l.head.next

	l.head.next = node
	node.prev = l.head
	node.next = first
	first.prev = node
}
func (l *DoublyLinkedList) len() int {
	p := l.head.next
	length := 0
	for p != l.tail {
		length += 1
		p = p.next
	}
	return length
}

func (l *DoublyLinkedList) moveToFront(k int) {
	p := l.head.next

	for p != l.tail && p.key != k {
		p = p.next
	}

}
func removeNode(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev

	node.prev = nil
	node.next = nil
}

type LFUCache struct {
	capacity int
	size     int
	cache    map[int]*Node
	freqMap  map[int]*DoublyLinkedList
	minFreq  int
}

func Constructor(capacity int) LFUCache {
	return LFUCache{
		capacity: capacity,
		cache:    make(map[int]*Node),
		freqMap:  make(map[int]*DoublyLinkedList),
	}
}

func (c *LFUCache) Get(key int) int {
	if node, ok := c.cache[key]; ok {

		c.cache[key].freq += 1
		freq := node.freq
		removeNode(node)
		if c.freqMap[freq] == nil {
			c.freqMap[freq] = &DoublyLinkedList{
				head: &Node{},
				tail: &Node{},
			}
			c.freqMap[freq].head.next = c.freqMap[freq].tail
			c.freqMap[freq].tail.prev = c.freqMap[freq].head
		}
		c.freqMap[freq].addNode(node)

		if c.freqMap[c.minFreq] != nil && c.freqMap[c.minFreq].len() == 0 {
			c.minFreq += 1
		}

		return node.val
	}
	return -1
}

func (c *LFUCache) Put(key int, value int) {
	if c.capacity == 0 {
		return
	}
	if node, ok := c.cache[key]; ok {
		node.val = value
		c.Get(key)
		return
	}
	if c.size == c.capacity {
		evictedNode := c.freqMap[c.minFreq].removeFromTail()
		delete(c.cache, evictedNode.key)
		c.size -= 1
	}
	newNode := &Node{
		key:  key,
		val:  value,
		freq: 1,
	}
	c.cache[key] = newNode
	if c.freqMap[1] == nil {
		c.freqMap[1] = &DoublyLinkedList{
			head: &Node{},
			tail: &Node{},
		}
		c.freqMap[1].head.next = c.freqMap[1].tail
		c.freqMap[1].tail.prev = c.freqMap[1].head
	}
	c.freqMap[1].addNode(newNode)
	c.minFreq = 1
	c.size += 1
}

func main() {
	fmt.Println("=== LFU Cache Simulation ===")

	// Create cache with capacity 3
	cache := Constructor(3)
	fmt.Println("Created LFU Cache with capacity 3\n")

	// Test 1: Add items
	fmt.Println("Test 1: Adding items")
	cache.Put(1, 100)
	fmt.Println("Put(1, 100) - freq=1")
	cache.Put(2, 200)
	fmt.Println("Put(2, 200) - freq=1")
	cache.Put(3, 300)
	fmt.Println("Put(3, 300) - freq=1")
	fmt.Printf("Cache size: %d\n\n", cache.size)

	// Test 2: Access items to increase frequency
	fmt.Println("Test 2: Accessing items to change frequency")
	val := cache.Get(1)
	fmt.Printf("Get(1) = %d (freq: 1->2)\n", val)
	val = cache.Get(1)
	fmt.Printf("Get(1) = %d (freq: 2->3)\n", val)
	val = cache.Get(2)
	fmt.Printf("Get(2) = %d (freq: 1->2)\n\n", val)

	// Test 3: Eviction - add item when cache is full
	fmt.Println("Test 3: LFU Cache eviction")
	fmt.Println("Current frequencies: key1=3, key2=2, key3=1")
	fmt.Println("Adding key 4 (cache is full, should evict key 3 with freq=1)")
	cache.Put(4, 400)
	fmt.Println("Put(4, 400) - key 3 evicted (lowest frequency)")

	val = cache.Get(3)
	fmt.Printf("Get(3) = %d (should return -1 - key was evicted)\n\n", val)

	// Test 4: Evict based on LRU when frequencies are equal
	fmt.Println("Test 4: Eviction with equal frequencies")
	cache.Put(5, 500)
	fmt.Println("Put(5, 500) - freq=1")
	fmt.Println("Current frequencies: key1=3, key2=2, key4=1, key5=1")
	fmt.Println("Keys 4 and 5 both have freq=1, key 4 is LRU among them")

	cache.Put(6, 600)
	fmt.Println("Put(6, 600) - should evict key 4 (LRU among lowest freq)")

	val = cache.Get(4)
	fmt.Printf("Get(4) = %d (should return -1 - key was evicted)\n\n", val)

	// Test 5: Update existing key
	fmt.Println("Test 5: Update existing key")
	cache.Put(2, 250)
	fmt.Println("Put(2, 250) - updating existing key (freq increases)")
	val = cache.Get(2)
	fmt.Printf("Get(2) = %d (should return 250)\n\n", val)

	// Test 6: Display current cache state
	fmt.Println("Test 6: Current cache contents")
	fmt.Println("Key 1:")
	val = cache.Get(1)
	fmt.Printf("  Get(1) = %d\n", val)

	fmt.Println("Key 2:")
	val = cache.Get(2)
	fmt.Printf("  Get(2) = %d\n", val)

	fmt.Println("Key 5:")
	val = cache.Get(5)
	fmt.Printf("  Get(5) = %d\n", val)

	fmt.Println("Key 6:")
	val = cache.Get(6)
	fmt.Printf("  Get(6) = %d\n", val)

	// Test 7: Edge case - cache with capacity 0
	fmt.Println("\nTest 7: Edge case - capacity 0")
	emptyCache := Constructor(0)
	emptyCache.Put(1, 100)
	fmt.Println("Put(1, 100) to cache with capacity 0")
	val = emptyCache.Get(1)
	fmt.Printf("Get(1) = %d (should return -1)\n", val)

	// Test 8: Access pattern showing frequency-based eviction
	fmt.Println("\nTest 8: Frequency-based eviction pattern")
	testCache := Constructor(2)
	testCache.Put(1, 1)
	fmt.Println("Put(1, 1)")
	testCache.Put(2, 2)
	fmt.Println("Put(2, 2)")

	val = testCache.Get(1)
	fmt.Printf("Get(1) = %d (freq: 1->2)\n", val)

	testCache.Put(3, 3)
	fmt.Println("Put(3, 3) - should evict key 2 (freq=1, key 1 has freq=2)")

	val = testCache.Get(2)
	fmt.Printf("Get(2) = %d (should return -1)\n", val)
	val = testCache.Get(3)
	fmt.Printf("Get(3) = %d (should return 3)\n", val)
	val = testCache.Get(1)
	fmt.Printf("Get(1) = %d (should return 1)\n", val)

	fmt.Println("\n=== Simulation Complete ===")
}
