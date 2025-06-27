//go:build xreal && !xreal_debug_logging
// +build xreal,!xreal_debug_logging

package xreal

// #cgo CFLAGS: -DNDEBUG
import "C"
