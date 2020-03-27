package util

import (
	"errors"
	"io"
)

var (
	ErrSymbolNotFound      = errors.New("Symbol not found")
	ErrInvalidBuffSizeArgs = errors.New("Invalid buffer size arguments")
)

func ReadAheadTo(reader io.Reader, buffer []byte, n int, symbol byte) (i int, err error) {
	_, err = reader.Read(buffer)
	if err != nil {
		return -1, err
	}

	for i = 0; i < n; i++ {
		if buffer[i] == symbol {
			return i, nil
		}
	}

	return -1, nil
}

func ReadBehindTo(reader io.ReadSeeker, buffer []byte, n int, symbol byte) (i int, err error) {
	if 0 >= n || n > len(buffer) {
		return -1, ErrInvalidBuffSizeArgs
	}

	_, err = reader.Seek(-int64(n), 1)
	if err != nil {
		return -1, err
	}

	_, err = reader.Read(buffer)
	if err != nil {
		return -1, err
	}

	for i := n - 1; i >= 0; i-- {
		if buffer[i] == symbol {
			return i, nil
		}
	}

	return -1, ErrSymbolNotFound
}
