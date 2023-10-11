package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestBufferOverflow(t *testing.T) {
	logFilePath := fmt.Sprintf("/tmp/logalert_test_%d", time.Now().UnixNano()%1000)
	data := make1KbExtraString()
	createFileWithData(logFilePath, data)
	defer func() {
		err := os.Remove(logFilePath)
		if err != nil {
			t.Error(err)
		}
	}()

	logCfgTest := FileConfig{
		Path:           logFilePath,
		ReadBufferSize: "1kb",
	}

	logWatcher, err := NewWatcher(logCfgTest, nil)
	if err != nil {
		t.Error(err)
	}

	lines, err := logWatcher.getNewLines()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, []byte(lines[0])) {
		t.Errorf("Expected string: '%s'\n received: '%s'", data, lines[0])
	}
}

func make1KbExtraString() []byte {
	data := make([]byte, 1024)
	for i, _ := range data {
		data[i] = '0'
	}
	return append(data, []byte("Extra test string")...)
}

func createFileWithData(path string, data []byte) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
}
