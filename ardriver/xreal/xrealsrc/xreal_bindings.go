//go:build xreal
// +build xreal

package xreal

// #cgo CFLAGS: -w
// #cgo pkg-config: json-c libusb-1.0 hidapi-libusb
// #include "go_ffi.h"
// #include "device.h"
// #include "device_imu.h"
// #include "device_mcu.h"
import "C"
import (
	"fmt"
	"sync"
	"time"

	"git.terah.dev/UnrealXR/unrealxr/ardriver/commons"
)

var (
	deviceEventHandlerMutex = sync.Mutex{}
	deviceEventListener     *commons.AREventListener
)

//export goIMUEventHandler
func goIMUEventHandler(_ C.uint64_t, event_type C.device_imu_event_type, ahrs *C.struct_device_imu_ahrs_t) {
	if deviceEventListener == nil {
		return
	}

	if event_type != C.DEVICE_IMU_EVENT_UPDATE {
		return
	}

	orientation := C.device_imu_get_orientation(ahrs)
	euler := C.device_imu_get_euler(orientation)

	deviceEventListener.PitchCallback(float32(euler.pitch))
	deviceEventListener.RollCallback(float32(euler.roll))
	deviceEventListener.YawCallback(float32(euler.yaw))
}

// Implements commons.ARDevice
type XrealDevice struct {
	eventListener *commons.AREventListener
	imuDevice     *C.struct_device_imu_t
	deviceIsOpen  bool
}

func (device *XrealDevice) Initialize() error {
	if device.deviceIsOpen {
		return fmt.Errorf("device is already open")
	}

	device.imuDevice = &C.struct_device_imu_t{}

	// (*[0]byte) is a FUBAR way to cast a pointer to a function, but unsafe.Pointer doesn't work:
	// cannot use unsafe.Pointer(_Cgo_ptr(_Cfpvar_fp_imuEventHandler)) (value of type unsafe.Pointer) as *[0]byte value in variable declaration
	if C.DEVICE_IMU_ERROR_NO_ERROR != C.device_imu_open(device.imuDevice, (*[0]byte)(C.imuEventHandler)) {
		return fmt.Errorf("failed to open IMU device")
	}

	C.device_imu_clear(device.imuDevice)
	C.device_imu_calibrate(device.imuDevice, 1000, true, true, false)

	device.deviceIsOpen = true

	// let's hope this doesn't cause race conditions
	go func() {
		for device.eventListener == nil {
			time.Sleep(time.Millisecond * 10)
		}

		for {
			if !device.deviceIsOpen {
				break
			}

			// I'm sorry.
			deviceEventHandlerMutex.Lock()
			deviceEventListener = device.eventListener
			status := C.device_imu_read(device.imuDevice, -1)
			deviceEventHandlerMutex.Unlock()

			if status != C.DEVICE_IMU_ERROR_NO_ERROR {
				break
			}
		}

		device.deviceIsOpen = false
		C.device_imu_close(device.imuDevice)
	}()

	return nil
}

func (device *XrealDevice) End() error {
	if !device.deviceIsOpen {
		return fmt.Errorf("device is not open")
	}

	C.device_imu_close(device.imuDevice)
	device.deviceIsOpen = false
	return nil
}

func (device *XrealDevice) IsPollingLibrary() bool {
	return false
}

func (device *XrealDevice) IsEventBasedLibrary() bool {
	return true
}

func (device *XrealDevice) Poll() error {
	return fmt.Errorf("not supported")
}

func (device *XrealDevice) RegisterEventListeners(listener *commons.AREventListener) {
	device.eventListener = listener
}
