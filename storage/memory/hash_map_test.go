package memory

import (
	"github.com/stretchr/testify/assert"
	"github.com/tz3/url-shortner/storage/memory/test"
	"net/url"
	"testing"
)

func BenchmarkHashTable10(b *testing.B) {
	test.BenchmarkStorage(b, NewHashMap(), 10)
}

func BenchmarkHashTable100(b *testing.B) {
	test.BenchmarkStorage(b, NewHashMap(), 100)
}

func BenchmarkHashTable100000(b *testing.B) {
	test.BenchmarkStorage(b, NewHashMap(), 100000)
}

func TestNewHashMap(t *testing.T) {
	t.Parallel()

	hm := NewHashMap()
	hm.Set(&url.URL{
		Host: "short1",
	}, &url.URL{
		Host: "long1.com",
	})

	res, err := hm.Get(&url.URL{
		Host: "short1",
	})

	assert.Nil(t, err)
	assert.Equal(t, &url.URL{
		Host: "long1.com",
	}, res)
}
