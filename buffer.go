package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

const bufferSizeLimit = 10_000_000

// newBuffer allocates and returns new bytes buffer
// available values: 1 Kb - bufferSizeLimit
// bufSize example: "10Kb", "1mb", "50KB", etc
func newBuffer(bufSize string) []byte {
	r := regexp.MustCompile(`^(\d+)\s?([bBkKmM]{1,2})$`)

	matches := r.FindStringSubmatch(bufSize)
	if len(matches) == 0 {
		return []byte{}
	}

	size, err := strconv.Atoi(matches[1])
	if err != nil || size <= 0 {
		log.Println("bufSize incorrect number", matches[1])
		return []byte{}
	}

	switch strings.ToLower(matches[2]) {
	case "k", "kb":
		size *= 1024
	case "m", "mb":
		size *= 1024 * 1024
	default:
		log.Println("[ERROR] bufSize incorrect representation", matches[2])
		return []byte{}
	}

	if size > bufferSizeLimit {
		log.Printf("[ERROR] buffer size %d exceeds the limit %d", size, bufferSizeLimit)
		return []byte{}
	}

	return make([]byte, size)
}
