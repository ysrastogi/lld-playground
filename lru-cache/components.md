# Components
The `lru-cache` package is built on top of a few core components that work together to provide efficient caching functionality. These components include:

1. **Cache Storage**: This component is responsible for storing the key-value pairs in the cache. It typically uses a hash map (or dictionary) to allow for O(1) time complexity for get and put operations.
2. **Doubly Linked List**: This component maintains the order of usage for the cache items. It allows for efficient eviction of the least recently used item when the cache reaches its capacity.
3. **Eviction Policy**: This component implements the logic for evicting items from the cache when the capacity is exceeded. In this case, it follows the Least Recently Used (LRU) eviction policy.
4. **Concurrency Control**: This component ensures that the cache can be safely accessed by multiple threads simultaneously without
