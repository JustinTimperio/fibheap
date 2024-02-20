package fibheap

import (
	"container/list"
	"sync"
)

type Value[t any] interface {
	Tag() t
	Key() float64
}

type FibHeap[t any] struct {
	roots       *list.List
	index       map[interface{}]*node[t]
	treeDegrees map[uint]*list.Element
	min         *node[t]
	num         uint
	mutex       sync.Mutex
}

type node[t any] struct {
	self     *list.Element
	parent   *node[t]
	children *list.List
	marked   bool
	degree   uint
	position uint
	tag      t
	key      float64
	value    Value[t]
}