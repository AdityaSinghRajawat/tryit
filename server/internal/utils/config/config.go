// Package config loads runtime configuration from environment variables per
// IMPL §10.1. Defaults are encoded here; nothing reaches the rest of the app
// except through a Config value.
package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port                int
	SecretsBackend      string // "env" (Phase 1), "keychain", "file" (Phase 2)
	ExecTimeout         time.Duration
	ExecMaxResponseSize int64
	ExecInsecureTLS     bool
	LogLevel            string
	PairFile            string // ~/.tryit/pair.json (Phase 1)
	HomeDir             string // resolved $HOME (or empty)
}

func Load() (Config, error) {
	c := Config{
		Port:                8765,
		SecretsBackend:      "env",
		ExecTimeout:         30 * time.Second,
		ExecMaxResponseSize: 10 * 1024 * 1024,
		ExecInsecureTLS:     false,
		LogLevel:            "info",
	}
	if v := os.Getenv("TRYIT_PORT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 65535 {
			return c, fmt.Errorf("TRYIT_PORT: invalid %q", v)
		}
		c.Port = n
	}
	if v := os.Getenv("TRYIT_SECRETS_BACKEND"); v != "" {
		c.SecretsBackend = v
	}
	if v := os.Getenv("TRYIT_EXEC_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return c, fmt.Errorf("TRYIT_EXEC_TIMEOUT: %w", err)
		}
		c.ExecTimeout = d
	}
	if v := os.Getenv("TRYIT_EXEC_MAX_RESPONSE_BYTES"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			return c, fmt.Errorf("TRYIT_EXEC_MAX_RESPONSE_BYTES: invalid %q", v)
		}
		c.ExecMaxResponseSize = n
	}
	if v := os.Getenv("TRYIT_EXEC_INSECURE_TLS"); v != "" {
		c.ExecInsecureTLS = strings.EqualFold(v, "true") || v == "1"
	}
	if v := os.Getenv("TRYIT_LOG_LEVEL"); v != "" {
		c.LogLevel = strings.ToLower(v)
	}
	c.HomeDir, _ = os.UserHomeDir()
	c.PairFile = os.Getenv("TRYIT_PAIR_FILE")
	if c.PairFile == "" && c.HomeDir != "" {
		c.PairFile = c.HomeDir + "/.tryit/pair.json"
	}
	return c, nil
}

func (c Config) ListenAddr() string {
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(c.Port))
}

func (c Config) HostHeader() string {
	return "127.0.0.1:" + strconv.Itoa(c.Port)
}
