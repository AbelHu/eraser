//go:build !windows
// +build !windows

package utils

import (
	"golang.org/x/sys/unix"
)

// MkNamedPipe creates a named pipe (FIFO) on Unix-like systems.
// It receives the path of the pipe and the mode permissions.
func MkNamedPipe(path string, mode uint32) error {
	return unix.Mkfifo(path, mode)
}
