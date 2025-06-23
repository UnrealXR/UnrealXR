//go:build xreal
// +build xreal

package xreal

// #include "evdi_lib.h"
// #include "go_ffi.h"
// #cgo CFLAGS: -w
// #cgo pkg-config: json-c libusb-1.0 hidapi-libusb
import "C"

var IsXrealEnabled = true
