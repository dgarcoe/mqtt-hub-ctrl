package main

// #cgo LDFLAGS: -lusb
// #include "hub-ctrl.h"
import (
	"C"
)
import (
	"fmt"
	"log"
)

const USB_REQ_CLEAR_FEATURE = 1
const USB_REQ_SET_FEATURE = 3

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

func getHub(busnum, devnum int) int {
	return int(C.get_hub(C.int(busnum), C.int(devnum)))
}

func sendCommandToHub(hub, request, feature, index int) error {

	log.Printf("Sending command to hub")
	if C.send_command(C.int(hub), C.int(request), C.int(feature), C.int(index)) < 0 {
		return fmt.Errorf("Error sending commang to hub")
	}

	return nil
}
