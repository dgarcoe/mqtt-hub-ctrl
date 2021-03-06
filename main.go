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
	topic      = flag.String("topic", "", "Topic where hub-ctrl messages will be received (mandatory)")
	user       = flag.String("user", "", "MQTT username")
	pwd        = flag.String("password", "", "MQTT password")
	verbose    = flag.Int("verbose", 0, "Set it to 1 to print more information from the hubs")
)

var (
	listing int
	busNum  int
	devNum  int
	hub     = -1
)

//Message Used to hold MQTT JSON messages
type Message struct {
	Hub   int
	Port  int
	Power bool
}

//Connect to the MQTT broker
func connectMQTT() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + *mqttBroker)

	if *user != "" && *pwd != "" {
		opts.SetUsername(*user).SetPassword(*pwd)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("%s", token.Error())
	}

	return client, nil
}

//Callback for MQTT messages received through the subscribed topic
func mqttCallback(client mqtt.Client, msg mqtt.Message) {

	var jsonMessage Message
	log.Printf("Message received: %s", msg.Payload())

	err := json.Unmarshal(msg.Payload(), &jsonMessage)
	if err != nil {
		log.Printf("Error parsing JSON: %s", err)
	}

	hub := jsonMessage.Hub
	port := jsonMessage.Port
	power := jsonMessage.Power

	findHubs(listing, *verbose, busNum, devNum, hub)
	if power == true {
		log.Printf("Powering up hub. Hub: %d, Port: %d", hub, port)
		sendPowerCommandToHub(hub, UsbReqSetFeature, port)
	} else {
		log.Printf("Powering down hub. Hub: %d, Port: %d", hub, port)
		sendPowerCommandToHub(hub, UsbReqClearFeature, port)
	}

}

func init() {
	initUsb()
	flag.Parse()
	if *verbose == 1 {
		listing = 1
	}
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
