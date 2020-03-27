package chunker_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"large_file_sorter/src/chunker"
	"large_file_sorter/src/util"
)

type FromAndToResults struct {
	data string
	from int64
	to   int64
	want string
	err  error
}

func runFromAndToTest(ctx context.Context, data string, tc FromAndToResults, t *testing.T) {
	r := strings.NewReader(data)
	chunks, err := chunker.ReadInChunks(ctx, r, tc.from, tc.to)
	if tc.err != nil {
		util.CheckExpectedError(err, tc.err, t)
		return
	}

	got := strings.Builder{}
	for c := range chunks {
		if c.Err != nil {
			util.CheckExpectedError(c.Err, tc.err, t)
			return
		}
		got.WriteString(c.Word)
	}

	if strings.Compare(tc.want, got.String()) != 0 {
		t.Fatalf("want=%s got=%s", tc.want, got.String())
	}
}

func Test_ReadInChunks_FromAndTo_OneSymbolWords_OffsetsAtNL(t *testing.T) {
	var data = "1\n2\n3\n4\n5"
	dLen := int64(len(data))

	cases := []FromAndToResults{
		{data, -1, 0, "", chunker.ErrInvalidFromOrTo},
		{data, 0, 0, "1", chunker.ErrInvalidFromOrTo},
		{data, 0, -5, "", chunker.ErrInvalidFromOrTo},
		{data, -1, 1, "", chunker.ErrInvalidFromOrTo},
		{data, -1, -1, "", chunker.ErrInvalidFromOrTo},

		{data, 0, 1, "1", nil},
		{data, 0, 3, "12", nil},
		{data, 0, 5, "123", nil},
		{data, 0, 7, "1234", nil},
		{data, 0, 8, "12345", nil},
		{data, 0, dLen, "12345", nil},
		{data, 0, 100, "12345", nil},

		{data, 1, 3, "2", nil},
		{data, 1, 5, "23", nil},
		{data, 1, 7, "234", nil},
		{data, 1, 8, "2345", nil},
		{data, 1, dLen, "2345", nil},
		{data, 2, 3, "2", nil},
		{data, 2, 5, "23", nil},
		{data, 2, 7, "234", nil},
		{data, 2, 8, "2345", nil},
		{data, 2, dLen, "2345", nil},

		{data, 3, 5, "3", nil},
		{data, 3, 7, "34", nil},
		{data, 3, 8, "345", nil},
		{data, 3, dLen, "345", nil},
		{data, 4, 5, "3", nil},
		{data, 4, 7, "34", nil},
		{data, 4, 8, "345", nil},
		{data, 4, dLen, "345", nil},

		{data, 5, 7, "4", nil},
		{data, 5, 8, "45", nil},
		{data, 5, dLen, "45", nil},
		{data, 6, 7, "4", nil},
		{data, 6, 8, "45", nil},
		{data, 6, dLen, "45", nil},

		{data, 7, 8, "5", nil},
		{data, 7, dLen, "5", nil},

		{data, 8, 10, "5", nil},

		{data, dLen, 10, "", nil},
		{data, 10, 11, "", nil},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("from=%d, to=%d", tc.from, tc.to)
		t.Run(name, func(t *testing.T) {
			runFromAndToTest(context.Background(), data, tc, t)
		})
	}
}

func Test_ReadInChunks_FromAndTo_LongerWords_OffsetAtAnyPlace(t *testing.T) {
	var data = "1234\n5678\n9"
	dLen := int64(len(data))

	cases := []FromAndToResults{
		{data, 0, 4, "1234", nil},
		{data, 1, 4, "234", nil},
		{data, 2, 4, "34", nil},
		{data, 3, 4, "4", nil},
		{data, 4, 4, "", chunker.ErrInvalidFromOrTo},

		{data, 0, 9, "12345678", nil},
		{data, 1, 9, "2345678", nil},
		{data, 2, 9, "345678", nil},
		{data, 3, 9, "45678", nil},
		{data, 4, 9, "5678", nil},
		{data, 5, 9, "5678", nil},
		{data, 6, 9, "678", nil},
		{data, 7, 9, "78", nil},
		{data, 8, 9, "8", nil},
		{data, 9, 9, "", chunker.ErrInvalidFromOrTo},

		{data, 0, dLen, "123456789", nil},
		{data, 1, dLen, "23456789", nil},
		{data, 2, dLen, "3456789", nil},
		{data, 3, dLen, "456789", nil},
		{data, 4, dLen, "56789", nil},
		{data, 5, dLen, "56789", nil},
		{data, 6, dLen, "6789", nil},
		{data, 7, dLen, "789", nil},
		{data, 8, dLen, "89", nil},
		{data, 9, dLen, "9", nil},
		{data, 10, dLen, "9", nil},
		{data, 11, dLen, "", chunker.ErrInvalidFromOrTo},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("from=%d, to=%d", tc.from, tc.to)
		t.Run(name, func(t *testing.T) {
			runFromAndToTest(context.Background(), data, tc, t)
		})
	}
}

func Test_ReadInChunks_FromAndTo_ContextDoneBeforeRead(t *testing.T) {
	data := "1234\n5678\n9\n10\n11"
	dLen := int64(len(data))
	cases := []FromAndToResults{
		{data, 0, dLen, "", nil},
		{data, 0, dLen, "", nil},
		{data, 0, dLen, "123456789", nil},
	}

	t.Run("Context cancel before first read", func(t *testing.T) {
		ctx2, cancel := context.WithCancel(context.Background())
		cancel()

		tc := cases[0]
		r := strings.NewReader(data)
		chunks, err := chunker.ReadInChunks(ctx2, r, tc.from, tc.to)
		if tc.err != nil {
			util.CheckExpectedError(err, tc.err, t)
			return
		}

		for c := range chunks {
			if errors.Is(c.Err, context.Canceled) {
				return
			}
			t.Fatalf("First element in chunk chanel has no Canceled error!")
		}
	})

	t.Run("Context exceed deadline before first read", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()

		time.Sleep(time.Millisecond * 5)

		tc := cases[1]
		r := strings.NewReader(data)
		chunks, err := chunker.ReadInChunks(ctx, r, tc.from, tc.to)
		if tc.err != nil {
			util.CheckExpectedError(err, tc.err, t)
			return
		}

		for c := range chunks {
			if errors.Is(c.Err, context.DeadlineExceeded) {
				return
			}
			t.Fatalf("First element in chunk chanel has no DeadlineExceeded error after timeout exceeded!")
		}
	})
}

func Test_ReadInChunks_FromAndTo_ContextCancelWhileRead(t *testing.T) {
	var data = "1\n22\n3333\n44\n555\n"
	dLen := int64(len(data))

	cases := []struct {
		stop int
		want string
	}{
		{0, ""}, {1, "1"}, {2, "122"}, {3, "1223333"}, {4, "122333344"},
		{5, "122333344555"}, {6, "122333344555"},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("stop at %d", tc.stop)
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			r := strings.NewReader(data)
			chunks, err := chunker.ReadInChunks(ctx, r, 0, dLen)
			if err != nil {
				t.Fatal(err)
				return
			}

			count := tc.stop
			got := strings.Builder{}
			for c := range chunks {
				switch {
				case count > 0:
					count--
				case count == 0:
					count--
					cancel()
					time.Sleep(time.Millisecond * 1) // give the context some time to cancel
					continue
				case errors.Is(c.Err, context.Canceled):
					break
				default:
					t.Fatalf("Next element in chanel, after close, has no \"context canceled\" error!")
				}

				got.WriteString(c.Word)
			}

			if strings.Compare(tc.want, got.String()) != 0 {
				t.Fatalf("want=%s got=%s", tc.want, got.String())
			}
		})
	}
}
