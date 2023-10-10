package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

const rotationDepth = 3

type LogWatcher struct {
	hostname        string
	buf             []byte
	dateReg         *regexp.Regexp
	state           State
	stateFilePath   string
	logFilePath     string
	logFileInfoPrev os.FileInfo
	logFileInfoCurr os.FileInfo
	posCurr         int64
	checkInterval   time.Duration
	filters         []*Filter
	sender          *Sender
	needToSave      bool
}

func NewLogWatcher(hostname string, cfg LogFileConfig, sender *Sender) *LogWatcher {
	var statePath string

	switch runtime.GOOS {
	case "linux":
		statePath = "/var/lib/logalert"
	case "darwin":
		statePath = "/usr/local/var/lib/logalert"
	default:
		log.Fatal(runtime.GOOS, " OS is not supported")
	}

	_, err := os.Stat(statePath)
	if err != nil {
		log.Fatal(err)
	}

	w := LogWatcher{
		hostname:      hostname,
		buf:           newBuffer(cfg.ReadBufferSize),
		logFilePath:   cfg.Path,
		checkInterval: time.Second * time.Duration(cfg.IntervalSec),
		sender:        sender,
	}

	logFilePathHash := md5.Sum([]byte(cfg.Path))
	w.stateFilePath = fmt.Sprintf("%s/%x", statePath, logFilePathHash)

	if len(w.buf) == 0 {
		log.Fatalf("Can't create read buffer with size %s for logfile %s",
			cfg.ReadBufferSize,
			cfg.Path,
		)
	}

	file, err := os.Open(cfg.Path)
	if err != nil {
		log.Fatal(err)
	}

	w.logFileInfoCurr, err = file.Stat()
	if err != nil {
		file.Close()
		log.Fatal(err)
	}

	file.Close()

	err = w.saveState(w.logFileInfoCurr, 0)
	if err != nil {
		log.Fatal(err)
	}

	w.dateReg, err = regexp.Compile(cfg.DateFormat)
	if err != nil {
		log.Fatalf("LogFile %s date pattern compile error: %v", cfg.Path, err)
	}

	for _, filterCfg := range cfg.Filters {
		filter, err := NewFilter(filterCfg, w.hostname)
		if err != nil {
			log.Fatalf("NewFilter error: %v", err)
		}

		w.filters = append(w.filters, filter)
	}

	return &w
}

func (w *LogWatcher) saveState(fileInfo os.FileInfo, pos int64) error {
	if !w.needToSave {
		return nil
	}

	file, err := os.Create(w.stateFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sys := fileInfo.Sys().(*syscall.Stat_t)

	state := State{
		Pos: pos,
		Dev: sys.Dev,
		Ino: sys.Ino,
	}

	b, err := json.Marshal(state)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(file, b)
	if err != nil {
		return err
	}

	w.state = state
	w.needToSave = false
	w.logFileInfoPrev = fileInfo

	return nil
}

func (w *LogWatcher) watch(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	err := w.logParsingAndSendMessages(ctx)
	if err != nil {
		log.Fatal("[ERROR] logParsingAndSendNotifications: ", err)
	}

	err = w.saveState(w.logFileInfoCurr, w.posCurr)
	if err != nil {
		log.Fatalf("[ERROR] state update error: %v log file: %s", err, w.logFilePath)
	}

	ticker := time.NewTicker(w.checkInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err = w.logParsingAndSendMessages(ctx)
			if err != nil {
				log.Println("[ERROR] logParsingAndSendNotifications: ", err)
				continue
			}

			err = w.saveState(w.logFileInfoCurr, w.posCurr)
			if err != nil {
				log.Fatalf("[ERROR] state update error: %v log file: %s", err, w.logFilePath)
			}
		}
	}
}

func (w *LogWatcher) logParsingAndSendMessages(ctx context.Context) error {
	lines, err := w.getNewLines()
	if err != nil {
		return fmt.Errorf("getNewLines error: %v logFile: %s", err, w.logFilePath)
	}

	messages := processLines(lines, w.filters, w.dateReg)

	for _, msg := range messages {
		if err = w.sender.SendMessage(ctx, msg); err != nil {
			return fmt.Errorf("SendMessage error: %v", err)
		}
	}

	if len(messages) > 0 {
		w.needToSave = true
	}

	return nil
}

func (w *LogWatcher) getNewLines() ([]string, error) {
	fileIndex, fileInfo := searchFile(w.logFilePath, w.logFileInfoPrev)
	if fileIndex == -1 {
		return nil, fmt.Errorf("Can't find log file '%s' state. The file was unexpectedly changed.", w.logFilePath)
	}

	var (
		pos         = w.state.Pos
		linesResult []string
	)

	for i := fileIndex; i >= 0; i-- {
		filePath := w.logFilePath
		if i > 0 {
			filePath = fmt.Sprintf("%s.%d", w.logFilePath, i)
		}

		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		var (
			readErr  error
			partLine string
			n        int
		)

		for readErr != io.EOF {
			n, readErr = file.ReadAt(w.buf, pos)
			if readErr != nil && readErr != io.EOF {
				file.Close()
				return nil, fmt.Errorf("file read error: %v file: %s", err, filePath)
			}

			pos += int64(n)

			lines := strings.Split(string(w.buf[:n]), "\n")

			if partLine != "" {
				lines[0] = partLine + lines[0]
				partLine = ""
			}

			if readErr != io.EOF && w.buf[len(w.buf)-1] != '\n' {
				partLine = lines[len(lines)-1]
				lines = lines[:len(lines)-1]
			}

			linesResult = append(linesResult, lines...)
		}

		fileInfo, err = file.Stat()
		if err != nil {
			return nil, err
		}

		file.Close()

		if i > 0 {
			pos = 0
		}
	}

	w.logFileInfoCurr = fileInfo
	w.posCurr = pos

	return linesResult, nil
}

func processLines(lines []string, filters []*Filter, dateReg *regexp.Regexp) []Message {
	matchMaps := make([]map[string]int, len(filters))

	for fIndex, _ := range filters {
		matchMaps[fIndex] = make(map[string]int)
	}

	for _, line := range lines {
		for fIndex, filter := range filters {
			if filter.Match(line) {
				line, _ = lineRemoveDate(line, dateReg)

				if _, ok := matchMaps[fIndex][line]; !ok {
					matchMaps[fIndex][line] = 0
				}

				matchMaps[fIndex][line]++
			}
		}
	}

	var messages []Message

	for fIndex, filter := range filters {
		for line, count := range matchMaps[fIndex] {
			messages = append(messages, Message{
				Text:   line,
				Count:  count,
				Filter: filter,
			})
		}
	}

	return messages
}

// Returns file index and os.FileInfo
// Example: filePath: /foo/bar/file
// index = 0 on path /foo/bar/file founded
// index = 1 on path /foo/bar/file.1 founded
// index = -1 if not founded
func searchFile(filePath string, fileInfo os.FileInfo) (int, os.FileInfo) {
	for i := 0; i <= rotationDepth; i++ {
		path := filePath
		if i > 0 {
			path = fmt.Sprintf("%s.%d", filePath, i)
		}

		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		fInfo, err := f.Stat()
		if err != nil {
			f.Close()
			log.Fatal(err)
		}

		f.Close()

		if os.SameFile(fileInfo, fInfo) {
			return i, fInfo
		}
	}

	return -1, fileInfo
}

func lineRemoveDate(str string, reDate *regexp.Regexp) (string, bool) {
	if str == "" || reDate == nil {
		return str, false
	}

	matches := reDate.FindStringSubmatch(str)
	if len(matches) == 1 {
		return strings.Replace(str, matches[0], "", 1), true
	}

	return str, false
}

func removeDuplicates(strSlice []string) []string {
	keys := make(map[string]struct{}, len(strSlice))
	dedup := make([]string, 0, len(strSlice))
	for _, str := range strSlice {
		if _, ok := keys[str]; !ok {
			keys[str] = struct{}{}
			dedup = append(dedup, str)
		}
	}
	return dedup
}
