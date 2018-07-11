# mqtt-hub-ctrl
Hub-ctrl with MQTT interface for remote control of USB power. Programmed in Go and using the C library with CGo. The reason for doing this is simply to learn and test how C integration works with Go.

## Installation

A Dockerfile for ARM is provided although it has not been thoroughly tested for the moment and it will probably fail.

Otherwise, you can install the software with the following commands:

* $ sudo apt-get install libusb-dev
* $ go get -u github.com/eclipse/paho.mqtt.golang
* $ go install

## Running the application

To run the application you can use the following command:

```
sudo ./mqtt-hub-ctrl -mqttBroker localhost:1883 -topic hubPower
```

Using an MQTT client you can send commands to the topic in JSON format following the conventions of the C hub-ctrl. For example, to power off all the USB ports in a Raspberry you should send the following JSON message through the specified topic:

```
{"hub":0,"port":2,"power":false}
```
To power them on again you should send:

```
{"hub":0,"port":2,"power":true}
```

## Original copyright
Copyright (C) 2006 Free Software Initiative of Japan

Author: NIIBE Yutaka  <gniibe at fsij.org>

This file can be distributed under the terms and conditions of the GNU General Public License version 2 (or later).

Updated from the C version published on Github by Joel Dare (codazoda)
