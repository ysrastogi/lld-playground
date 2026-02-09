package main

import "testing"

func TestNewLRUCache(t *testing.T) {
	cache := NewLRUCache(5)
	if cache == nil {
		t.Fatal("NewLRUCache returned nil")
	}
	if cache.capacity != 5 {
		t.Errorf("Expected capacity 5, got %d", cache.capacity)
	}
	if cache.dll.len() != 0 {
		t.Errorf("Expected empty cache, got length %d", cache.dll.len())
	}
}

func TestPutAndGet(t *testing.T) {
	cache := NewLRUCache(3)

	cache.put(1, 100)
	cache.put(2, 200)
	cache.put(3, 300)

	if val := cache.get(1); val != 100 {
		t.Errorf("Expected 100, got %d", val)
	}
	if val := cache.get(2); val != 200 {
		t.Errorf("Expected 200, got %d", val)
	}
	if val := cache.get(3); val != 300 {
		t.Errorf("Expected 300, got %d", val)
	}
}

func TestGetNonExistentKey(t *testing.T) {
	cache := NewLRUCache(3)

	cache.put(1, 100)

	if val := cache.get(999); val != 0 {
		t.Errorf("Expected 0 for non-existent key, got %d", val)
	}
}

func TestUpdateExistingKey(t *testing.T) {
	cache := NewLRUCache(3)

	cache.put(1, 100)
	cache.put(1, 150)

	if val := cache.get(1); val != 150 {
		t.Errorf("Expected 150 after update, got %d", val)
	}
}

func TestEviction(t *testing.T) {
	cache := NewLRUCache(3)

	cache.put(1, 100)
	cache.put(2, 200)
	cache.put(3, 300)

	cache.put(4, 400)

	if val := cache.get(1); val != 0 {
		t.Errorf("Expected key 1 to be evicted, but got value %d", val)
	}

	if val := cache.get(2); val != 200 {
		t.Errorf("Expected key 2 to exist with value 200, got %d", val)
	}
	if val := cache.get(3); val != 300 {
		t.Errorf("Expected key 3 to exist with value 300, got %d", val)
	}
	if val := cache.get(4); val != 400 {
		t.Errorf("Expected key 4 to exist with value 400, got %d", val)
	}
}

func TestLRUOrdering(t *testing.T) {
	cache := NewLRUCache(3)

	cache.put(1, 100)
	cache.put(2, 200)
	cache.put(3, 300)

	cache.get(1)
	cache.put(4, 400)

	if val := cache.get(2); val != 0 {
		t.Errorf("Expected key 2 to be evicted, got %d", val)
	}
	if val := cache.get(1); val != 100 {
		t.Errorf("Expected key 1 to exist, got %d", val)
	}
	if val := cache.get(3); val != 300 {
		t.Errorf("Expected key 3 to exist, got %d", val)
	}
	if val := cache.get(4); val != 400 {
		t.Errorf("Expected key 4 to exist, got %d", val)
	}
}

func TestCapacityOne(t *testing.T) {
	cache := NewLRUCache(1)

	cache.put(1, 100)
	if val := cache.get(1); val != 100 {
		t.Errorf("Expected 100, got %d", val)
	}

	cache.put(2, 200)

	if val := cache.get(1); val != 0 {
		t.Errorf("Expected key 1 to be evicted, got %d", val)
	}
	if val := cache.get(2); val != 200 {
		t.Errorf("Expected key 2 with value 200, got %d", val)
	}
}

func TestMultipleUpdates(t *testing.T) {
	cache := NewLRUCache(2)

	cache.put(1, 100)
	cache.put(2, 200)
	cache.put(1, 150)
	cache.put(3, 300)

	if val := cache.get(2); val != 0 {
		t.Errorf("Expected key 2 to be evicted, got %d", val)
	}
	if val := cache.get(1); val != 150 {
		t.Errorf("Expected key 1 with value 150, got %d", val)
	}
	if val := cache.get(3); val != 300 {
		t.Errorf("Expected key 3 with value 300, got %d", val)
	}
}

func TestCacheLength(t *testing.T) {
	cache := NewLRUCache(5)

	if cache.dll.len() != 0 {
		t.Errorf("Expected initial length 0, got %d", cache.dll.len())
	}

	cache.put(1, 100)
	if cache.dll.len() != 1 {
		t.Errorf("Expected length 1, got %d", cache.dll.len())
	}

	cache.put(2, 200)
	cache.put(3, 300)

	if cache.dll.len() != 3 {
		t.Errorf("Expected length 3, got %d", cache.dll.len())
	}

	cache.put(4, 400)
	cache.put(5, 500)
	cache.put(6, 600)

	if cache.dll.len() != 5 {
		t.Errorf("Expected length to stay at capacity 5, got %d", cache.dll.len())
	}
}

func TestSequentialOperations(t *testing.T) {
	cache := NewLRUCache(2)

	cache.put(1, 1)
	cache.put(2, 2)

	if val := cache.get(1); val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}

	cache.put(3, 3)

	if val := cache.get(2); val != 0 {
		t.Errorf("Expected key 2 to be evicted, got %d", val)
	}

	cache.put(4, 4)

	if val := cache.get(1); val != 0 {
		t.Errorf("Expected key 1 to be evicted, got %d", val)
	}
	if val := cache.get(3); val != 3 {
		t.Errorf("Expected 3, got %d", val)
	}
	if val := cache.get(4); val != 4 {
		t.Errorf("Expected 4, got %d", val)
	}
}
