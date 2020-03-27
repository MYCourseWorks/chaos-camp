package main

import (
	"os"

	"large_file_sorter/src/util"
)

func main() {
	genCount := util.GetEnvIntVarOrPanic("GEN_COUNT_MB")
	outDirName := util.GetEnvVarOrPanic("OUT_DIR")
	fileName := util.GetEnvVarOrPanic("FILE_NAME")

	outDir, err := os.Open(outDirName)
	if err != nil {
		panic(err)
	}

	newFilePath := outDir.Name() + "/" + fileName
	f, err := os.OpenFile(
		newFilePath,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, // Open/Cread and truncate
		os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = util.WriteRndStr(f, genCount)
	if err != nil {
		panic(err)
	}
}
