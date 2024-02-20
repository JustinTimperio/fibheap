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
	tag   int
	key   float64
	value string
}

func (demo *demoStruct) Tag() interface{} {
	return demo.tag
}

func (demo *demoStruct) Key() float64 {
	return demo.key
}

type student struct {
	Name string
	Age  float64
}

func (s *student) Tag() interface{} {
	return s.Name
}

func (s *student) Key() float64 {
	return s.Age
}

func TestBasic(t *testing.T) {

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

	fmt.Println(heap.Num())     // 4
	fmt.Println(heap.Minimum()) // &{John 18.3}
	fmt.Println(heap.Num())     // 4

	heap.IncreaseKey(s1, 20.0)
	fmt.Println(heap.ExtractMin()) // &{Jessica 19.4}

	fmt.Println(heap.ExtractMin()) // &{John 20.0}
	fmt.Println(heap.Num())        // 2

	heap.DecreaseKey(s4, 16.5)
	fmt.Println(heap.ExtractMin()) // &{Amy 16.5}

	fmt.Println(heap.Num())            // 1
	fmt.Println(heap.ExtractValue(s2)) // &{Tom 21.0}
	fmt.Println(heap.Num())            // 0

	heap.Insert(s1, s1.Age)
	heap.Insert(s2, s2.Age)
	heap.Insert(s3, s3.Age)
	heap.Insert(s4, s4.Age)
	heap2.Insert(t1, t1.Age)
	heap2.Insert(t2, t2.Age)
	heap2.Insert(t3, t3.Age)

	heap.Union(heap2)

	fmt.Println(heap.Num()) // 7
	fmt.Println(heap.Minimum())
	for heap.Num() > 1 {
		heap.ExtractMin()
	}
	fmt.Println(heap.Num())     // 1
	fmt.Println(heap.Minimum()) // &{Ryan 28.0}
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

	Context("behaviour tests of tag/key interfaces", func() {
		BeforeEach(func() {
			heap = fibheap.NewFibHeap[int]()
		})

		AfterEach(func() {
			heap = nil
		})

		It("Given an empty fibHeap, when call Minimum api, it should return nil.", func() {
			tag, key := heap.Minimum()
			Expect(tag).Should(BeEquivalentTo(0))
			Expect(key).Should(BeEquivalentTo(math.Inf(-1)))
		})

		It("Given a empty fibHeap, when call Insert api with a negative infinity key, it should return error.", func() {
			Expect(heap.Insert(1000, math.Inf(-1))).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call Minimum api, it should return the minimum value inserted.", func() {
			min := math.Inf(1)
			for i := 0; i < 10000; i++ {
				key := rand.Float64()
				heap.Insert(i, key)
				if key < min {
					min = key
				}
			}

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			_, minKey := heap.Minimum()
			Expect(minKey).Should(BeEquivalentTo(min))
			Expect(heap.Num()).Should(BeEquivalentTo(10000))
		})

		It("Given an empty fibHeap, when call ExtractMin api, it should return nil.", func() {
			tag, _ := heap.ExtractMin()
			Expect(tag).Should(BeEquivalentTo(0))
		})

		It("Given a fibHeap inserted multiple values, when call ExtractMin api, it should extract the minimum value inserted.", func() {
			for i := 0; i < 10000; i++ {
				key := rand.Float64()
				heap.Insert(i, key)
			}

			Expect(heap.Num()).Should(BeEquivalentTo(10000))
			_, lastKey := heap.Minimum()
			for i := 0; i < 10000; i++ {
				_, key := heap.ExtractMin()
				Expect(key).Should(BeNumerically(">=", lastKey))
				Expect(heap.Num()).Should(BeEquivalentTo(9999 - i))
				lastKey = key
			}
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})

		It("Given a fibHeap inserted multiple values, when call DecreaseKey api with a non-exists value, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}

			Expect(heap.DecreaseKey(v1, float64(999))).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap with a value, when call DecreaseKey api with a negative infinity key, it should return error.", func() {
			heap.Insert(1000, float64(1000))
			Expect(heap.DecreaseKey(v1, math.Inf(-1))).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call DecreaseKey api with a larger key, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}

			Expect(heap.DecreaseKey(v2, float64(1000))).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap inserted multiple values, when call IncreaseKey api with a non-exists value, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}

			Expect(heap.IncreaseKey(v1, float64(999))).Should(HaveOccurred())
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given a fibHeap with a value, when call IncreaseKey api with a negative infinity key, it should return error.", func() {
			heap.Insert(1000, float64(1000))
			Expect(heap.IncreaseKey(v1, math.Inf(-1))).Should(HaveOccurred())
		})

		It("Given a fibHeap inserted multiple values, when call IncreaseKey api with a smaller key, it should return error.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}

			Expect(heap.IncreaseKey(v2, float64(998))).Should(HaveOccurred())
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
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				anotherHeap.Insert(demo.tag, demo.key)
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
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.Insert(demo.tag, demo.key)
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
				demo.tag = i
				demo.key = rand.Float64()
				demo.value = fmt.Sprint(demo.key)
				heap.Insert(demo.tag, demo.key)
			}

			for i := 5000; i < 10000; i++ {
				anotherdemo := new(demoStruct)
				anotherdemo.tag = i
				anotherdemo.key = rand.Float64()
				anotherdemo.value = fmt.Sprint(anotherdemo.key)
				anotherHeap.Insert(anotherdemo.tag, anotherdemo.key)
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
				_, key := heap.ExtractMin()
				Expect(key).Should(BeNumerically(">=", lastKey))
				lastKey = key
			}
			Expect(heap.Num()).Should(BeEquivalentTo(0))
		})
	})

	Context("index tests of tag/key interfaces", func() {
		BeforeEach(func() {
			heap = fibheap.NewFibHeap[int]()
			anotherHeap = fibheap.NewFibHeap[int]()
		})

		AfterEach(func() {
			heap = nil
			anotherHeap = nil
		})

		It("Given one fibHeap, when Insert values with same tag, it should return an error.", func() {
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

		It("Given two fibHeaps which both has value with same tag, when call Union, it should return an error.", func() {
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

		It("Given one fibHeaps which has not a value with TAG, when GetTag this TAG, it should return nil.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, rand.Float64())
			}

			v4 := 10000
			Expect(heap.GetTag(v4)).Should(BeEquivalentTo(math.Inf(-1)))
		})

		It("Given one fibHeaps which has a value with TAG, when GetTag this TAG, it should return the value with this TAG.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, rand.Float64())
			}
			heap.Insert(10000, float64(10000))

			v4 := 10000
			Expect(heap.GetTag(v4)).Should(BeEquivalentTo(10000))
			Expect(heap.Num()).Should(BeEquivalentTo(1001))
		})

		It("Given one fibHeaps which has not a value with TAG, when ExtractTag this TAG, it should return nil.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, rand.Float64())
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			Expect(heap.ExtractTag(1000)).Should(BeEquivalentTo(math.Inf(-1)))
			Expect(heap.Num()).Should(BeEquivalentTo(1000))
		})

		It("Given one fibHeaps which has a value with TAG, when ExtractTag this TAG, it should extract the value with this TAG from the heap.", func() {
			for i := 0; i < 1000; i++ {
				heap.Insert(i, float64(i))
			}
			Expect(heap.Num()).Should(BeEquivalentTo(1000))

			Expect(heap.ExtractTag(999)).Should(BeEquivalentTo(999))
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
					demo.tag = i * j
					demo.key = float64(i * j)
					demo.value = fmt.Sprint(demo.key)
					heap.Insert(demo.tag, demo.key)
				}
				heap.ExtractMin()
			}

			debugMsg := "Total number: 16, Root Size: 1, Index size: 16,\n" +
				"Current min: key(14.000000), tag(14), value(<nil>),\n" +
				"Heap detail:\n" +
				"< 14.000000 < 56.000000 28.000000 < 42.000000 > 30.000000 < 33.000000 36.000000 < 39.000000 > > 20.000000 < 22.000000 24.000000 < 26.000000 > 40.000000 < 44.000000 48.000000 < 52.000000 > > > > > \n"
			Expect(heap.Stats()).Should(BeEquivalentTo(debugMsg))
		})

		It("Given one fibHeaps which one normal and multi +inf keys, when call ExtractMin, it should update min value correctly.", func() {
			heap.Insert(0, 0)
			heap.Insert(1, math.Inf(1))
			heap.Insert(2, math.Inf(1))
			heap.Insert(3, math.Inf(1))

			_, key := heap.ExtractMin()
			Expect(key).Should(BeEquivalentTo(0))
			_, key = heap.ExtractMin()
			Expect(key).Should(BeEquivalentTo(math.Inf(1)))
		})
	})
})
