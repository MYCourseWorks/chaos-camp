package util

import (
	"bufio"
	"io"
	"math/rand"
	"time"
)

func WriteRndStr(w io.Writer, maxInMB int) error {
	writer := bufio.NewWriter(w)
	maxBytes := int64(maxInMB) * int64(MegaByte)
	var byteCount int64 = 0

	for {
		// Average len of a word in the english language is 4.7,
		// so the random generation range is [0..9], or 5 on average.
		rndStr := rndStr(seededRand.Intn(9) + 1)
		bytes, err := writer.WriteString(rndStr)
		if err != nil {
			return err
		}

		byteCount += int64(bytes)
		if byteCount > maxBytes {
			break
		}
	}

	return nil
}

const charset = "abcdefghijklmnopqrstuvwxyz"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func rndStringWithCharset(length int, charset string) string {
	b := make([]byte, length+1)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	b[length] = '\n'
	return string(b)
}

func rndStr(length int) string {
	return rndStringWithCharset(length, charset)
}
