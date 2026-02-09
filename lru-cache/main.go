package main

import (
	"fmt"
	"sync"
)

type Node struct {
	key  int
	val  int
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
	prev := node.prev
	prev.next = l.tail
	l.tail.prev = prev

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

	if p == l.tail {
		return
	}

	if p.prev == l.head {
		return
	}

	p.prev.next = p.next
	p.next.prev = p.prev

	first := l.head.next
	l.head.next = p
	p.prev = l.head
	p.next = first
	first.prev = p
}

type LRUCache struct {
	capacity int
	cache    map[int]*Node
	dll      *DoublyLinkedList
	mu       sync.RWMutex
}

func (c *LRUCache) get(key int) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache[key] == nil {
		fmt.Println("No Key found")
		return 0
	}
	node := c.cache[key]
	val := node.val
	c.dll.moveToFront(key)
	return val
}

func (c *LRUCache) put(key int, val int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache[key] != nil {
		node := c.cache[key]
		node.val = val
		c.dll.moveToFront(key)
		return
	}
	if c.dll.len() == c.capacity {
		tail := c.dll.removeFromTail()
		delete(c.cache, tail.key)
	}

	node := &Node{key: key, val: val}
	c.cache[key] = node
	c.dll.addNode(node)
}

func NewLRUCache(capacity int) *LRUCache {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*Node),
		dll:      &DoublyLinkedList{head: head, tail: tail},
	}
}

func main() {
	fmt.Println("=== LRU Cache Simulation ===")

	// Create cache with capacity 3
	cache := NewLRUCache(3)
	fmt.Println("Created LRU Cache with capacity 3")

	// Test 1: Add items
	fmt.Println("Test 1: Adding items")
	cache.put(1, 100)
	fmt.Println("Put(1, 100)")
	cache.put(2, 200)
	fmt.Println("Put(2, 200)")
	cache.put(3, 300)
	fmt.Println("Put(3, 300)")
	fmt.Printf("Cache size: %d\n\n", cache.dll.len())

	// Test 2: Get items
	fmt.Println("Test 2: Getting items")
	val := cache.get(1)
	fmt.Printf("Get(1) = %d (moves key 1 to front)\n", val)
	val = cache.get(2)
	fmt.Printf("Get(2) = %d (moves key 2 to front)\n\n", val)

	// Test 3: Eviction - add item when cache is full
	fmt.Println("Test 3: Cache eviction")
	fmt.Println("Adding key 4 (cache is full, should evict LRU item)")
	cache.put(4, 400)
	fmt.Println("Put(4, 400) - key 3 should be evicted")

	val = cache.get(3)
	fmt.Printf("Get(3) = %d (should return 0 - key was evicted)\n\n", val)

	// Test 4: Update existing key
	fmt.Println("Test 4: Update existing key")
	cache.put(2, 250)
	fmt.Println("Put(2, 250) - updating existing key")
	val = cache.get(2)
	fmt.Printf("Get(2) = %d (should return 250)\n\n", val)

	// Test 5: Add more items to test further eviction
	fmt.Println("Test 5: More evictions")
	cache.put(5, 500)
	fmt.Println("Put(5, 500) - should evict key 1 (LRU)")
	val = cache.get(1)
	fmt.Printf("Get(1) = %d (should return 0 - key was evicted)\n", val)

	fmt.Println("\nCurrent cache contents (MRU to LRU): 5, 2, 4")
	val = cache.get(5)
	fmt.Printf("Get(5) = %d\n", val)
	val = cache.get(2)
	fmt.Printf("Get(2) = %d\n", val)
	val = cache.get(4)
	fmt.Printf("Get(4) = %d\n", val)

	fmt.Println("\n=== Simulation Complete ===")
}
