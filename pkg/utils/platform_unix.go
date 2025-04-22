//go:build !windows
// +build !windows
package utils

import (
	"context"
	"net"
)

const unixProtocol = "unix"

func defaultProtocol() string {
	return unixProtocol
}

func getPlatformDialer(protocol string) func(ctx context.Context, addr string) (net.Conn, error) {
	if protocol == "tcp" {
		return func(ctx context.Context, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
		}
	}
	return func(ctx context.Context, addr string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, unixProtocol, addr)
	}
}
