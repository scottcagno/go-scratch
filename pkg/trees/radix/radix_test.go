package radix

import (
	"fmt"
	"testing"
)

func TestNewTree(t *testing.T) {

	t.Logf("Creating new radix tree...")
	rt := NewTree()
	if rt.Len() != 0 {
		t.Fatalf("Bad length, expected %v, got %v", 0, rt.Len())
	}

	var entries = []struct {
		Key string
		Val any
	}{
		{"http://www.example.com/api/{version}/", true},
		{"www.example.com", true},
		{"http://www.example.com/api/", true},
		{"www.example", true},
		{"www", true},
		{"http://www.example.com", true},
		{"www.example.org", true},
		{"http://", true},
		{"http://www.example.com/api/v2/", true},
		{"www.example.com/api/v2/users", true},
	}

	t.Logf("Inserting a few keys and values....")
	for _, entry := range entries {
		rt.Insert(entry.Key, entry.Val)
	}

	t.Logf("Searching by closest prefix...")
	searchKey := "http://www.example.com/api/foo/"
	foundKey, _, _ := rt.FindLongestPrefix(searchKey)
	fmt.Printf("FindLongestPrefix(%q) => %q\n", searchKey, foundKey)

	t.Logf("Walking the tree...")
	rt.Walk(
		func(k string, v interface{}) bool {
			fmt.Printf("key: %q\n", k)
			return false
		},
	)

	t.Logf("Checking out the minimum, and the maximum...")
	outMin, _, _ := rt.Min()
	fmt.Printf("min: %q\n", outMin)
	// if outMin != min {
	// 	t.Fatalf("bad minimum: %v %v", outMin, min)
	// }
	outMax, _, _ := rt.Max()
	fmt.Printf("max: %q\n", outMax)
	// if outMax != max {
	// 	t.Fatalf("bad maximum: %v %v", outMax, max)
	// }

	t.Logf("Checking the length...\n")
	if rt.Len() != len(entries) {
		t.Fatalf("Bad length, expected %v, got %v", len(entries), rt.Len())
	}

	t.Logf("Deleting %q...\n", "www.example.com")
	rt.Delete("www.example.com")

	t.Logf("Checking the length again...\n")
	if rt.Len() != len(entries)-1 {
		t.Fatalf("Bad length, expected %v, got %v", len(entries)-1, rt.Len())
	}
}
