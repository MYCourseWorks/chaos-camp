package util_test

import (
	"errors"
	"fmt"
	"io"
	"large_file_sorter/src/util"
	"strings"
	"testing"
)

func seekReaderTo(r io.ReadSeeker, to int64, t *testing.T) {
	_, err := r.Seek(to, 0)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ReadBehindTo(t *testing.T) {
	data := "12\n345\n6789\n0"
	largeBuff := make([]byte, 300)
	// exactBuff := make([]byte, 4)
	smallerBuff := make([]byte, 3)

	cases := []struct {
		offset int64
		buffer []byte
		n      int
		want   string
		err    error
	}{
		{offset: 7, buffer: largeBuff, n: -1, want: "", err: util.ErrInvalidBuffSizeArgs},
		{offset: 7, buffer: largeBuff, n: 0, want: "", err: util.ErrInvalidBuffSizeArgs},
		{offset: 7, buffer: smallerBuff, n: 5, want: "", err: util.ErrInvalidBuffSizeArgs},

		{offset: 6, buffer: largeBuff, n: 2, want: "345", err: util.ErrSymbolNotFound},
		{offset: 6, buffer: largeBuff, n: 3, want: "345", err: util.ErrSymbolNotFound},
		{offset: 6, buffer: largeBuff, n: 4, want: "345", err: nil},
		{offset: 6, buffer: largeBuff, n: 5, want: "345", err: nil},
		{offset: 6, buffer: largeBuff, n: 6, want: "345", err: nil},
		{
			offset: 6, buffer: largeBuff, n: 7, want: "345",
			err: errors.New("strings.Reader.Seek: negative position"),
		},

		{offset: 7, buffer: largeBuff, n: 1, want: "", err: nil},
		{offset: 7, buffer: largeBuff, n: 2, want: "", err: nil},
		{offset: 7, buffer: largeBuff, n: 3, want: "", err: nil},
		{offset: 7, buffer: largeBuff, n: 4, want: "", err: nil},
		{offset: 7, buffer: largeBuff, n: 5, want: "", err: nil},
		{offset: 7, buffer: largeBuff, n: 6, want: "", err: nil},
		{offset: 7, buffer: largeBuff, n: 7, want: "", err: nil},
		{
			offset: 7, buffer: largeBuff, n: 8, want: "",
			err: errors.New("strings.Reader.Seek: negative position"),
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("nz oshte")
		t.Run(name, func(t *testing.T) {
			reader := strings.NewReader(data)
			seekReaderTo(reader, tc.offset, t)
			iOfNL, err := util.ReadBehindTo(reader, tc.buffer, tc.n, '\n')
			if tc.err != nil {
				if err == nil || tc.err.Error() != err.Error() {
					t.Fatalf("Not the wanted error. want=%v, got=%v", tc.err, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
				return
			}

			got := string(tc.buffer[iOfNL+1 : tc.n])
			if strings.Compare(got, tc.want) != 0 {
				t.Fatalf("got=%s, want=%s", got, tc.want)
				return
			}
		})
	}
}
