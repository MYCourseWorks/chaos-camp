package main

import (
	"fmt"

	"large_file_sorter/src/constants"
	sortverify "large_file_sorter/src/sort_verifier"
)

func main() {
	err := sortverify.ParallelVerify("./out/test", 2, constants.MAX_WORD_SIZE)
	fmt.Println(err)
}
