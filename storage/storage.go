package storage

import (
	"errors"
	"net/url"
)

var (
	ErrorNotFound          = errors.New("NOT_FOUND_ERROR")
	ErrorStorageInitFailed = errors.New("FAILED_STORAGE_INIT")
)

type Storage interface {
	Get(url *url.URL) (*url.URL, error)
}
