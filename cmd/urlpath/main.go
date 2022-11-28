package main

import (
	"fmt"
	"strings"

	"github.com/scottcagno/go-scratch/pkg/trees/radix"
)

func insert(tree *radix.Tree, path string) {
	i := strings.IndexByte(path, '{')
	j := strings.IndexByte(path, '}')
	if i > 0 && j > 0 {
		tree.Insert(path[:i], path[i+1:j])
	}
	tree.Insert(path, false)
}

func main() {
	tree := radix.NewTree()
	insert(tree, "/")
	insert(tree, "/api/users")
	insert(tree, "/api/users/{id}")
	// insert(tree, "/")
	// insert(tree, "/")
	// insert(tree, "/")
	// insert(tree, "/")
	// tree.Insert("/", true)
	// tree.Insert("/api", true)
	// tree.Insert("/api/", true)
	// tree.Insert("/api/users", true)
	// tree.Insert("/api/users/", "id")
	// tree.Insert("/api/users/id", false)
	// tree.Insert("/api/users/id/", true)
	// tree.Insert("/api/users/id/home", true)
	// tree.Insert("/api/users/id/home/", true)

	searchKey := "/api/users/12"
	k, v, found := tree.FindLongestPrefix(searchKey)
	fmt.Printf("FindLongestPrefix(%q): k=%v, v=%v, found=%v\n", searchKey, k, v, found)

	searchKey = "/api/users/12"
	tree.WalkPath(
		searchKey, func(k string, v any) bool {
			fmt.Printf("WalkPath(%q): k=%v, v=%v\n", searchKey, k, v)
			if v == "id" {
				return true
			}
			return false
		},
	)

	searchKey = "/api/users/12"
	tree.WalkPrefix(
		searchKey, func(k string, v any) bool {
			fmt.Printf("WalkPrefix(%q): k=%v, v=%v\n", searchKey, k, v)
			return false
		},
	)

}
