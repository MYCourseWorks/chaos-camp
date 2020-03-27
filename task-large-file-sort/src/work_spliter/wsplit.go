package wsplit

import (
	"errors"
	"fmt"
	"io"

	"large_file_sorter/src/util"
)

var (
	ErrInvalidEncoding = errors.New("Invalid input file encoding")
	ErrNumCPUsBelowOne = errors.New("numCPUs must be > 1")
)

// SplitForParallelWork reads a byte stream and creates an array of offsets,
// such that the stream is evenly split for the provided numCPU. Those
// offsets can then be used to proccess the file in parallel.
func SplitForParallelWork(r io.ReadSeeker, size int64, maxWordLen, numCPUs int) (offsets []int64, err error) {
	if numCPUs < 1 {
		return nil, ErrNumCPUsBelowOne
	}
	if size/int64(numCPUs) < 2 {
		return nil, fmt.Errorf("inSize %d too small to split in %d parts", size, numCPUs)
	}

	offsets = make([]int64, 0)
	onePart := size / int64(numCPUs)
	readAheadBuf := make([]byte, maxWordLen)
	currSplits := int64(1)
	nlOffsets := int64(0)

	offsets = append(offsets, 0)
	for currSplits < int64(numCPUs) {
		currOffset := currSplits*onePart + nlOffsets
		_, err = r.Seek(currOffset, 0)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		var indexOfNL int
		indexOfNL, err = util.ReadAheadTo(r, readAheadBuf, maxWordLen, '\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// The situation where there is no new line in the readAheadBuf is considered an encoding error
		if indexOfNL == -1 {
			if currOffset+int64(maxWordLen) > size {
				// Not an error if we go past the EOF
				break
			}
			return nil, ErrInvalidEncoding
		}

		// Add the index of new next new line and append to the offset
		currOffset += int64(indexOfNL)
		if currOffset >= size-1 {
			break
		}
		offsets = append(offsets, currOffset)

		nlOffsets += int64(indexOfNL)
		currSplits++
	}
	offsets = append(offsets, size-1)

	return offsets, nil
}
