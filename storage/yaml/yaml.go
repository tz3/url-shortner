package yaml

import (
	"fmt"
	"github.com/tz3/url-shortner/storage"
	parser "gopkg.in/yaml.v3"
	"io"
	"log/slog"
	"net/url"
)

var Log *slog.Logger = slog.Default()

// yaml represents a URL shortener backed by YAML storage.
type yaml struct {
	repo storage.Repository
}

// row represents a YAML entry with Short and Long URL mappings.
type row struct {
	Short string `yaml:"short"` // Short is the source URL that will be redirected.
	Long  string `yaml:"long"`  // Long is the destination URL to which the source will be redirected.
}

// New creates a new YAML-backed URL shortener using the provided storage repository and input reader.
func New(repo storage.Repository, src io.Reader) (*yaml, error) {
	y := &yaml{repo: repo}

	// Read the content into a structure that we can convert it to URLs
	rows := make([]row, 0)
	dec := parser.NewDecoder(src)
	err := dec.Decode(&rows)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", storage.ErrorStorageInitFailed, err)
	}

	// Fill up the storage with the links
	for _, r := range rows {
		from, err := url.Parse(r.Short)
		if err != nil {
			Log.Debug("failed to parse from log", "url", r.Short, "err", err)
			continue
		}

		to, err := url.Parse(r.Long)
		if err != nil {
			Log.Debug("failed to parse from log", "url", r.Long, "err", err)
			continue
		}

		if err := y.repo.Set(from, to); err != nil {
			return nil, fmt.Errorf("%w: %s", storage.ErrorStorageInitFailed, err)
		}
	}

	return y, nil
}

// Get retrieves the destination URL for the given source URL from the storage.
func (y *yaml) Get(u *url.URL) (*url.URL, error) {
	return y.repo.Get(u)
}

// Set is not supported for YAML storage and always returns an error indicating read-only mode.
func (y *yaml) Set(*url.URL, *url.URL) error {
	return storage.ErrReadOnlyStorage
}
