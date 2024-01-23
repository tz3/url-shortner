package test

import (
	"encoding/base64"
	"github.com/tz3/url-shortner/storage"
	"math/rand"
	"net/url"
	"testing"
)

const seed = 42

func BenchmarkStorage(b *testing.B, str storage.Repository, iter int64) {
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
		str.Set(short, long)

	}

	// Iterate through the whole list, finding them all. Actual benchmark.
	b.ResetTimer()
	for _, u := range urls {
		str.Get(u)
	}
}
