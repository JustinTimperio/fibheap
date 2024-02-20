<p align="center">
  <img width=400 src="./docs/fibpattern.jpg">
</p>

<h1 align="center">
    FibHeap - A pure Go implementation of Fibonacci Heaps
</h1>



This package was originally based on the work of a fairly old lib created by starwander called [GoFibonacciHeap](https://github.com/starwander/GoFibonacciHeap). The original package is coming up on almost a decade old now and hasn't been touched in the past 5 years. For this reason, I have detached a fork and updated the guts to include some priority features:
- Instead of using a slow and costly `interface{}` to store values we now use the generics in Go to allow for native types in the heap.
- Concurrency safety was missing in the original but now single heaps are protected with a mutex.
- Previously data structs had to conform to a interface spec. Now data can be arbitrarily added and removed with no need to conform to a interface spec.
- The test suites had fallen quite out of date and are now fully upgraded to work with ginkgo/v2.
- Code layout, organization and ergonomics have been greatly improved.
- The original was created before the standardization of go.mod and go.sum for packages. These have been added.

This implementation is a bit different from the traditional Fibonacci Heap with an index map inside. Thanks to the index map, the internal struct 'node' no longer need to be exposed outsides the package. The index map also makes the random access to the values in the heap possible. The union operation of this implementation is O(n) rather than O(1) of the traditional implementation.

| Operations                 | Insert | Minimum | ExtractMin | Union | DecreasePriority | IncreasePriority | Delete    | Get  |
| :------------------------: | :----: | :-----: | :--------: | :---: | :--------------: | :--------------: | :-------: | :--: |
| Traditional Implementation | O(1)   | O(1)    | O(log n)¹  | O(1)  | O(1)¹            | O(1)¹            | O(log n)¹ | N/A  |
| This Implementation        | O(1)   | O(1)    | O(log n)¹  | O(n)  | O(1)¹            | O(1)¹            | O(log n)¹ | O(1) |


## Operations

- `NewFibHeap[t any]() *FibHeap[t]`: Creates and initializes a new Fibonacci Heap.
- `Num() uint`: Returns the total number of values in the heap.
- `Insert(data t, priority float64) error`: Inserts a new value with the given data and priority into the heap.
- `Minimum() (data t, f float64)`: Returns the current minimum data and priority in the heap.
- `ExtractMin() (data t, f float64)`: Returns the current minimum data and priority in the heap and then extracts them from the heap.
- `Union(anotherHeap *FibHeap[t]) error`: Merges the input heap into the target heap.
- `DecreasePriority(data t, priority float64) error`: Decreases the priority of the value with the given data in the heap.
- `IncreasePriority(data t, priority float64) error`: Increases the priority of the value with the given data in the heap.
- `Delete(data t) error`: Removes the value with the given data from the heap.
- `GetPriority(data t) (priority float64)`: Returns the priority of the value with the given data in the heap.
- `ExtractPriority(data t) (priority float64)`: Returns the priority of the value with the given data in the heap and then extracts it from the heap.
- `Extract(data t) (t, float64)`: Returns the data and priority of the value with the given data in the heap and then extracts it from the heap.
- `Stats() string`: Returns some basic debug information about the heap.



## Example
```go

type SchoolEntry struct {
	Name string
	Age  float64
	Type string
}

func main() {

	heap := fibheap.NewFibHeap[SchoolEntry]()
	heap2 := fibheap.NewFibHeap[SchoolEntry]()

	s1 := SchoolEntry{"John", 18.3, "student"}
	s2 := SchoolEntry{"Tom", 21.0, "student"}
	s3 := SchoolEntry{"Jessica", 19.4, "student"}
	s4 := SchoolEntry{"Amy", 23.1, "student"}

	t1 := SchoolEntry{"Jason", 10.0, "teacher"}
	t2 := SchoolEntry{"Jack", 25.0, "teacher"}
	t3 := SchoolEntry{"Ryan", 28.0, "teacher"}

	heap.Insert(s1, s1.Age)
	heap.Insert(s2, s2.Age)
	heap.Insert(s3, s3.Age)
	heap.Insert(s4, s4.Age)

	fmt.Println(heap.Num())     // 4
	fmt.Println(heap.Minimum()) // {John 18.3 student} 18.3
	fmt.Println(heap.Num())     // 4

	heap.IncreasePriority(s1, 20.0)
	fmt.Println(heap.ExtractMin()) // {Jessica 19.4 student} 19.4

	fmt.Println(heap.ExtractMin()) // {John 18.3 student} 20
	fmt.Println(heap.Num())        // 2

	heap.DecreasePriority(s4, 16.5)
	fmt.Println(heap.ExtractMin()) // {Amy 23.1 student} 16.5

	fmt.Println(heap.Num())       // 1
	fmt.Println(heap.Extract(s2)) // {Tom 21 student} 21
	fmt.Println(heap.Num())       // 0

	heap.Insert(s1, s1.Age)
	heap.Insert(s2, s2.Age)
	heap.Insert(s3, s3.Age)
	heap.Insert(s4, s4.Age)
	heap2.Insert(t1, t1.Age)
	heap2.Insert(t2, t2.Age)
	heap2.Insert(t3, t3.Age)

	heap.Union(heap2)

	fmt.Println(heap.Num())     // 7
	fmt.Println(heap.Minimum()) // {Jason 10 teacher} 10
	for heap.Num() > 1 {
		heap.ExtractMin()
	}
	fmt.Println(heap.Num())     // 1
	fmt.Println(heap.Minimum()) // {Ryan 28 teacher} 28
}
```

## Why You Should NOT Use This!

This package is astronomically slow compared to the [standard package](https://pkg.go.dev/container/heap#example-package-PriorityQueue) and there really isn't any point to use this. Why then did you do this? Because it was fun and sometimes that is all that matters. Performance is about 5x slower than the standard heap which was honestly better than I had expected.


```
goos: linux
goarch: amd64
pkg: github.com/JustinTimperio/fibheap
cpu: AMD Ryzen 9 5900X 12-Core Processor            
BenchmarkFibHeap-24             1000000000               0.3370 ns/op
BenchmarkStandardLibHeap-24     1000000000               0.06133 ns/op
PASS
ok      github.com/JustinTimperio/fibheap       6.526s
```


Below is a much better priority queue just using the standard lib:

```go

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
```