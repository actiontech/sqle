package queue

// Queue is a fast unique items queue which stores positive integers up to a
// fixed bound.
type Queue struct {
	list []uint32
	set  []uint32
}

// New creates a Queue where n-1 is the maximum positive integer which can be
// stored in the queue.
func New(n int) *Queue {
	q := new(Queue)
	q.list = make([]uint32, 0, 10)
	q.set = make([]uint32, n)
	return q
}

// Empty returns true if the queue is empty
func (q *Queue) Empty() bool { return len(q.list) <= 0 }

// Has checks the queue to see if pc is in it
func (q *Queue) Has(pc uint32) bool {
	idx := q.set[pc]
	return idx < uint32(len(q.list)) && q.list[idx] == pc
}

// Clear clears the queue
func (q *Queue) Clear() {
	q.list = q.list[:0]
}

// Push adds an item to the queue
func (q *Queue) Push(pc uint32) {
	if q.Has(pc) {
		return
	}
	q.set[pc] = uint32(len(q.list))
	q.list = append(q.list, pc)
}

// Pop removes an item from the queue
func (q *Queue) Pop() uint32 {
	pc := q.list[len(q.list)-1]
	q.list = q.list[:len(q.list)-1]
	return pc
}
