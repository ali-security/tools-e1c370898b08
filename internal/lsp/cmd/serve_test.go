// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"net"
	"testing"
)

func TestListenParsing(t *testing.T) {
	tests := []struct {
		input, wantNetwork, wantAddr string
	}{
		{"127.0.0.1:0", "tcp", "127.0.0.1:0"},
		{"unix;/tmp/sock", "unix", "/tmp/sock"},
		{"auto", "auto", ""},
		{"auto;foo", "auto", "foo"},
	}

	for _, test := range tests {
		gotNetwork, gotAddr := parseAddr(test.input)
		if gotNetwork != test.wantNetwork {
			t.Errorf("network = %q, want %q", gotNetwork, test.wantNetwork)
		}
		if gotAddr != test.wantAddr {
			t.Errorf("addr = %q, want %q", gotAddr, test.wantAddr)
		}
	}
}

// TestListenAddrImplicitHostIsLoopback checks that an address with no host, as
// produced by "gopls serve -port=N" or passed as "-listen=:N", is bound to
// the loopback interface instead of all network interfaces.
func TestListenAddrImplicitHostIsLoopback(t *testing.T) {
	ln, err := net.Listen("tcp", listenAddr("tcp", ":0"))
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	ip := ln.Addr().(*net.TCPAddr).IP
	if !ip.IsLoopback() {
		t.Errorf("listenAddr(\"tcp\", \":0\") bound %v, want a loopback address", ln.Addr())
	}
}

// TestListenAddrExplicitHost checks that an explicit host is honored, including
// 0.0.0.0 for callers that deliberately want all network interfaces, and that
// non-TCP addresses are left unchanged.
func TestListenAddrExplicitHost(t *testing.T) {
	tests := []struct {
		network, addr string
	}{
		{"tcp", "127.0.0.1:0"},
		{"tcp", "0.0.0.0:0"},
		{"tcp4", "0.0.0.0:0"},
		{"unix", ":sock"},
	}
	for _, test := range tests {
		if got := listenAddr(test.network, test.addr); got != test.addr {
			t.Errorf("listenAddr(%q, %q) = %q, want %q", test.network, test.addr, got, test.addr)
		}
	}
}
