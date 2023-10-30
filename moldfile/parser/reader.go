package parser

import (
	"bytes"
	"fmt"
	"io"
)

type reader interface {
	ReadBytes() ([]byte, error)
	Empty() bool
	Len() int
	io.Seeker
	io.Reader
}

type readerImpl struct {
	br *bytes.Reader
}

func newReader(ior io.Reader) (reader, error) {
	b, err := io.ReadAll(ior)
	if err != nil {
		return nil, fmt.Errorf("failed to read all: %v", err)
	}

	r := &readerImpl{
		br: bytes.NewReader(b),
	}

	return r, nil
}

func (r *readerImpl) ReadBytes() ([]byte, error) {
	b, err := r.br.ReadByte()
	if err != nil {
		return nil, err
	}

	return []byte{b}, nil
}

func (r *readerImpl) Empty() bool {
	return r.br.Len() == 0
}

func (r *readerImpl) Len() int {
	return r.br.Len()
}

func (r *readerImpl) Seek(offset int64, whence int) (int64, error) {
	return r.br.Seek(offset, whence)
}

func (r *readerImpl) Read(p []byte) (n int, err error) {
	return r.br.Read(p)
}
