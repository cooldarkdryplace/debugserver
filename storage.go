package debugserver

import (
	"sync"
	"time"
)

type Storage struct {
	sync.Mutex
	buckets map[string][]Request
	records *list
}

func NewStorage() *Storage {
	s := &Storage{
		buckets: make(map[string][]Request),
		records: &list{},
	}

	go func() {
		for t := range time.Tick(5 * time.Second) {
			s.Expire(t.Unix() - 20)
		}
	}()

	return s
}

func (s *Storage) Get(id string) []Request {
	s.Lock()
	defer s.Unlock()
	return s.buckets[id]
}

func (s *Storage) Add(id string, r Request) {
	s.Lock()
	s.records.add(id)
	s.buckets[id] = append(s.buckets[id], r)
	s.Unlock()
}

func (s *Storage) Del(id string) {
	s.Lock()
	defer s.Unlock()
	delete(s.buckets, id)
	s.records.del(id)
}

func (s *Storage) Expire(timestamp int64) {
	s.Lock()
	defer s.Unlock()

	item := s.records.First
	for {
		if item == nil {
			s.records.First = nil
			s.records.Last = nil
			return
		}

		if item.CreatedAt > timestamp {
			s.records.First = item
			return
		}

		delete(s.buckets, item.Value)
		item = item.Next
	}
}

type node struct {
	CreatedAt int64
	Value     string
	Next      *node
}

type list struct {
	First *node
	Last  *node
}

func (l *list) add(s string) {
	newNode := &node{
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

func (l *list) del(s string) {
	item := l.First
	for {
		if item == nil || item.Next == nil {
			return
		}

		if item.Next.Value == s {
			item.Next = item.Next.Next
			return
		}

		item = item.Next
	}
}
