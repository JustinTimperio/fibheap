<p align="center">
  <img width=400 src="./docs/fibpattern.jpg">
</p>

<h1 align="center">
    FibHeap - A pure Go implementation of Fibonacci Heaps
</h1>

This package was originally based on the work of a fairly old lib that by starwander called [GoFibonacciHeap](https://github.com/starwander/GoFibonacciHeap). The original package is coming up on almost a decade old now and hasn't been touched in the past 5 years. For this reason, I have detached a fork and updated the guts to include some key features:
- Instead of using a slow and costly `interface{}` to store values we now use the generics in Go to allow for native types in the heap.:w
- 
- The test suites had fallen quite out of date and are now fully upgraded to work with ginkgo/v2.
- Code layout, organization and ergonomics have been greatly improved.
- The original was created before the standardization of go.mod and go.sum for packages. These have been added.

This implementation is a bit different from the traditional Fibonacci Heap with an index map inside. Thanks to the index map, the internal struct 'node' no longer need to be exposed outsides the package. The index map also makes the random access to the values in the heap possible. The union operation of this implementation is O(n) rather than O(1) of the traditional implementation.

| Operations                 | Insert | Minimum | ExtractMin | Union | DecreaseKey | IncreaseKey | Delete    | Get  |
| :------------------------: | :----: | :-----: | :--------: | :---: | :---------: | :---------: | :-------: | :--: |
| Traditional Implementation | O(1)   | O(1)    | O(log n)¹  | O(1)  | O(1)¹       | O(1)¹       | O(log n)¹ | N/A  |
| This Implementation        | O(1)   | O(1)    | O(log n)¹  | O(n)  | O(1)¹       | O(1)¹       | O(log n)¹ | O(1) |


## Operations

- `NewFibHeap[t any]() *FibHeap[t]`: Creates and initializes a new Fibonacci Heap.
- `Num() uint`: Returns the total number of values in the heap.
- `Insert(tag t, key float64) error`: Inserts a new value with the given tag and key into the heap.
- `Minimum() (tag t, f float64)`: Returns the current minimum tag and key in the heap.
- `ExtractMin() (tag t, f float64)`: Returns the current minimum tag and key in the heap and then extracts them from the heap.
- `Union(anotherHeap *FibHeap[t]) error`: Merges the input heap into the target heap.
- `DecreaseKey(tag t, key float64) error`: Decreases the key of the value with the given tag in the heap.
- `IncreaseKey(tag t, key float64) error`: Increases the key of the value with the given tag in the heap.
- `Delete(tag t) error`: Removes the value with the given tag from the heap.
- `GetTag(tag t) (key float64)`: Returns the key of the value with the given tag in the heap.
- `ExtractTag(tag t) (key float64)`: Returns the key of the value with the given tag in the heap and then extracts it from the heap.
- `ExtractValue(tag t) (t, float64)`: Returns the tag and key of the value with the given tag in the heap and then extracts it from the heap.
- `Stats() string`: Returns some basic debug information about the heap.



## Example
```go

def main() {
    heap := fibheap.NewFibHeap[student]()
	heap2 := fibheap.NewFibHeap[student]()

	s1 := student{"John", 18.3}
	s2 := student{"Tom", 21.0}
	s3 := student{"Jessica", 19.4}
	s4 := student{"Amy", 23.1}

	t1 := student{"Jason", 10.0}
	t2 := student{"Jack", 25.0}
	t3 := student{"Ryan", 28.0}

	heap.Insert(s1, s1.Age)
	heap.Insert(s2, s2.Age)
	heap.Insert(s3, s3.Age)
	heap.Insert(s4, s4.Age)

	fmt.Println(heap.Num())     
	fmt.Println(heap.Minimum()) 
	fmt.Println(heap.Num())     

	heap.IncreaseKey(s1, 20.0)
	fmt.Println(heap.ExtractMin()) 

	fmt.Println(heap.ExtractMin()) 
	fmt.Println(heap.Num())        

	heap.DecreaseKey(s4, 16.5)
	fmt.Println(heap.ExtractMin()) 

	fmt.Println(heap.Num())            
	fmt.Println(heap.ExtractValue(s2))
	fmt.Println(heap.Num())            

	heap.Insert(s1, s1.Age)
	heap.Insert(s2, s2.Age)
	heap.Insert(s3, s3.Age)
	heap.Insert(s4, s4.Age)
	heap2.Insert(t1, t1.Age)
	heap2.Insert(t2, t2.Age)
	heap2.Insert(t3, t3.Age)

	heap.Union(heap2)

	fmt.Println(heap.Num())
	fmt.Println(heap.Minimum())
	for heap.Num() > 1 {
		heap.ExtractMin()
	}
	fmt.Println(heap.Num()) 
	fmt.Println(heap.Minimum()) 
}

```