package memory

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

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
