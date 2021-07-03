package connectionQueue

import (
	"container/heap"
	"echoServer/models"
	"time"
)

type ConnectionQueue []*models.Connection

func NewConnectionQueue() ConnectionQueue {
	return ConnectionQueue{}
}

func (cq ConnectionQueue) Len() int {
	return len(cq)
}

func (cq ConnectionQueue) Less(i, j int) bool {
	return cq[i].LastUpdate.Unix() > cq[j].LastUpdate.Unix()
}

func (cq ConnectionQueue) Swap(i, j int) {
	cq[i], cq[j] = cq[j], cq[i]
	cq[i].Index = i
	cq[j].Index = j
}

func (cq *ConnectionQueue) Push(x interface{}) {
	n := len(*cq)
	item := x.(*models.Connection)
	item.Index = n
	*cq = append(*cq, item)
	if n == 0 {
		heap.Init(cq)
	}
	heap.Fix(cq, n)
}

func (cq *ConnectionQueue) Pop() interface{} {
	if len(*cq) == 0 {
		return nil
	}
	old := *cq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*cq = old[0 : n-1]
	return item
}

// Update modifies the priority and value of an Item in the queue.
func (cq *ConnectionQueue) Update(item *models.Connection, lastUpdate time.Time) {
	if len(*cq) == 0 {
		return
	}
	item.LastUpdate = lastUpdate
	heap.Fix(cq, item.Index)
}
