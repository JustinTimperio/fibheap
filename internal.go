package fibheap

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"math"
)

func probeTree[t any](buffer *bytes.Buffer, tree *list.List) {
	buffer.WriteString(fmt.Sprintf("< "))
	for e := tree.Front(); e != nil; e = e.Next() {
		buffer.WriteString(fmt.Sprintf("%f ", e.Value.(*node[t]).priority))
		if e.Value.(*node[t]).children.Len() != 0 {
			probeTree[t](buffer, e.Value.(*node[t]).children)
		}
	}
	buffer.WriteString(fmt.Sprintf("> "))
}

func (heap *FibHeap[t]) deleteNode(n *node[t]) {
	heap.decreaseKey(n, math.Inf(-1))
	heap.extractMin()
}

func (heap *FibHeap[t]) link(parent, child *node[t]) {
	child.marked = false
	child.parent = parent
	child.self = parent.children.PushBack(child)
	parent.degree++
}

func (heap *FibHeap[t]) resetMin() {
	heap.min = heap.roots.Front().Value.(*node[t])
	for tree := heap.min.self.Next(); tree != nil; tree = tree.Next() {
		if tree.Value.(*node[t]).priority < heap.min.priority {
			heap.min = tree.Value.(*node[t])
		}
	}
}

func (heap *FibHeap[t]) cut(n *node[t]) {
	n.parent.children.Remove(n.self)
	n.parent.degree--
	n.parent = nil
	n.marked = false
	n.self = heap.roots.PushBack(n)
}

func (heap *FibHeap[t]) cascadingCut(n *node[t]) {
	if n.parent != nil {
		if !n.marked {
			n.marked = true
		} else {
			parent := n.parent
			heap.cut(n)
			heap.cascadingCut(parent)
		}
	}
}

func (heap *FibHeap[t]) consolidate() {
	for tree := heap.roots.Front(); tree != nil; tree = tree.Next() {
		heap.treeDegrees[tree.Value.(*node[t]).position] = nil
	}

	for tree := heap.roots.Front(); tree != nil; {
		if heap.treeDegrees[tree.Value.(*node[t]).degree] == nil {
			heap.treeDegrees[tree.Value.(*node[t]).degree] = tree
			tree.Value.(*node[t]).position = tree.Value.(*node[t]).degree
			tree = tree.Next()
			continue
		}

		if heap.treeDegrees[tree.Value.(*node[t]).degree] == tree {
			tree = tree.Next()
			continue
		}

		for heap.treeDegrees[tree.Value.(*node[t]).degree] != nil {
			anotherTree := heap.treeDegrees[tree.Value.(*node[t]).degree]
			heap.treeDegrees[tree.Value.(*node[t]).degree] = nil
			if tree.Value.(*node[t]).priority <= anotherTree.Value.(*node[t]).priority {
				heap.roots.Remove(anotherTree)
				heap.link(tree.Value.(*node[t]), anotherTree.Value.(*node[t]))
			} else {
				heap.roots.Remove(tree)
				heap.link(anotherTree.Value.(*node[t]), tree.Value.(*node[t]))
				tree = anotherTree
			}
		}
		heap.treeDegrees[tree.Value.(*node[t]).degree] = tree
		tree.Value.(*node[t]).position = tree.Value.(*node[t]).degree
	}

	heap.resetMin()
}

func (heap *FibHeap[t]) insert(data t, priority float64) error {
	if math.IsInf(priority, -1) {
		return errors.New("Negative infinity priority is reserved for internal usage ")
	}

	heap.mutex.Lock()
	defer heap.mutex.Unlock()

	if _, exists := heap.index[data]; exists {
		return errors.New("Duplicate data is not allowed ")
	}

	node := new(node[t])
	node.children = list.New()
	node.data = data
	node.priority = priority

	node.self = heap.roots.PushBack(node)
	heap.index[node.data] = node
	heap.num++

	if heap.min == nil || heap.min.priority > node.priority {
		heap.min = node
	}

	return nil
}

func (heap *FibHeap[t]) extractMin() *node[t] {
	heap.mutex.Lock()
	defer heap.mutex.Unlock()

	min := heap.min

	children := heap.min.children
	if children != nil {
		for e := children.Front(); e != nil; e = e.Next() {
			e.Value.(*node[t]).parent = nil
			e.Value.(*node[t]).self = heap.roots.PushBack(e.Value.(*node[t]))
		}
	}

	heap.roots.Remove(heap.min.self)
	heap.treeDegrees[min.position] = nil
	delete(heap.index, heap.min.data)
	heap.num--

	if heap.num == 0 {
		heap.min = nil
	} else {
		heap.consolidate()
	}

	return min
}

func (heap *FibHeap[t]) decreaseKey(n *node[t], priority float64) error {
	heap.mutex.Lock()
	defer heap.mutex.Unlock()

	if priority >= n.priority {
		return errors.New("New priority is not smaller than current priority ")
	}

	n.priority = priority
	if n.parent != nil {
		parent := n.parent
		if n.priority < n.parent.priority {
			heap.cut(n)
			heap.cascadingCut(parent)
		}
	}

	if n.parent == nil && n.priority < heap.min.priority {
		heap.min = n
	}

	return nil
}

func (heap *FibHeap[t]) increaseKey(n *node[t], priority float64) error {
	heap.mutex.Lock()
	defer heap.mutex.Unlock()

	if priority <= n.priority {
		return errors.New("New priority is not larger than current priority ")
	}

	n.priority = priority

	child := n.children.Front()
	for child != nil {
		childNode := child.Value.(*node[t])
		child = child.Next()
		if childNode.priority < n.priority {
			heap.cut(childNode)
			heap.cascadingCut(n)
		}
	}

	if heap.min == n {
		heap.resetMin()
	}

	return nil
}
