package wsplit_test

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	wsplit "large_file_sorter/src/work_spliter"
)

type SplitForParallelWorkResult struct {
	name    string
	data    string
	numCPUs int
	want    []int64
}

func verifySplitForParallelWork(got []int64, tc SplitForParallelWorkResult, t *testing.T) {
	if !reflect.DeepEqual(got, tc.want) {
		t.Fatalf("got %v want %v", got, tc.want)
	}

	if len(got) > tc.numCPUs+1 {
		t.Fatalf("Split in too many parts, len of got = %v numCPUs = %v", len(got), tc.numCPUs+1)
	}
}

func callSplitForParallelWork(tc SplitForParallelWorkResult, t *testing.T) []int64 {
	var f = strings.NewReader(tc.data)
	size := f.Size()
	got, err := wsplit.SplitForParallelWork(f, size, 500, tc.numCPUs)
	if err != nil {
		t.Fatal(err)
	}

	return got
}

func genSplitForParallelWorkResults(data string, expectedTable [][]int64) []SplitForParallelWorkResult {
	var cases = make([]SplitForParallelWorkResult, len(expectedTable))

	for i := 0; i < len(expectedTable); i++ {
		cases[i] = SplitForParallelWorkResult{
			name:    strconv.Itoa(i+1) + " CPUs",
			data:    data,
			numCPUs: i + 1,
			want:    expectedTable[i],
		}
	}

	return cases
}

func Test_SplitForParallelWork_WordsWithLen1(t *testing.T) {
	var data = "1\n2\n3\n4\n5\n6\n7\n8\n9\n"
	var dLen int64 = int64(len(data)) - 1

	expected := [][]int64{
		{0, dLen}, // 1
		{0, 9, dLen},
		{0, 7, 13, dLen},
		{0, 5, 9, 13, dLen},
		{0, 3, 7, 11, 15, dLen},
		{0, 3, 7, 11, 15, dLen},
		{0, 3, 5, 7, 9, 11, 13, dLen},
		{0, 3, 5, 7, 9, 11, 13, 15, dLen},
		{0, 3, 5, 7, 9, 11, 13, 15, dLen}, // 9
	}

	var cases = genSplitForParallelWorkResults(data, expected)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := callSplitForParallelWork(tc, t)
			verifySplitForParallelWork(got, tc, t)
		})
	}
}

func Test_SplitForParallelWork_AverageWordLen(t *testing.T) {
	var data = "01234\n678\n101112"

	expected := [][]int64{
		{0, 15}, // 1
		{0, 9, 15},
		{0, 5, 15},
		{0, 5, 9, 15},
		{0, 5, 9, 15},
		{0, 5, 9, 15},
		{0, 5, 9, 15},
		{0, 5, 9, 15}, // 8
	}

	var cases = genSplitForParallelWorkResults(data, expected)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := callSplitForParallelWork(tc, t)
			verifySplitForParallelWork(got, tc, t)
		})
	}
}

func Test_SplitForParallelWork_LongerWords(t *testing.T) {
	var data = "01234\n678\n56723123216213512312351233\n1273\n"

	expected := [][]int64{
		{0, 41}, // 1
		{0, 36, 41},
		{0, 36, 41},
		{0, 36, 41},
		{0, 9, 36, 41}, // 5
		{0, 9, 36, 41},
		{0, 9, 36, 41},
		{0, 5, 36, 41},
		{0, 5, 9, 36, 41}, // 9
		{0, 5, 9, 36, 41},
		{0, 5, 9, 36, 41},
		{0, 5, 9, 36, 41},
		{0, 5, 9, 36, 41},
		{0, 5, 9, 36, 41},
		{0, 5, 9, 36, 41},
		{0, 5, 9, 36, 41}, // 16
	}

	var cases = genSplitForParallelWorkResults(data, expected)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := callSplitForParallelWork(tc, t)
			verifySplitForParallelWork(got, tc, t)
		})
	}
}

func Test_SplitForParallelWork_Errors(t *testing.T) {
	var data = "01234\n678\n101112"
	var err error

	var f = strings.NewReader(data)
	size := f.Size()

	numCPUs := 0
	_, err = wsplit.SplitForParallelWork(f, size, 500, numCPUs)
	if err != wsplit.ErrNumCPUsBelowOne {
		t.Fatalf("unexpected error = %v", err)
	}

	numCPUs = len(data)/2 + 1
	_, err = wsplit.SplitForParallelWork(f, size, 500, numCPUs)
	if err == nil &&
		err.Error() != fmt.Sprintf("inSize %d too small to split in %d parts", size, numCPUs) {
		t.Fatalf("invalid error = %v", err)
	}

	numCPUs = len(data) / 2
	_, err = wsplit.SplitForParallelWork(f, size, 500, numCPUs)
	if err != nil {
		t.Fatalf("Should work when numCPUs is len(data)/2, err = %v", err)
	}
}
