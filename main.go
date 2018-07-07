package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttBroker = flag.String("mqttBroker", "", "MQTT broker URI (mandatory). E.g.:192.168.1.1:1883")
	topic      = flag.String("topic", "", "Topic where hub-ctrl messages will be received")
)

var (
	verbose int
	listing int
	busNum  int
	devNum  int
	hub     = -1
)

type Message struct {
	hub   int
	port  int
	power bool
}

func connectMQTT() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + *mqttBroker)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("%s", token.Error())
	}

	return client, nil
}

func mqttCallback(client mqtt.Client, msg mqtt.Message) {

	var jsonMessage Message
	log.Printf("Message received: %s", msg.Payload())

	err := json.Unmarshal(msg.Payload(), &jsonMessage)

	if err != nil {
		log.Printf("Error parsing JSON: %s", err)
	}

	findHubs(1, 1, 0, 0, jsonMessage.hub)
	if jsonMessage.power == true {
		log.Printf("Powering up hub. Hub: %d, Port: %d", jsonMessage.hub, jsonMessage.port)
		sendCommandToHub(jsonMessage.hub, USB_REQ_SET_FEATURE, 8, jsonMessage.port)
	} else {
		log.Printf("Powering down hub. Hub: %d, Port: %d", jsonMessage.hub, jsonMessage.port)
		sendCommandToHub(jsonMessage.hub, USB_REQ_CLEAR_FEATURE, 8, jsonMessage.port)
	}

}

func init() {
	initUsb()
	flag.Parse()
}

func main() {

	//Check command line parameters
	if *mqttBroker == "" || *topic == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	//Channel used to block while receiving messages
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	clientMQTT, err := connectMQTT()
	if err != nil {
		log.Fatalf("Error connecting to MQTT broker: %s", err)
	}

	log.Printf("Connected to MQTT broker at %s", *mqttBroker)

	if token := clientMQTT.Subscribe(*topic, 0, mqttCallback); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to topic %s : %s", *topic, err)
	}

	log.Printf("Subscribed to topic %s", *topic)

	<-c

}
