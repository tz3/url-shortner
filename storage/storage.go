package storage

import (
	"errors"
	"net/url"
)

var (
	ErrorNotFound          = errors.New("NOT_FOUND_ERROR")
	ErrorDupKey            = errors.New("FAILED_DUPLICATED_ID")
	ErrorStorageInitFailed = errors.New("FAILED_STORAGE_INIT")
	ErrReadOnlyStorage     = errors.New("FAILED_READ_ONLY_STORAGE_TYPE")
)

type Repository interface {
	Get(url *url.URL) (*url.URL, error)
	Set(short *url.URL, long *url.URL) error
}
