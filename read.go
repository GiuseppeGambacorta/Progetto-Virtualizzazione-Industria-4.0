package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := "tcp://192.168.1.103:1883" // cambia con IP broker MQTT
	topic := "Palletizzatore/MainMachine/Conteggio"
	clientID := "go-mqtt-subscriber"

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	// Handler dei messaggi ricevuti
	var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("Messaggio ricevuto sul topic %s: %s\n", msg.Topic(), string(msg.Payload()))
	}

	opts.SetDefaultPublishHandler(messagePubHandler)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Errore connessione MQTT: %v", token.Error())
	}

	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Errore sottoscrizione topic: %v", token.Error())
	}

	fmt.Printf("Sottoscritto a topic %s, in attesa di messaggi...\n", topic)

	// Attendi segnale per chiudere (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Disconnetto e chiudo client MQTT")
	client.Disconnect(250)
}
