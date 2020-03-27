package sortverify

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"large_file_sorter/src/chunker"
	"large_file_sorter/src/util"
	wsplit "large_file_sorter/src/work_spliter"
)

var (
	ErrInvalidSort = errors.New("File is not sorted")
)

func ParallelVerify(path string, numCPUs int, maxWordLen int) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		return err
	}

	offsets, err := genFileOffsets(f, info.Size(), numCPUs, maxWordLen)
	if err != nil {
		return err
	}

	err = verifySegments(f, offsets, maxWordLen)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	splits := len(offsets)
	errCh := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	// FIXME: Debug code
	runtime.LockOSThread()

	for i := 0; i < splits-1; i++ {
		from, to := offsets[i], offsets[i+1]
		go func(ctx context.Context, path string, from, to int64) {
			var err error
			f, err := os.Open(path)
			if err != nil {
				errCh <- err
				return
			}

			err = verifyChunk(ctx, f, from, to)
			errCh <- err
		}(ctx, path, from, to)
	}

	for i := 0; i < splits-1; i++ {
		e := <-errCh
		if e != nil {
			err = e
			break
		}
	}

	close(errCh)
	cancel()
	return err
}

func genFileOffsets(reader io.ReadSeeker, size int64, numCPUs int, maxWordLen int) (offsets []int64, err error) {
	offsets, err = wsplit.SplitForParallelWork(
		reader,
		size,
		maxWordLen,
		numCPUs,
	)
	if err != nil {
		// TODO: Can't split for parallel work, fallback to linear verify then.
		return nil, err
	}

	// Rewind
	_, err = reader.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	return offsets, nil
}

func verifySegments(reader io.ReadSeeker, offsets []int64, maxWordLen int) (err error) {
	readAheadBuf := make([]byte, maxWordLen)
	readBehindBuf := make([]byte, maxWordLen)

	for i := 1; i < len(offsets)-1; i++ {
		off := offsets[i]
		prevOff := offsets[i-1]

		reader.Seek(off+1, 0)
		if err != nil {
			return err
		}
		ralni, err := util.ReadAheadTo(reader, readAheadBuf, len(readAheadBuf), '\n')
		if err != nil {
			return err
		}

		reader.Seek(off-1, 0)
		if err != nil {
			return err
		}
		rbn := maxWordLen
		if off-prevOff < int64(maxWordLen) {
			rbn = int(off - prevOff)
		}

		rblni, err := util.ReadBehindTo(reader, readBehindBuf, rbn, '\n')
		if err != nil {
			return err
		}

		fmt.Println(ralni, rblni)
		a := readAheadBuf[0:ralni]
		b := readBehindBuf[rblni:rbn]
		fmt.Println(string(a), string(b))
	}

	// Rewind
	_, err = reader.Seek(0, 0)
	if err != nil {
		return err
	}
	return nil
}

func verifyChunk(ctx context.Context, reader io.ReadSeeker, from, to int64) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	chunksCh, err := chunker.ReadInChunks(ctx, reader, from, to)
	if err != nil {
		return err
	}

	var prev, curr string
	var chunk chunker.FileChunk
	var ok bool
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case chunk, ok = <-chunksCh:
			curr = chunk.Word
		}

		switch {
		case !ok:
			break loop
		case chunk.Err != nil:
			err = chunk.Err
			break loop
		case prev != "" && strings.Compare(prev, curr) > 0:
			err = ErrInvalidSort
			break loop
		}

		prev = curr
	}

	return err
}
