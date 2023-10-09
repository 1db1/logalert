//go:build linux

package main

type State struct {
	Pos int64
	Dev uint64
	Ino uint64
}
