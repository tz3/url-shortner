package tests

import (
	"encoding/base64"
	"github.com/tz3/url-shortner/storage"
	"github.com/tz3/url-shortner/storage/memory"
	"math/rand"
	"net/url"
	"strconv"
	"testing"
)

const seed = 42 // is a constant used as the seed for random number generation

// benchmarkLengths contains different lengths for benchmarking.
var benchmarkLengths = []int64{10, 100, 1000, 100000, 5000000}

// sinkFactories is a map associating storage types with functions creating storage repositories.
var sinkFactories = map[string]func() storage.Repository{
	"hash table": func() storage.Repository { return memory.NewHashMap() },
}

// benchmark performs a benchmark for the given repository and iteration count.
// It generates random URLs, sets them in the repository, and then benchmarks the retrieval.
func benchmark(b *testing.B, repo storage.Repository, iter int64) {
	// Generate the URLs randomly. Uses Rand.Read() and Base64 URL safe encoding to generate
	// "fairly random" URLs, creating a large, unsorted array. All URLs point to the same result, as this is
	// outside the scope of the benchmark.
	urls := []*url.URL{}
	long := &url.URL{
		Host: "google.com",
		Path: "/benchmarks",
	}

	rand := rand.New(rand.NewSource(seed))
	var i int64
	for i = 0; i <= iter; i++ {

		bytes := make([]byte, 10)
		rand.Read(bytes)

		short := &url.URL{
			Host: "s3k",
			Path: base64.URLEncoding.EncodeToString(bytes),
		}
		urls = append(urls, long)
		err := repo.Set(short, long)
		if err != nil {
			return
		}
	}

	// Iterate through the whole list, finding them all. Actual benchmark.
	b.ResetTimer()
	for _, u := range urls {
		_, err := repo.Get(u)
		if err != nil {
			return
		}
	}
}

// race creates 1000 goroutines doing nothing, used for testing race conditions.
func race(repo storage.Repository) {
	for i := 0; i < 1000; i++ {
		go func() {}()
	}
}

// BenchmarkAll runs benchmarks for all storage types and lengths defined in sinkFactories and benchmarkLengths.
func BenchmarkAll(b *testing.B) {
	for n, f := range sinkFactories {
		for _, l := range benchmarkLengths {
			b.Run(n+"+"+strconv.Itoa(int(l)), func(b *testing.B) {
				benchmark(b, f(), l)
			})
		}
	}
}

// TestRaceAll tests goroutine races for all storage types defined in sinkFactories.
func TestRaceAll(t *testing.T) {
	for n, f := range sinkFactories {
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			race(f())
		})
	}
}
