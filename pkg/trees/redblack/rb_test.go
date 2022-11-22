package redblack

import (
	"fmt"
	"strings"
	"testing"
)

func Assert(t *testing.T, label string, expected, got any) {
	if expected != got {
		t.Errorf("%s: expected=%+v, got=%+v\n", label, expected, got)
	}
}

func AssertB(b *testing.B, label string, expected, got any) {
	if expected != got {
		b.Errorf("%s: expected=%+v, got=%+v\n", label, expected, got)
	}
}

type myItem struct {
	data string
}

func (e myItem) Compare(other Item) int {
	return strings.Compare(e.data, other.(myItem).data)
}

func (e myItem) Size() int {
	return len(e.data)
}

func (e myItem) String() string {
	return fmt.Sprintf("%q", e.data)
}

func TestTreeNIL(t *testing.T) {
	tree := NewTree()
	if tree.NIL != nil {
		t.Logf("tree.NIL==%#v\n", tree.NIL)
	}
}

func TestTree_Add(t *testing.T) {
	tree := NewTree()
	Assert(t, "count (before)", 0, tree.count)
	for i := 0; i < 32; i++ {
		it := myItem{
			data: fmt.Sprintf("myItem-%.3d", i),
		}
		added := tree.Add(it)
		Assert(t, "adding", true, added)
	}
	it := myItem{
		data: fmt.Sprintf("myItem-%.3d", 16),
	}
	added := tree.Add(it)
	Assert(t, "adding", true, added)
	Assert(t, "count (after)", 32, tree.count)
	for i := 0; i < 32; i++ {
		it := myItem{
			data: fmt.Sprintf("myItem-%.3d", i),
		}
		got, found := tree.Get(it)
		Assert(t, "getting", true, found)
		fmt.Println(got)
	}
	tree = nil
}

func TestTree_Scan(t *testing.T) {
	tree := NewTree()
	for i := 0; i < 32; i++ {
		tree.Add(myItem{fmt.Sprintf("myItem-%.3d", i)})
	}
	tree.Scan(
		func(it Item) bool {
			fmt.Println(it)
			return true
		},
	)
	tree = nil
}

func TestTree_Iter(t *testing.T) {
	tree := NewTree()
	for i := 0; i < 32; i++ {
		tree.Add(myItem{fmt.Sprintf("myItem-%.3d", i)})
	}
	tree.Print()
	tree = nil
}

func TestTree_Print(t *testing.T) {
	tree := NewTree()
	for i := 0; i < 18; i++ {
		tree.Add(myItem{fmt.Sprintf("myItem-%.3d", i)})
	}
	tree.Print()
}

func BenchmarkTree_Scan(b *testing.B) {
	tree := NewTree()
	for i := 0; i < 250; i++ {
		tree.Add(myItem{fmt.Sprintf("myItem-%.3d", i)})
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tree.Scan(
			func(it Item) bool {
				if it == nil {
					b.Error("got a nil myItem")
				}
				return it != nil
			},
		)
	}
	tree = nil
}

func BenchmarkTree_Iter(b *testing.B) {
	tree := NewTree()
	for i := 0; i < 250; i++ {
		tree.Add(myItem{fmt.Sprintf("myItem-%.3d", i)})
	}
	it := tree.NewIterator()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for e := it.First(); e != nil; e = it.Next() {
			if e == nil {
				b.Error("go a nil myItem")
			}
		}
	}
	tree = nil
}
