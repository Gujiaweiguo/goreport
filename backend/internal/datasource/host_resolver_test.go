package datasource

import "testing"

func TestResolveHost_NoAlias(t *testing.T) {
	t.Setenv(datasourceLocalhostAliasEnv, "")

	tests := []struct {
		name     string
		host     string
		expected string
	}{
		{name: "localhost unchanged", host: "localhost", expected: "localhost"},
		{name: "ipv4 loopback unchanged", host: "127.0.0.1", expected: "127.0.0.1"},
		{name: "ipv6 loopback unchanged", host: "::1", expected: "::1"},
		{name: "remote host unchanged", host: "db.internal", expected: "db.internal"},
		{name: "trim spaces", host: "  db.internal  ", expected: "db.internal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveHost(tt.host)
			if got != tt.expected {
				t.Fatalf("ResolveHost(%q) = %q, want %q", tt.host, got, tt.expected)
			}
		})
	}
}

func TestResolveHost_WithAlias(t *testing.T) {
	t.Setenv(datasourceLocalhostAliasEnv, "host.docker.internal")

	tests := []struct {
		name     string
		host     string
		expected string
	}{
		{name: "localhost rewritten", host: "localhost", expected: "host.docker.internal"},
		{name: "localhost rewritten case insensitive", host: "LOCALHOST", expected: "host.docker.internal"},
		{name: "ipv4 loopback rewritten", host: "127.0.0.1", expected: "host.docker.internal"},
		{name: "ipv6 loopback rewritten", host: "::1", expected: "host.docker.internal"},
		{name: "remote host unchanged", host: "goreport-mysql", expected: "goreport-mysql"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveHost(tt.host)
			if got != tt.expected {
				t.Fatalf("ResolveHost(%q) = %q, want %q", tt.host, got, tt.expected)
			}
		})
	}
}
