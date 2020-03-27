package rindexdb

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/homework/hw4/settings"
)

type shard struct {
	FileName string
	File     *os.File
}

// IReverseIndexDB comment
type IReverseIndexDB interface {
	StoreDocument(url string, data string) error
}

// ReverseIndexDB comment
type ReverseIndexDB struct {
	shards              []shard
	maxShardSizeInBytes int64
	nextShard           int
	mux                 sync.Mutex
}

// reverseIndex comment
type reverseIndex map[int]struct {
	Line int
	Row  int
}

// scrapedDocument comment
type scrapedDocument struct {
	ID   int
	URL  string
	Data string
}

func (sd *scrapedDocument) String() string {
	return fmt.Sprintf("%d%s%s%s%s", sd.ID, UniqueDataSepSeq, sd.URL, UniqueDataSepSeq, sd.Data)
}

var (
	// ErrFileIsNotDirecotry commnet
	ErrFileIsNotDirecotry error = errors.New("provided file is not a directory")
	// ErrShardsParamInvalid comment
	ErrShardsParamInvalid error = errors.New("invalid maxShards parameter")
	// ErrShardsMaxSizeIsToLow comment
	ErrShardsMaxSizeIsToLow error = errors.New("shardsMaxSize is to low")
	// ErrOutOfSpace comment
	ErrOutOfSpace error = errors.New("db is out of space")

	// Privates :
	autoID int = 0
)

const (
	// UniqueDataSepSeq comment
	UniqueDataSepSeq = " [*sep] "
	// UniqueEOFSeq comment
	UniqueEOFSeq      = "[*end]\n"
	maxPossibleShards = 50
)

var ridb *ReverseIndexDB

func init() {
	var err error
	st := settings.GetSettings()

	ridb, err = newReverseIndexDB(st.MaxResultFiles, st.MaxResultFileSizeInBytes, "./db")
	if err != nil {
		panic(err)
	}
}

// GetInstance comment
func GetInstance() IReverseIndexDB {
	return ridb
}

func newReverseIndexDB(shardsCount int, shardsMaxSize int64, directory string) (*ReverseIndexDB, error) {
	if shardsCount < 1 || shardsCount > maxPossibleShards {
		fmt.Fprintf(os.Stderr, "Error: %s", ErrShardsParamInvalid)
		return nil, ErrFileIsNotDirecotry
	}

	if shardsMaxSize < 200 {
		fmt.Fprintf(os.Stderr, "Error: %s", ErrShardsMaxSizeIsToLow)
		return nil, ErrShardsMaxSizeIsToLow
	}

	info, err := os.Stat(directory)
	if err != nil {
		return nil, err
	}
	if info.IsDir() == false {
		fmt.Fprintf(os.Stderr, "Error: %s", ErrFileIsNotDirecotry)
		return nil, ErrFileIsNotDirecotry
	}

	db := &ReverseIndexDB{
		shards:              make([]shard, shardsCount),
		maxShardSizeInBytes: shardsMaxSize,
	}

	for i := 0; i < shardsCount; i++ {
		name := fmt.Sprintf("%s_%d.ridb", settings.GetSettings().BaseResultFileName, i)
		file, err := os.OpenFile(path.Join(directory, name), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			return nil, err
		}

		err = file.Truncate(0) // clean file
		if err != nil {
			return nil, err
		}

		shard := shard{
			FileName: name,
			File:     file,
		}

		db.shards[i] = shard
	}

	return db, nil
}

// StoreDocument comment
func (db *ReverseIndexDB) StoreDocument(url string, data string) error {
	// We want every dataline to end with the uniqueEOFSeq
	data = fmt.Sprintf("%s%s", strings.TrimSpace(data), UniqueEOFSeq)

	db.mux.Lock()
	defer db.mux.Unlock()

	// NOTE: The below code should be resonantly fast, so we can lock safely here.

	autoID++
	doc := &scrapedDocument{ID: autoID, URL: url, Data: data}
	docStr := doc.String()

	s := db.shards[db.nextShard]
	db.nextShard = (db.nextShard + 1) % len(db.shards) // even if the below fails, we can move the next shard index

	info, err := s.File.Stat()
	if err != nil {
		return err
	}

	if info.Size() >= db.maxShardSizeInBytes {
		// If one of the shards is overflowing, then probably all of them are.
		// That won't be the case if some of the write operations fail often..
		return ErrOutOfSpace
	}

	n, err := s.File.WriteString(docStr)
	if err != nil {
		return err
	}
	if n != len(docStr) {
		return fmt.Errorf("Wrote %d bytes, but Wanted To Write %d bytes", n, len(docStr))
	}

	// Index document

	return nil
}
