package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	key   int
	value int
	left  *Node
	right *Node
}

type OrderedMap struct {
	root *Node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{} // need to implement
}

func (m *OrderedMap) Insert(key, value int) {
	if m.Contains(key) {
		return
	}
	m.root = m.insert(m.root, key, value)
	m.size++
}

func (m *OrderedMap) insert(node *Node, key, value int) *Node {
	if node == nil {
		return &Node{
			key:   key,
			value: value,
		}
	}
	if node.key > key {
		node.left = m.insert(node.left, key, value)
		return node
	}
	node.right = m.insert(node.right, key, value)
	return node
}

func (m *OrderedMap) Erase(key int) {
	var deleted bool
	m.root, deleted = deleteNode(m.root, key)
	if deleted {
		m.size--
	}
}

func deleteNode(node *Node, key int) (*Node, bool) {
	if node == nil {
		return nil, false
	}
	var deleted bool
	if key < node.key {
		left, del := deleteNode(node.left, key)
		node.left = left
		deleted = del
		return node, deleted
	}
	if key > node.key {
		right, del := deleteNode(node.right, key)
		node.right = right
		deleted = del
		return node, deleted
	}
	if node.left == nil {
		return node.right, true
	}
	if node.right == nil {
		return node.left, true
	}

	minRight := findMin(node.right)
	node.key = minRight.key
	node.value = minRight.value
	node.right, deleted = deleteNode(node.right, minRight.key)
	return node, deleted
}

func findMin(node *Node) *Node {
	for node != nil && node.left != nil {
		node = node.left
	}
	return node
}

func (m *OrderedMap) Contains(key int) bool {
	node := m.root
	for {
		if node == nil {
			return false
		}
		if node.key == key {
			return true
		}
		if node.key > key {
			node = node.left
			continue
		}
		node = node.right
	}
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	forEachInOrder(m.root, action)
}

func forEachInOrder(node *Node, action func(int, int)) {
	if node == nil {
		return
	}
	forEachInOrder(node.left, action)
	action(node.key, node.value)
	forEachInOrder(node.right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
