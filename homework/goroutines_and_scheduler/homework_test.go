package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap          []Task
	startPriority map[int]Task
	idxStore      map[int]int
}

func NewScheduler() Scheduler {
	return Scheduler{
		heap:          make([]Task, 0),
		startPriority: make(map[int]Task),
		idxStore:      make(map[int]int),
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.heap = append(s.heap, task)

	s.startPriority[task.Identifier] = task

	idx := len(s.heap) - 1
	s.idxStore[task.Identifier] = idx

	s.siftUp(idx)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	idx, exists := s.idxStore[taskID]
	if !exists {
		return
	}

	oldPriority := s.heap[idx].Priority
	if newPriority == oldPriority {
		return
	}

	s.heap[idx].Priority = newPriority

	if newPriority > oldPriority {
		s.siftUp(idx)
		return
	}
	s.siftDown(idx)
}

func (s *Scheduler) GetTask() Task {
	if len(s.heap) == 0 {
		return Task{}
	}

	s.swap(0, len(s.heap)-1)
	task := s.heap[len(s.heap)-1]
	delete(s.idxStore, task.Identifier)
	s.heap = s.heap[:len(s.heap)-1]

	if len(s.heap) > 0 {
		s.siftDown(0)
	}

	return s.startPriority[task.Identifier]
}

func (s *Scheduler) siftUp(elem int) {
	if elem == 0 {
		return
	}

	parentElem := parent(elem)
	if s.heap[elem].Priority <= s.heap[parentElem].Priority {
		return
	}

	s.swap(elem, parentElem)
	s.siftUp(parentElem)
}

func (s *Scheduler) siftDown(elem int) {
	left := leftChild(elem)
	right := rightChild(elem)
	biggest := elem

	if left < len(s.heap) && s.heap[left].Priority > s.heap[biggest].Priority {
		biggest = left
	}
	if right < len(s.heap) && s.heap[right].Priority > s.heap[biggest].Priority {
		biggest = right
	}

	if biggest != elem {
		s.swap(elem, biggest)
		s.siftDown(biggest)
	}
}

func (s *Scheduler) swap(i, j int) {
	s.heap[i], s.heap[j] = s.heap[j], s.heap[i]
	s.idxStore[s.heap[i].Identifier] = i
	s.idxStore[s.heap[j].Identifier] = j
}

func parent(elem int) int     { return (elem - 1) / 2 }
func leftChild(elem int) int  { return 2*elem + 1 }
func rightChild(elem int) int { return 2*elem + 2 }

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
