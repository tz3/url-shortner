package yaml

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/tz3/url-shortner/storage"
	"github.com/tz3/url-shortner/storage/memory"
	"github.com/tz3/url-shortner/storage/tests"
	"io"
	"log/slog"
	"net/url"
	"testing"
)

func TestNewYaml(t *testing.T) {
	te := errors.New("Testing_Error")

	type urls struct {
		f, t *url.URL
		err  error
	}

	for _, tc := range []struct {
		name string

		repo storage.Repository
		in   io.Reader

		err  error
		urls []urls
	}{
		{
			name: "Success",
			repo: memory.NewHashMap(),
			in: bytes.NewBufferString(`
---
- short: //tz3/shortLink
  long: //tz3/longLink
- short: //tz3/longLink
  long: //tz3/baz
`),
			urls: []urls{
				{
					f: &url.URL{Host: "tz3", Path: "/shortLink"},
					t: &url.URL{Host: "tz3", Path: "/longLink"},
				},
			},
		},
		{
			name: "single line corrupt (tab character)",
			repo: memory.NewHashMap(),
			in: bytes.NewBufferString(`
- short: "//	/shortLink"
  long: //tz3/longLink
- short: //tz3/longLink
  long: //tz3/baz				
`),
			urls: []urls{
				{
					f:   &url.URL{Host: "	", Path: "/shortLink"},
					err: storage.ErrorNotFound,
				},
			},
		},
		{
			name: "yaml corrupt",
			repo: memory.NewHashMap(),
			in:   bytes.NewBufferString(`I'm not yaml!`),
			err:  storage.ErrorStorageInitFailed,
		},
		{
			name: "storage failure",
			repo: tests.New(tests.WithError(te)),
			in: bytes.NewBufferString(`
---
- short: //tz3/shortLink
  long: //tz3/longLink
- short: //tz3/longLink
  long: //tz3/baz
`),
			err: storage.ErrorStorageInitFailed,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			y, err := New(tc.repo, tc.in)

			// Validate the error case
			assert.ErrorIs(t, err, tc.err)

			// Validate State
			for _, u := range tc.urls {
				nu, err := y.Get(u.f)
				assert.ErrorIs(t, err, u.err)
				assert.Equal(t, u.t, nu)
			}
		})
	}
}

func TestYamlRejectWrite(t *testing.T) {
	t.Parallel()
	y, err := New(memory.NewHashMap(), bytes.NewBufferString("---"))

	assert.Nil(t, err)

	err = y.Set(
		&url.URL{Host: "tz3", Path: "/shortLink"},
		&url.URL{Host: "tz3", Path: "/longLink"},
	)

	assert.ErrorIs(t, err, storage.ErrReadOnlyStorage)
}

func TestLoggerOverride(t *testing.T) {
	exist := Log
	defer func() { Log = exist }()

	// Set up a new logger
	b := &bytes.Buffer{}
	Log = slog.New(slog.NewTextHandler(b, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Get an error condition (failed URL)
	New(memory.NewHashMap(), bytes.NewBufferString(`
- short: "//	/shortLink"
  long: //tz3/longLink
`))

	// Validate the output
	assert.Contains(t, b.String(), "invalid control character")
}
