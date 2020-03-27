package chunker

import (
	"bufio"
	"context"
	"errors"
	"io"
)

var (
	ErrInvalidFromOrTo = errors.New("Invalid from or to")
)

type FileChunk struct {
	Word string
	Err  error
}

func newErrFileChunk(err error) FileChunk {
	return FileChunk{Err: err}
}

func ReadInChunks(ctx context.Context, f io.ReadSeeker, from, to int64) (<-chan FileChunk, error) {
	if 0 > from || from >= to {
		return nil, ErrInvalidFromOrTo
	}

	var err error
	chunksCh := make(chan FileChunk)

	go func() {
		defer close(chunksCh)

		if ctx.Err() != nil {
			chunksCh <- newErrFileChunk(ctx.Err())
			return
		}

		_, err = f.Seek(from, 0)
		if err != nil {
			chunksCh <- newErrFileChunk(err)
			return
		}

		var readBytes int64 = from
		reader := bufio.NewReader(f)

	loop:
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				if err == io.EOF {
					break loop
				}

				select {
				case <-ctx.Done():
					break loop
				case chunksCh <- newErrFileChunk(err):
					break loop
				}
			}

			s := string(line)
			if s != "" {
				select {
				case <-ctx.Done():
					break loop
				case chunksCh <- FileChunk{Word: s}: // send chunk
				}
			}

			// +1 because reader.ReadLine skips new lines.
			readBytes += int64(len(line)) + 1
			if readBytes > to {
				break loop
			}
		}
	}()

	return chunksCh, nil
}
