package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MotorData struct {
	Position int `json:"Position"`
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	var rootTopic string
	mqttBroker := "tcp://localhost:"
	if len(os.Args) < 3 {
		log.Fatal("inserire nome pod e porta broker MQTT")
	} else {
		rootTopic = os.Args[1]
		mqttBroker = mqttBroker + os.Args[2]
		fmt.Printf("Usando radice da parametro: %s\n", rootTopic)
	}

	clientID := "go-mqtt-simulation"

	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqttBroker)
	opts.SetClientID(clientID)

	client := MQTT.NewClient(opts)
	for {
		if token := client.Connect(); token.Wait() && token.Error() == nil {
			fmt.Println("Connected to Mosquitto!")
			break
		}
		fmt.Println("Waiting for Mosquitto...")
		time.Sleep(1 * time.Second)
	}

	topic := rootTopic + "/ProductInfeedConveyors/Objects/FirstConveyor"
	position := 0

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Gestione chiusura pulita (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Printf(" Pubblicazione su %s ogni 100ms\n", topic)

loop:
	for {
		select {
		case <-ticker.C:
			position = (position + 5) % 360
			data := MotorData{Position: position}
			payload, _ := json.Marshal(data)
			token := client.Publish(topic, 0, false, payload)
			token.Wait()
			fmt.Printf("Pubblicato: %s\n", string(payload))
		case <-sigChan:
			fmt.Println("\n🔌 Disconnetto e chiudo client MQTT")
			break loop
		}
	}

	client.Disconnect(250)
}
