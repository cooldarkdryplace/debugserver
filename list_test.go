package debugserver

import (
	"reflect"
	"testing"
)

func TestAddSingle(t *testing.T) {
	l := &List{}
	expected := "yo"
	l.Add(expected)

	actual := l.First

	if actual.Value != expected {
		t.Fatalf("Got: %s, expected: %s", actual.Value, expected)
	}

	if l.First != l.Last {
		t.Error("Last item not set")
	}
}

func TestAddMultiple(t *testing.T) {
	l := &List{}

	expected := []string{"1", "2", "3", "4", "5"}

	for _, s := range expected {
		l.Add(s)
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
}
