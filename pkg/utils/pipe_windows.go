//go:build windows
// +build windows
package utils

import (
	"os"

	"golang.org/x/sys/windows"
)

// MkNamedPipe creates a named pipe on Windows systems.
// It receives the path of the pipe and the mode permissions (ignored on Windows).
func MkNamedPipe(path string, mode uint32) error {
	// Windows uses a different naming convention for pipes
	// Convert the path to a Windows pipe format if needed
	pipeName := path
	if len(pipeName) > 0 && pipeName[0] != '\\' {
		// Format for Windows named pipes is: \\.\pipe\PipeName
		pipeName = `\\.\pipe\` + pipeName
	}

	// Default security attributes
	sa := &windows.SecurityAttributes{
		Length:             0,
		InheritHandle:      0,
		SecurityDescriptor: nil,
	}

	// Create the named pipe
	handle, err := windows.CreateNamedPipe(
		windows.StringToUTF16Ptr(pipeName),
		windows.PIPE_ACCESS_DUPLEX,
		windows.PIPE_TYPE_BYTE|windows.PIPE_READMODE_BYTE|windows.PIPE_WAIT,
		1,    // Max instances
		4096, // Out buffer size
		4096, // In buffer size
		0,    // Default timeout
		sa,
	)

	if err != nil {
		return os.NewSyscallError("CreateNamedPipe", err)
	}

	// Close the handle - since we're just creating the pipe
	// The actual file operations will use os.OpenFile
	windows.CloseHandle(handle)

	return nil
}
