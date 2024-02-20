package fibheap

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"math"
	"sync"
)

// NewFibHeap creates an initialized Fibonacci Heap.
func NewFibHeap[t any]() *FibHeap[t] {
	// Create a new instance of FibHeap
	heap := new(FibHeap[t])
	// Initialize the roots list
	heap.roots = list.New()
	// Initialize the index map
	heap.index = make(map[interface{}]*node[t])
	// Initialize the treeDegrees map
	heap.treeDegrees = make(map[uint]*list.Element)
	// Initialize the number of values in the heap
	heap.num = 0
	// Initialize the minimum node
	heap.min = nil
	// Initialize the mutex for thread-safety
	heap.mutex = sync.Mutex{}

	return heap
}

// Num returns the total number of values in the heap.
func (heap *FibHeap[t]) Num() uint {
	return heap.num
}

// Insert inserts a new value with the given tag and key into the heap.
// Returns an error if the insertion fails.
func (heap *FibHeap[t]) Insert(tag t, key float64) error {
	return heap.insert(tag, key)
}

// Minimum returns the current minimum tag and key in the heap.
// Returns -inf if the heap is empty.
func (heap *FibHeap[t]) Minimum() (tag t, f float64) {
	if heap.num == 0 {
		return tag, math.Inf(-1)
	}

	return heap.min.tag, heap.min.key
}

// ExtractMin returns the current minimum tag and key in the heap and then extracts them from the heap.
// Returns nil/-inf if the heap is empty.
func (heap *FibHeap[t]) ExtractMin() (tag t, f float64) {
	if heap.num == 0 {
		return tag, math.Inf(-1)
	}

	min := heap.extractMin()
	return min.tag, min.key
}

// Union merges the input heap into the target heap.
// Returns an error if any duplicate tags are found in the target heap.
func (heap *FibHeap[t]) Union(anotherHeap *FibHeap[t]) error {
	for tag := range anotherHeap.index {
		if _, exists := heap.index[tag]; exists {
			return errors.New("Duplicate tag is found in the target heap")
		}
	}

	for _, node := range anotherHeap.index {
		heap.insert(node.tag, node.key)
	}

	return nil
}

// DecreaseKey decreases the key of the value with the given tag in the heap.
// Returns an error if the value is not found or the key is negative infinity.
func (heap *FibHeap[t]) DecreaseKey(tag t, key float64) error {
	if math.IsInf(key, -1) {
		return errors.New("Negative infinity key is reserved for internal usage")
	}

	if node, exists := heap.index[tag]; exists {
		return heap.decreaseKey(node, key)
	}

	return errors.New("Value is not found")
}

// IncreaseKey increases the key of the value with the given tag in the heap.
// Returns an error if the value is not found or the key is negative infinity.
func (heap *FibHeap[t]) IncreaseKey(tag t, key float64) error {
	if math.IsInf(key, -1) {
		return errors.New("Negative infinity key is reserved for internal usage")
	}

	if node, exists := heap.index[tag]; exists {
		return heap.increaseKey(node, key)
	}

	return errors.New("Value is not found")
}

// Delete removes the value with the given tag from the heap.
// Returns an error if the tag is not found.
func (heap *FibHeap[t]) Delete(tag t) error {
	if _, exists := heap.index[tag]; !exists {
		return errors.New("Tag is not found")
	}

	heap.ExtractValue(tag)

	return nil
}

// GetTag returns the key of the value with the given tag in the heap.
// Returns -inf if the value is not found.
func (heap *FibHeap[t]) GetTag(tag t) (key float64) {
	if node, exists := heap.index[tag]; exists {
		return node.key
	}

	return math.Inf(-1)
}

// ExtractTag returns the key of the value with the given tag in the heap and then extracts it from the heap.
// Returns -inf if the value is not found.
func (heap *FibHeap[t]) ExtractTag(tag t) (key float64) {
	if node, exists := heap.index[tag]; exists {
		key = node.key
		heap.deleteNode(node)
		return
	}

	return math.Inf(-1)
}

// ExtractValue returns the tag and key of the value with the given tag in the heap and then extracts it from the heap.
// Returns the original tag and -inf if the value is not found.
func (heap *FibHeap[t]) ExtractValue(tag t) (t, float64) {
	if node, exists := heap.index[tag]; exists {
		k := node.key
		v := node.tag
		heap.deleteNode(node)
		return v, k
	}

	return tag, math.Inf(-1)
}

// Stats returns some basic debug information about the heap.
// It includes the total number of values, the size of the roots list, the size of the index map,
// and the current minimum value in the heap.
// It also includes the topology of the trees in the heap using a depth-first search.
func (heap *FibHeap[t]) Stats() string {
	var buffer bytes.Buffer

	if heap.num == 0 {
		buffer.WriteString(fmt.Sprintf("Heap is empty.\n"))
		return buffer.String()
	}

	buffer.WriteString(fmt.Sprintf("Total number: %d, Root Size: %d, Index size: %d,\n", heap.num, heap.roots.Len(), len(heap.index)))
	buffer.WriteString(fmt.Sprintf("Current min: key(%f), tag(%v),\n", heap.min.key, heap.min.tag))
	buffer.WriteString(fmt.Sprintf("Heap detail:\n"))
	probeTree[t](&buffer, heap.roots)
	buffer.WriteString(fmt.Sprintf("\n"))
	return buffer.String()
}
