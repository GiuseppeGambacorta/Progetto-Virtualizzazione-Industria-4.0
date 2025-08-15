package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := "tcp://localhost:1883"
	clientID := "go-mqtt-subscriber"

	// Parametro radice MQTT (default se non specificato)
	var rootTopic string
	if len(os.Args) < 2 {
		log.Fatal("inserire nome radice MQTT")
	} else {
		rootTopic = os.Args[1]
		fmt.Printf("Usando radice da parametro: %s\n", rootTopic)
	}

	// Auto-discovery con wildcard
	wildcardTopic := rootTopic + "/#" // # = tutti i sottotopic

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	// Handler dei messaggi ricevuti
	var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("[%s] %s\n", msg.Topic(), string(msg.Payload()))
	}

	opts.SetDefaultPublishHandler(messagePubHandler)

	client := MQTT.NewClient(opts)
	for {
		if token := client.Connect(); token.Wait() && token.Error() == nil {
			fmt.Println("Connected to Mosquitto!")
			break
		}
		fmt.Println("Waiting for Mosquitto...")
		time.Sleep(1 * time.Second)
	}

	if token := client.Subscribe(wildcardTopic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Errore sottoscrizione topic: %v", token.Error())
	}

	fmt.Printf("🔍 Auto-discovery attivo su: %s\n", wildcardTopic)
	fmt.Println("🎧 In attesa di messaggi... (Ctrl+C per uscire)")

	// Attendi segnale per chiudere (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n🔌 Disconnetto e chiudo client MQTT")
	client.Disconnect(250)
}
