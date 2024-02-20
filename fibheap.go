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

// Insert inserts a new value with the given data and priority into the heap.
// Returns an error if the insertion fails.
func (heap *FibHeap[t]) Insert(data t, priority float64) error {
	return heap.insert(data, priority)
}

// Minimum returns the current minimum data and priority in the heap.
// Returns -inf if the heap is empty.
func (heap *FibHeap[t]) Minimum() (data t, f float64) {
	if heap.num == 0 {
		return data, math.Inf(-1)
	}

	return heap.min.data, heap.min.priority
}

// ExtractMin returns the current minimum data and priority in the heap and then extracts them from the heap.
// Returns nil/-inf if the heap is empty.
func (heap *FibHeap[t]) ExtractMin() (data t, f float64) {
	if heap.num == 0 {
		return data, math.Inf(-1)
	}

	min := heap.extractMin()
	return min.data, min.priority
}

// Union merges the input heap into the target heap.
// Returns an error if any duplicate data are found in the target heap.
func (heap *FibHeap[t]) Union(anotherHeap *FibHeap[t]) error {
	for data := range anotherHeap.index {
		if _, exists := heap.index[data]; exists {
			return errors.New("Duplicate data is found in the target heap")
		}
	}

	for _, node := range anotherHeap.index {
		heap.insert(node.data, node.priority)
	}

	return nil
}

// DecreaseKey decreases the priority of the value with the given data in the heap.
// Returns an error if the value is not found or the priority is negative infinity.
func (heap *FibHeap[t]) DecreaseKey(data t, priority float64) error {
	if math.IsInf(priority, -1) {
		return errors.New("Negative infinity priority is reserved for internal usage")
	}

	if node, exists := heap.index[data]; exists {
		return heap.decreaseKey(node, priority)
	}

	return errors.New("Value is not found")
}

// IncreaseKey increases the priority of the value with the given data in the heap.
// Returns an error if the value is not found or the priority is negative infinity.
func (heap *FibHeap[t]) IncreaseKey(data t, priority float64) error {
	if math.IsInf(priority, -1) {
		return errors.New("Negative infinity priority is reserved for internal usage")
	}

	if node, exists := heap.index[data]; exists {
		return heap.increaseKey(node, priority)
	}

	return errors.New("Value is not found")
}

// Delete removes the value with the given data from the heap.
// Returns an error if the data is not found.
func (heap *FibHeap[t]) Delete(data t) error {
	if _, exists := heap.index[data]; !exists {
		return errors.New("Tag is not found")
	}

	heap.Extract(data)

	return nil
}

// GetPriority returns the priority of the value with the given data in the heap.
// Returns -inf if the value is not found.
func (heap *FibHeap[t]) GetPriority(data t) (priority float64) {
	if node, exists := heap.index[data]; exists {
		return node.priority
	}

	return math.Inf(-1)
}

// ExtractPriority returns the priority of the value with the given data in the heap and then extracts it from the heap.
// Returns -inf if the value is not found.
func (heap *FibHeap[t]) ExtractPriority(data t) (priority float64) {
	if node, exists := heap.index[data]; exists {
		priority = node.priority
		heap.deleteNode(node)
		return
	}

	return math.Inf(-1)
}

// ExtractValue returns the data and priority of the value with the given data in the heap and then extracts it from the heap.
// Returns the original data and -inf if the value is not found.
func (heap *FibHeap[t]) Extract(data t) (t, float64) {
	if node, exists := heap.index[data]; exists {
		k := node.priority
		v := node.data
		heap.deleteNode(node)
		return v, k
	}

	return data, math.Inf(-1)
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
	buffer.WriteString(fmt.Sprintf("Current min: priority(%f), data(%v),\n", heap.min.priority, heap.min.data))
	buffer.WriteString(fmt.Sprintf("Heap detail:\n"))
	probeTree[t](&buffer, heap.roots)
	buffer.WriteString(fmt.Sprintf("\n"))
	return buffer.String()
}
