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

func initUsb() {
	C.usb_init()
	C.usb_find_busses()
	C.usb_find_devices()
}

func findHubs(listing, verbose, busnum, devnum, hub int) error {
	//	if C.usb_find_hubs(C.int(listing), C.int(verbose), C.int(busnum), C.int(devnum), C.int(hub)) <= 0 {
	//		return fmt.Errorf("No hubs found")
	//	}
	return nil
}

/*func getHub(busnum, devnum int) int {
	return int(C.get_hub(C.int(busnum), C.int(devnum)))
}*/

func sendCommandToHub() error {

	log.Printf("Sending command to hub")
	if C.send_command(C.int(0), C.USB_REQ_SET_FEATURE, C.int(8), C.int(2)) < 0 {
		return fmt.Errorf("Error sending commang to hub")
	}

	return nil
}
