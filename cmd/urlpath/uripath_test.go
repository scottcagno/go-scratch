package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Jonny-Burkholder/streaming-example/pkg/netkit"
	"github.com/scottcagno/go-scratch/pkg/trees/radix"
)

var urlDataSet = []struct {
	Pattern     string
	RequestURI  string
	ShouldMatch bool
}{
	{
		Pattern:     "/api/users",
		RequestURI:  "/api/users",
		ShouldMatch: true,
	},
	{
		Pattern:     "/api/users",
		RequestURI:  "/api/users/",
		ShouldMatch: false,
	},
	{
		Pattern:     "/api/users/{id}",
		RequestURI:  "/api/users/123890",
		ShouldMatch: true,
	},
	{
		Pattern:     "/api/users/{id}",
		RequestURI:  "/api/users/foobar",
		ShouldMatch: true,
	},
}

func insert(tree *radix.Tree, path string) {
	i := strings.IndexByte(path, '{')
	j := strings.IndexByte(path, '}')
	if i > 0 && j > 0 {
		tree.Insert(path[:i], path[i+1:j])
	}
	tree.Insert(path, false)
}

func assert(b *testing.B, msg, uri, key string, val any, found, shouldMatch bool) {
	if found == shouldMatch || val == shouldMatch {
		return
	}
	if val == "id" && strings.HasPrefix(uri, key) {
		return
	}
	// if val == "id" && strings.HasPrefix(uri, key) && found == shouldMatch {
	// 	return
	// }
	b.Errorf("[%s] uri=%q, key=%q, val=%v, found=%v, shouldMatch=%v\n", msg, uri, key, val, found, shouldMatch)
}

func BenchmarkURLPath_Radix_FindLongestPrefix(b *testing.B) {
	b.ReportAllocs()
	tree := radix.NewTree()
	for _, data := range urlDataSet {
		insert(tree, data.Pattern)
	}
	var key string
	var val any
	var found bool
	for i := 0; i < b.N; i++ {
		for _, data := range urlDataSet {
			key, val, found = tree.FindLongestPrefix(data.RequestURI)
			assert(b, "FindLongestPrefix", data.RequestURI, key, val, found, data.ShouldMatch)
			// if key != data.RequestURI || found != data.ShouldMatch {
			// 	b.Errorf("FindLongestPrefix(%q): k=%v, v=%v, found=%v\n", data.RequestURI, key, val, found)
			// }
		}
	}
}

func BenchmarkURLPath_Radix_WalkPath(b *testing.B) {
	b.ReportAllocs()
	tree := radix.NewTree()
	for _, data := range urlDataSet {
		insert(tree, data.Pattern)
	}
	var key string
	var val any
	var found bool
	for i := 0; i < b.N; i++ {
		for _, data := range urlDataSet {
			tree.WalkPath(
				data.RequestURI, func(k string, v any) bool {
					if k == data.RequestURI || v == "id" {
						key = k
						val = v
						found = true
						return true
					}
					return false
				},
			)
			assert(b, "WalkPathBelow", data.RequestURI, key, val, found, data.ShouldMatch)
		}
	}
}

func BenchmarkURLPath_Radix_WalkPrefix(b *testing.B) {
	b.ReportAllocs()
	tree := radix.NewTree()
	for _, data := range urlDataSet {
		insert(tree, data.Pattern)
	}
	var key string
	var val any
	var found bool
	for i := 0; i < b.N; i++ {
		for _, data := range urlDataSet {
			tree.WalkPrefix(
				data.RequestURI, func(k string, v any) bool {
					if v == "id" {
						key = k
						val = v
						found = true
						return true
					}
					return false
				},
			)
			assert(b, "WalkPathAbove", data.RequestURI, key, val, found, data.ShouldMatch)
		}
	}
}

func BenchmarkURLPath_RawStringMatch(b *testing.B) {
	b.ReportAllocs()
	var err error
	var foundMatch bool
	for i := 0; i < b.N; i++ {
		for _, data := range urlDataSet {
			foundMatch, err = filepath.Match(data.Pattern, data.RequestURI)
			if foundMatch && data.ShouldMatch == foundMatch {
				b.Errorf("[RawStringMatch] %s\n", err)
			}
		}
	}
}

func BenchmarkURLPath_RegexMatch(b *testing.B) {
	b.ReportAllocs()
	var err error
	var foundMatch bool
	for i := 0; i < b.N; i++ {
		for _, data := range urlDataSet {
			foundMatch, err = regexp.Match(`\/(\{)?\w+(\})?`, []byte(data.RequestURI))
			if err != nil || foundMatch != data.ShouldMatch {
				b.Errorf(
					"[regexp.Match] foundMatch=%v, shouldMatch=%v, pattern=%q, request=%q, error=%s\n",
					foundMatch, data.ShouldMatch, data.Pattern, data.RequestURI, err,
				)
			}
		}
	}
}

func getPathParam(uri string) (string, bool) {
	if uri[len(uri)-1] == '/' {
		uri = uri[:len(uri)-1]
	}
	i := len(uri) - 1
	for i >= 0 && uri[i] != '/' {
		i--
	}
	param := uri[i+1:]
	return param, param != "" && param != uri[(len(uri)-1)/2:]
}

func getPathID(uri string) int {
	param, found := getPathParam(uri)
	if !found {
		return -1
	}
	id, err := strconv.ParseInt(param, 10, 0)
	if err != nil {
		return -1
	}
	return int(id)
}

func Benchmark_StandardLibraryHTTP(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			http.RedirectHandler("/api", http.StatusTemporaryRedirect)
			return
		},
	)
	mux.HandleFunc(
		"/api", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "API ROOT")
			return
		},
	)
	mux.HandleFunc(
		"/api/users", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "USERS LIST ROOT")
			return
		},
	)
	mux.HandleFunc(
		"/api/users/", func(w http.ResponseWriter, r *http.Request) {
			id := getPathID(r.URL.Path)
			if id < 0 {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "USERS %d ROOT", id)
			return
		},
	)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	testURLs := []struct {
		URL      string
		Response string
	}{
		{"/", ""},
		{"/api", "API ROOT"},
		{"/api/users", "USERS LIST ROOT"},
		{"/api/users/1", "USERS 1 ROOT"},
		{"/api/users/16", "USERS 16 ROOT"},
		{"/api/users/32", "USERS 32 ROOT"},
		{"/api/users/256", "USERS 256 ROOT"},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, testURL := range testURLs {
			res, err := http.Get(ts.URL + testURL.URL)
			if err != nil {
				b.Errorf("%s", err)
				log.Fatal(err)
			}
			resp, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				b.Errorf("%s", err)
				log.Fatal(err)
			}
			if string(resp) != testURL.Response {
				b.Errorf("bad response: got=%q, wanted=%q\n", string(resp), testURL.Response)
			}
		}
	}
}

func Benchmark_NetKitHTTP(b *testing.B) {
	mux := netkit.NewRouter(nil)
	mux.HandleFunc(
		http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
			http.RedirectHandler("/api", http.StatusTemporaryRedirect)
			return
		},
	)
	mux.HandleFunc(
		http.MethodGet,
		"/api", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "API ROOT")
			return
		},
	)
	mux.HandleFunc(
		http.MethodGet,
		"/api/users", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "USERS LIST ROOT")
			return
		},
	)
	mux.HandleFunc(
		http.MethodGet,
		"/api/users/", func(w http.ResponseWriter, r *http.Request) {
			id := getPathID(r.URL.Path)
			if id < 0 {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "USERS %d ROOT", id)
			return
		},
	)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	testURLs := []struct {
		URL      string
		Response string
	}{
		{"/", ""},
		{"/api", "API ROOT"},
		{"/api/users", "USERS LIST ROOT"},
		{"/api/users/1", "USERS 1 ROOT"},
		{"/api/users/16", "USERS 16 ROOT"},
		{"/api/users/32", "USERS 32 ROOT"},
		{"/api/users/256", "USERS 256 ROOT"},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, testURL := range testURLs {
			res, err := http.Get(ts.URL + testURL.URL)
			if err != nil {
				b.Errorf("%s", err)
				log.Fatal(err)
			}
			resp, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				b.Errorf("%s", err)
				log.Fatal(err)
			}
			if string(resp) != testURL.Response {
				b.Errorf("bad response: got=%q, wanted=%q\n", string(resp), testURL.Response)
			}
		}
	}
}
