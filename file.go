package main

import (
	"os"
	"regexp"
	"time"
)

type File struct {
	buf           []byte
	dateReg       *regexp.Regexp
	state         State
	stateFilePath string
	filePath      string
	fileInfoPrev  os.FileInfo
	fileInfoCurr  os.FileInfo
	posCurr       int64
	checkInterval time.Duration
	needToSave    bool
}
