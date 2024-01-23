package memory

import (
	"github.com/tz3/url-shortner/storage"
	"net/url"
	"sync"
)

// HashMap represents a data structure for mapping short strings to long URLs.
type HashMap struct {
	shortToLongMap map[string]*url.URL // Map to store short-to-long URL mappings.
	mux            sync.Mutex          // Mutex for concurrent access control.
}

// NewHashMap creates a new instance of HashMap with an initialized map and mutex.
func NewHashMap() *HashMap {
	return &HashMap{
		shortToLongMap: make(map[string]*url.URL),
		mux:            sync.Mutex{},
	}
}

// Get retrieves the long URL associated with a given short URL.
func (hm *HashMap) Get(short *url.URL) (*url.URL, error) {
	hm.mux.Lock()
	defer hm.mux.Unlock()

	v, ok := hm.shortToLongMap[short.String()]
	if !ok {
		return nil, storage.ErrorNotFound
	}
	return v, nil
}

// Set associates a short URL with a long URL in the HashMap.
func (hm *HashMap) Set(short *url.URL, long *url.URL) error {
	hm.mux.Lock()
	defer hm.mux.Unlock()

	if _, ok := hm.shortToLongMap[short.String()]; ok {
		return storage.ErrorDupKey
	}

	// Set the short-to-long URL mapping.
	hm.shortToLongMap[short.String()] = long
	return nil
}
