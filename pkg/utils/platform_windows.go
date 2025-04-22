//go:build windows
// +build windows
package utils

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
)

func defaultProtocol() string {
	return "npipe"
}

func getPlatformDialer(protocol string) func(ctx context.Context, addr string) (net.Conn, error) {
	if protocol == "tcp" {
		return func(ctx context.Context, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
		}
	}
	return func(ctx context.Context, addr string) (net.Conn, error) {
		// Ensure the address is in the correct format for named pipes
		if len(addr) > 0 && addr[0] != '\\' {
			addr = `\\.\pipe\` + addr
		}
		return winio.DialPipeContext(ctx, addr)
	}
}
