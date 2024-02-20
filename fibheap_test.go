package fibheap_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/JustinTimperio/fibheap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type demoStruct struct {
	data     int
	priority float64
	value    string
}

type SchoolEntry struct {
	Name string
	Age  float64
	Type string
}

func TestBasic(t *testing.T) {

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

func TestProxy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoFibonacciHeap Suite")
}

var _ = Describe("Tests of fibHeap", func() {
	var (
		heap        *fibheap.FibHeap[int]
		anotherHeap *fibheap.FibHeap[int]
		v1          = 1000
		v2          = 999
	)

	Context("behaviour tests of data/priority interfaces", func() {
		BeforeEach(func() {
			heap = fibheap.NewFibHeap[int]()
		})

		AfterEach(func() {
			heap = nil
		})

		It("Given an empty fibHeap, when call Minimum api, it should return nil.", func() {
			data, priority := heap.Minimum()
			Expect(data).Should(BeEquivalentTo(0))
			Expect(priority).Should(BeEquivalentTo(math.Inf(-1)))
		})

		It("Given a empty fibHeap, when call Insert api with a negative infinity priority, it should return error.", func() {
			Expect(heap.Insert(1000, math.Inf(-1))).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call Minimum api, it should return the minimum value inserted.", func() {
			min := math.Inf(1)
			for i := 0; i < 10000; i++ {
				priority := rand.Float64()
				heap.Insert(i, priority)
				if priority < min {
					min = priority
				}
			}

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			_, minKey := heap.Minimum()
			Expect(minKey).Should(BeEquivalentTo(min))
			Expect(heap.Num()).Should(BeEquivalentTo(10000))
		})

		It("Given an empty fibHeap, when call ExtractMin api, it should return nil.", func() {
			data, _ := heap.ExtractMin()
			Expect(data).Should(BeEquivalentTo(0))
		})

		It("Given a fibHeap inserted multiple values, when call ExtractMin api, it should extract the minimum value inserted.", func() {
			for i := 0; i < 10000; i++ {
				priority := rand.Float64()
				heap.Insert(i, priority)
			}

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			_, lastKey := heap.Minimum()
			for i := 0; i < 10000; i++ {
				_, priority := heap.ExtractMin()
				Expect(priority).Should(BeNumerically(">=", lastKey))
				Expect(heap.Num()).Should(BeEquivalentTo(9999 - i))
				lastKey = priority
			}
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})

		It("Given a fibHeap inserted multiple values, when call DecreasePriority api with a non-exists value, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}

			Expect(heap.DecreasePriority(v1, float64(999))).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap with a value, when call DecreasePriority api with a negative infinity priority, it should return error.", func() {
			heap.Insert(1000, float64(1000))
			Expect(heap.DecreasePriority(v1, math.Inf(-1))).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call DecreasePriority api with a larger priority, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}

			Expect(heap.DecreasePriority(v2, float64(1000))).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap inserted multiple values, when call IncreasePriority api with a non-exists value, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}

			Expect(heap.IncreasePriority(v1, float64(999))).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap with a value, when call IncreasePriority api with a negative infinity priority, it should return error.", func() {
			heap.Insert(1000, float64(1000))
			Expect(heap.IncreasePriority(v1, math.Inf(-1))).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call IncreasePriority api with a smaller priority, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}

			Expect(heap.IncreasePriority(v2, float64(998))).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap inserted multiple values, when call Delete api with a non-exists value, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
			v3 := 10000

			Expect(heap.Delete(v3)).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap inserted multiple values, when call Delete api, it should remove the value from the heap.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			for i := 0; i < 1000; i++ {
				Expect(heap.Delete(i)).ShouldNot(HaveOccurred())
			}
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})
	})

	Context("union tests", func() {
		BeforeEach(func() {
			heap = fibheap.NewFibHeap[int]()
			anotherHeap = fibheap.NewFibHeap[int]()
		})

		AfterEach(func() {
			heap = nil
			anotherHeap = nil
		})

		It("Given two empty fibHeaps, when call Union api, it should return an empty fibHeap.", func() {
			heap.Union(anotherHeap)
			n, f := heap.Minimum()
			Expect(f).Should(BeEquivalentTo(math.Inf(-1)))
			Expect(n).Should(BeEquivalentTo(0))
		})

		It("Given one empty fibHeap and one non-empty fibHeap, when Union the non-empty one into the empty one, it should retern a new heap with the number and min of the non-empty heap.", func() {
			for i := 0; i < int(rand.Int31n(1000)); i++ {
				demo := new(demoStruct)
				demo.data = i
				demo.priority = rand.Float64()
				demo.value = fmt.Sprint(demo.priority)
				anotherHeap.Insert(demo.data, demo.priority)
			}
			number := anotherHeap.Num()
			min, _ := anotherHeap.Minimum()

			heap.Union(anotherHeap)
			min2, _ := heap.Minimum()
			Expect(min2).Should(BeEquivalentTo(min))
			Expect(heap.Num()).Should(BeEquivalentTo(number))
		})

		It("Given one empty fibHeap and one non-empty fibHeap, when Union the empty one into the non-empty one, it should retern a new heap with the number and min of the non-empty heap.", func() {
			for i := 0; i < int(rand.Int31n(1000)); i++ {
				demo := new(demoStruct)
				demo.data = i
				demo.priority = rand.Float64()
				demo.value = fmt.Sprint(demo.priority)
				heap.Insert(demo.data, demo.priority)
			}
			number := heap.Num()
			min, _ := heap.Minimum()

			heap.Union(anotherHeap)
			min2, _ := heap.Minimum()
			Expect(min2).Should(BeEquivalentTo(min))
			Expect(heap.Num()).Should(BeEquivalentTo(number))
		})

		It("Given two fibHeap with multiple values, when call ExtractMin api after unioned, it should extract the minimum value inserted into both heaps.", func() {
			for i := 0; i < 5000; i++ {
				demo := new(demoStruct)
				demo.data = i
				demo.priority = rand.Float64()
				demo.value = fmt.Sprint(demo.priority)
				heap.Insert(demo.data, demo.priority)
			}

			for i := 5000; i < 10000; i++ {
				anotherdemo := new(demoStruct)
				anotherdemo.data = i
				anotherdemo.priority = rand.Float64()
				anotherdemo.value = fmt.Sprint(anotherdemo.priority)
				anotherHeap.Insert(anotherdemo.data, anotherdemo.priority)
			}

			_, min := heap.Minimum()
			_, amin := anotherHeap.Minimum()

			if amin < min {
				min = amin
			}
			heap.Union(anotherHeap)

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			_, lastKey := heap.Minimum()
			Expect(lastKey).Should(BeEquivalentTo(min))

			for i := 0; i < 10000; i++ {
				_, priority := heap.ExtractMin()
				Expect(priority).Should(BeNumerically(">=", lastKey))
				lastKey = priority
			}
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})
	})

	Context("index tests of data/priority interfaces", func() {
		BeforeEach(func() {
			heap = fibheap.NewFibHeap[int]()
			anotherHeap = fibheap.NewFibHeap[int]()
		})

		AfterEach(func() {
			heap = nil
			anotherHeap = nil
		})

		It("Given one fibHeap, when Insert values with same data, it should return an error.", func() {
			err := heap.Insert(1, float64(1))
			Expect(err).ShouldNot(HaveOccurred())
			_, minKey := heap.Minimum()
			Expect(minKey).Should(BeEquivalentTo(1))
			Expect(heap.Num()).Should(BeEquivalentTo(1))
			err = heap.Insert(1, float64(10))
			Expect(err).Should(HaveOccurred())
			_, minKey = heap.Minimum()
			Expect(minKey).Should(BeEquivalentTo(1))
			Expect(heap.Num()).Should(BeEquivalentTo(1))
		})

		It("Given two fibHeaps which both has value with same data, when call Union, it should return an error.", func() {
			heap.Insert(1, float64(1))
			anotherHeap.Insert(1, float64(10))

			err := heap.Union(anotherHeap)
			Expect(err).Should(HaveOccurred())
			_, minKey := heap.Minimum()
			Expect(minKey).Should(BeEquivalentTo(1))
			Expect(heap.Num()).Should(BeEquivalentTo(1))
			_, anotherMinKey := anotherHeap.Minimum()
			Expect(anotherMinKey).Should(BeEquivalentTo(10))
			Expect(anotherHeap.Num()).Should(BeEquivalentTo(1))
		})

		It("Given one fibHeaps which has not a value with TAG, when GetPriority this TAG, it should return nil.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, rand.Float64())
			}

			v4 := 10000
			Expect(heap.GetPriority(v4)).Should(BeEquivalentTo(math.Inf(-1)))
		})

		It("Given one fibHeaps which has a value with TAG, when GetPriority this TAG, it should return the value with this TAG.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, rand.Float64())
			}
			heap.Insert(10000, float64(10000))

			v4 := 10000
			Expect(heap.GetPriority(v4)).Should(BeEquivalentTo(10000))
			Expect(heap.Num()).Should(BeEquivalentTo(1001))
		})

		It("Given one fibHeaps which has not a value with TAG, when ExtractPriority this TAG, it should return nil.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, rand.Float64())
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			Expect(heap.ExtractPriority(1000)).Should(BeEquivalentTo(math.Inf(-1)))
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given one fibHeaps which has a value with TAG, when ExtractPriority this TAG, it should extract the value with this TAG from the heap.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			Expect(heap.ExtractPriority(999)).Should(BeEquivalentTo(999))
			Expect(heap.Num()).Should(BeEquivalentTo(999))
			_, minKey := heap.Minimum()
			Expect(minKey).Should(BeEquivalentTo(0))
		})
	})

	Context("debug test", func() {
		BeforeEach(func() {
			heap = fibheap.NewFibHeap[int]()
		})

		AfterEach(func() {
			heap = nil
		})

		It("Given one fibHeaps which some values, when call String api, it should return the internal debug string.", func() {
			Expect(heap.Stats()).Should(BeEquivalentTo("Heap is empty.\n"))
			for i := 1; i < 5; i++ {
				for j := 10; j < 15; j++ {
					demo := new(demoStruct)
					demo.data = i * j
					demo.priority = float64(i * j)
					demo.value = fmt.Sprint(demo.priority)
					heap.Insert(demo.data, demo.priority)
				}
				heap.ExtractMin()
			}

			debugMsg := "Total number: 16, Root Size: 1, Index size: 16,\n" +
				"Current min: priority(14.000000), data(14),\n" +
				"Heap detail:\n" +
				"< 14.000000 < 56.000000 28.000000 < 42.000000 > 30.000000 < 33.000000 36.000000 < 39.000000 > > 20.000000 < 22.000000 24.000000 < 26.000000 > 40.000000 < 44.000000 48.000000 < 52.000000 > > > > > \n"
			Expect(heap.Stats()).Should(BeEquivalentTo(debugMsg))
		})

		It("Given one fibHeaps which one normal and multi +inf prioritys, when call ExtractMin, it should update min value correctly.", func() {
			heap.Insert(0, 0)
			heap.Insert(1, math.Inf(1))
			heap.Insert(2, math.Inf(1))
			heap.Insert(3, math.Inf(1))

			_, priority := heap.ExtractMin()
			Expect(priority).Should(BeEquivalentTo(0))
			_, priority = heap.ExtractMin()
			Expect(priority).Should(BeEquivalentTo(math.Inf(1)))
		})
	})
})
