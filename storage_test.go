package debugserver

import (
	"reflect"
	"testing"
)

func contains(s string, l *list) bool {
	item := l.First

	for {
		if item == nil {
			return false
		}

		if item.Value == s {
			return true
		}

		item = item.Next
	}
}

func TestStoreSingleRequest(t *testing.T) {
	storage := NewStorage()
	r := Request{Body: "test"}
	key := "yo"

	storage.Add(key, r)
	records := storage.Get(key)

	if len(records) != 1 {
		t.Fatalf("Got %d records, expected: 1", len(records))
	}

	if records[0].Body != r.Body {
		t.Error("Added item is not equal to retrieved")
	}
}

func TestAddSingleRecord(t *testing.T) {
	l := &list{}
	expected := "yo"
	l.add(expected)

	actual := l.First

	if actual.Value != expected {
		t.Fatalf("Got: %s, expected: %s", actual.Value, expected)
	}

	if l.First != l.Last {
		t.Error("Last item not set")
	}
}

func TestDeleteTheOnlyRecord(t *testing.T) {
	l := &list{}
	value := "test"
	l.add(value)

	l.del(value)
	if l.First != nil {
		t.Error("Reference to first item was not updated")
	}

	if l.Last != nil {
		t.Error("Reference to last item was not updated")
	}
}

func TestDeleteMiddleRecord(t *testing.T) {
	var (
		l      = &list{}
		first  = "test1"
		target = "test2"
		last   = "test3"
	)

	l.add(first)
	l.add(target)
	l.add(last)

	l.del(target)
	if l.First.Value != first {
		t.Error("Unexpected first item")
	}

	if l.Last.Value != last {
		t.Error("Unexpected last item")
	}

	if contains(target, l) {
		t.Error("Item was not deleted")
	}
}

func TestAddMultipleRecords(t *testing.T) {
	l := &list{}

	expected := []string{"1", "2", "3", "4", "5"}

	for _, s := range expected {
		l.add(s)
	}

	var actual []string

	item := l.First

	for {
		if item == nil {
			break
		}

		actual = append(actual, item.Value)
		item = item.Next
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Got: %s, expected: %s", actual, expected)
	}

	if l.First.Value != expected[0] {
		t.Error("Unexpected first item")
	}

	if l.Last.Value != expected[len(expected)-1] {
		t.Error("Unexpected last item")
	}
}
