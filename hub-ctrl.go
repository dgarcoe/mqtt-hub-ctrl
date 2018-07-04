package main

// #cgo LDFLAGS: -lusb
// #include "hub-ctrl.h"
import (
	"C"
)
import (
	"fmt"
)

func initUsb() {
	C.usb_init()
	C.usb_find_busses()
	C.usb_find_devices()
}

func findHubs(listing, verbose, busnum, devnum, hub int) error {
	if C.usb_find_hubs(C.int(listing), C.int(verbose), C.int(busnum), C.int(devnum), C.int(hub)) <= 0 {
		return fmt.Errorf("No hubs found")
	}
	return nil
}
