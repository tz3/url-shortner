package tests

import (
	"github.com/tz3/url-shortner/storage"
	"net/url"
)

type TestOption func(t *ts)

type ts struct {
	r map[string]*url.URL

	// error will modify the test structure to return an error for all operations.
	err error
}

// New generates the test storage implementation
func New(opts ...TestOption) *ts {
	n := &ts{}

	for _, o := range opts {
		o(n)
	}

	return n
}

// WithError modifies the storage implementation such that any operation executed against it will return
// an error.
func WithError(err error) TestOption {
	return func(t *ts) {
		t.err = err
	}
}

// Get see storage.Repository
func (ts *ts) Get(u *url.URL) (*url.URL, error) {
	if ts.err != nil {
		return nil, ts.err
	}

	if v, ok := ts.r[u.String()]; ok {
		return v, nil
	}

	return nil, storage.ErrorNotFound
}

// Set see storage.Repository
func (ts *ts) Set(f *url.URL, t *url.URL) error {
	if ts.err != nil {
		return ts.err
	}

	ts.r[f.String()] = t
	return nil
}
