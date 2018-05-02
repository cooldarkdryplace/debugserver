package debugserver

import (
	"sync"
	"time"
)

type List struct {
	sync.Mutex
	First *Node
	Last  *Node
}

func (l *List) Add(s string) {
	l.Lock()
	defer l.Unlock()

	newNode := &Node{
		CreatedAt: time.Now().Unix(),
		Value:     s,
	}
	if l.First == nil {
		l.First = newNode
		l.Last = l.First
		return
	}

	l.Last.Next = newNode
	l.Last = newNode
}

func (l *List) DeleteExpired(timestamp int64) {
	l.Lock()
	defer l.Unlock()

	item := l.First
	for {
		if item == nil {
			return
		}

		if item.CreatedAt > timestamp {
			l.First = item
			return
		}

		item = item.Next
	}
}

type Node struct {
	CreatedAt int64
	Value     string
	Next      *Node
}
