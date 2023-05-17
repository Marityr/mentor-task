package src

import (
	"errors"
	"io"
	"strings"
)

type Reader interface {
	Read(p []byte) (int, error)
	ReadAll(bufSize int) (string, error)
	BytesRead() int64
}

type CountingToLowerReaderImpl struct {
	Reader         io.Reader
	TotalBytesRead int64
}

func toLower(s []byte, n int) {
	for i := 0; i < n; i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		s[i] = c
	}
}

func (cr *CountingToLowerReaderImpl) Read(p []byte) (int, error) {
	n, err := cr.Reader.Read(p)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, err
		}
	}

	toLower(p, n)

	cr.TotalBytesRead += int64(n)

	return n, nil
}

func (cr *CountingToLowerReaderImpl) ReadAll(bufSize int) (string, error) {
	buffer := make([]byte, bufSize)
	var data strings.Builder

	for {
		n, err := cr.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				data.Grow(n)
				if _, err = data.Write(buffer[:n]); err != nil {
					return data.String(), err
				}
				return data.String(), err
			}
			return data.String(), err
		}
		data.Grow(n)
		if _, err := data.Write(buffer[:n]); err != nil {
			return data.String(), err
		}
	}
}

func (cr *CountingToLowerReaderImpl) BytesRead() int64 {
	if cr.TotalBytesRead == 0 {
		return 0
	}
	return cr.TotalBytesRead
}

func NewCountingReader(r io.Reader) *CountingToLowerReaderImpl {
	return &CountingToLowerReaderImpl{
		Reader: r,
	}
}
