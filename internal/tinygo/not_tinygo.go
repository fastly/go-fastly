//go:build !tinygo

// Package tinygo reports whether the code is being built with TinyGo, so
// that tests relying on unsupported features can skip themselves.
package tinygo

// Enabled is true when built with TinyGo.
const Enabled = false
