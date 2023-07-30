package priorityqueue

import(
    "container/heap"
)

type FileHash struct {
    Priority    int     // frequency
    Value       string  // filename
    Index       int     
}

type PQueue []*FileHash

func (pq PQueue) Len() int { return len(pq) }

func (pq PQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq PQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PQueue) Push(x any) {
	n := len(*pq)
	item := x.(*FileHash)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PQueue) update(item *FileHash, Value string, Priority int) {
	item.Value = Value
	item.Priority = Priority
	heap.Fix(pq, item.Index)
}
