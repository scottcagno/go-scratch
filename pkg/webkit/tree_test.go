package webkit

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func populateWords(tree *Tree) {
	tree.Insert("ex", true)
	tree.Insert("relax", true)
	tree.Insert("extract", true)
	tree.Insert("rewind", true)
	tree.Insert("subside", true)
	tree.Insert("review", true)
	tree.Insert("extinct", true)
	tree.Insert("exhale", true)
	tree.Insert("subs", true)
	tree.Insert("exhale", true)
	tree.Insert("subway", true)
	tree.Insert("sub", true)
	tree.Insert("antique", true)
	tree.Insert("re", true)
	tree.Insert("anti", true)
	tree.Insert("extra", true)
	tree.Insert("submerge", true)
	tree.Insert("ant", true)
	tree.Insert("submarine", true)
	tree.Insert("antiquate", true)
}

func makeKey(method, k string) (string, string) {
	var val string
	i, j := strings.IndexByte(k, '{'), strings.IndexByte(k, '}')
	if i > 0 && j > 0 {
		val = k[i+1 : j]
	}
	return fmt.Sprintf("%s:%s", method, k), val
}

func populatePaths(tree *Tree) {
	tree.Insert(makeKey(http.MethodGet, "/v1/entries"))
	tree.Insert(makeKey(http.MethodPost, "/v1/entries"))
	tree.Insert(makeKey(http.MethodGet, "/v1/posts"))
	tree.Insert(makeKey(http.MethodPost, "/v1/posts"))
	tree.Insert(makeKey(http.MethodGet, "/v1/posts/{date}"))
	tree.Insert(makeKey(http.MethodGet, "/v1/posts/{date}/comments"))
	tree.Insert(makeKey(http.MethodPost, "/v1/posts/{date}/comments"))
	tree.Insert(makeKey(http.MethodDelete, "/v1/posts/{date}/comments"))
	tree.Insert(makeKey(http.MethodPut, "/v1/posts/{date}/comments"))
}

func TestTree_Walk_WithWords(t *testing.T) {
	fmt.Println("Populating tree...")
	tree := NewTree()
	populateWords(tree)
	fmt.Println("Walking tree...")
	tree.Walk(
		func(k string, v any) bool {
			fmt.Printf("%s\n", k)
			return false
		},
	)
}

func TestTree_Walk_WithPaths(t *testing.T) {
	fmt.Println("Populating tree...")
	tree := NewTree()
	populatePaths(tree)

	// fmt.Println("Walking tree...")
	// tree.Walk(
	// 	func(k string, v any) bool {
	// 		fmt.Printf("%s (contains wildard=%v)\n", k, v)
	// 		return false
	// 	},
	// )

	prefix, _ := makeKey(http.MethodGet, "/v1/posts/{date}/comments")
	fmt.Printf("Walking path below %q\n", prefix)
	tree.WalkPathBelow(
		prefix,
		func(k string, v any) bool {
			fmt.Printf("%s (contains wildard=%v)\n", k, v)
			return false
		},
	)

	// prefix, _ = makeKey(http.MethodGet, "/v1/p")
	// fmt.Printf("Walking path above %q\n", prefix)
	// tree.WalkPathAbove(
	// 	prefix,
	// 	func(k string, v any) bool {
	// 		fmt.Printf("%s (contains wildard=%v)\n", k, v)
	// 		return false
	// 	},
	// )
}
