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

//UsbReqClearFeature Clear a USB feature
const UsbReqClearFeature = 1

//UsbReqSetFeature Set a USB feature
const UsbReqSetFeature = 3

//Initialize USB
func initUsb() {
	C.usb_init()
	C.usb_find_busses()
	C.usb_find_devices()
}

//Find available hubs
func findHubs(listing, verbose, busnum, devnum, hub int) error {
	if C.usb_find_hubs(C.int(listing), C.int(verbose), C.int(busnum), C.int(devnum), C.int(hub)) <= 0 {
		return fmt.Errorf("No hubs found")
	}
	return nil
}

//Get one specific hub
func getHub(busnum, devnum int) int {
	return int(C.get_hub(C.int(busnum), C.int(devnum)))
}

//Send a command to modify the power of a port in a USB hub
func sendPowerCommandToHub(hub, request, index int) error {

	log.Printf("Sending command to hub")
	if C.send_command(C.int(hub), C.int(request), C.USB_PORT_FEAT_POWER, C.int(index)) < 0 {
		return fmt.Errorf("Error sending commang to hub")
	}

	return nil
}
