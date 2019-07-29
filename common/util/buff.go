package util

import (
	"bytes"
)

type BuffReadCloser struct {
	*bytes.Reader
}

func NewBuffReadCloser(b []byte) *BuffReadCloser {
	return &BuffReadCloser{Reader: bytes.NewReader(b)}
}

func (b *BuffReadCloser) Close() error {
	return nil
}
