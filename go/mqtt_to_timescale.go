package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq" // È come dire: "Carica il plugin PostgreSQL" , il compilatore con _ non si lamenta
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func waitForService(host, port, serviceName string, maxRetries int) {
	for i := 0; i < maxRetries; i++ {
		log.Printf("🔄 Tentativo connessione %s %d/%d...", serviceName, i+1, maxRetries)

		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 2*time.Second)
		if err == nil {
			conn.Close()
			log.Printf("%s è pronto!", serviceName)
			return
		}

		log.Printf("⏳ %s non pronto, aspetto 3 secondi...", serviceName)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("Impossibile connettersi a %s dopo %d tentativi", serviceName, maxRetries)
}

func connectDB(connStr string, maxRetries int) *sql.DB {
	for i := 0; i < maxRetries; i++ {
		log.Printf("🔄 Tentativo connessione database %d/%d...", i+1, maxRetries)

		// sql.Open usa automaticamente il driver "postgres" registrato da lib/pq
		db, err := sql.Open("postgres", connStr)
		if err == nil {
			if err := db.Ping(); err == nil {
				log.Printf("✅ Connesso a TimescaleDB!")
				return db
			}
		}

		log.Printf("⏳ Database non pronto, aspetto 3 secondi...")
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("❌ Impossibile connettersi al database dopo %d tentativi", maxRetries)
	return nil
}

func insertMQTTData(db *sql.DB, topic, payload, podName string) error {
	query := `INSERT INTO mqtt_data (time, topic, value, pod_name) VALUES (NOW(), $1, $2, $3)`

	// Verifica se è JSON valido
	var jsonValue json.RawMessage
	if json.Valid([]byte(payload)) {
		jsonValue = json.RawMessage(payload)
	} else {
		// Se non è JSON, wrappa come stringa
		jsonStr, _ := json.Marshal(payload)
		jsonValue = json.RawMessage(jsonStr)
	}

	_, err := db.Exec(query, topic, jsonValue, podName)
	return err
}

func main() {

	var rootTopic string
	if len(os.Args) < 2 {
		log.Fatal("inserire nome radice MQTT")
	} else {
		rootTopic = os.Args[1]
		fmt.Printf("Usando radice da parametro: %s\n", rootTopic)
	}

	// Auto-discovery con wildcard
	wildcardTopic := rootTopic + "/#" // # = tutti i sottotopic

	clientID := "go-mqtt-subscriber"
	mqttBroker := getEnv("MQTT_BROKER", "tcp://localhost:1883")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "mqtt_data")
	dbUser := getEnv("DB_USER", "mqtt_user")
	dbPassword := getEnv("DB_PASSWORD", "mqtt_password")

	waitForService(dbHost, dbPort, "TimeScaleDB", 30)

	// Connessione database
	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db := connectDB(dbConnStr, 30)
	defer db.Close()

	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqttBroker)
	opts.SetClientID(clientID)

	// Handler dei messaggi ricevuti
	var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		payload := string(msg.Payload())

		// Stampa per debug
		fmt.Printf("[%s] %s = %s\n", time.Now().Format("15:04:05"), msg.Topic(), payload)

		// Salva in TimescaleDB
		if err := insertMQTTData(db, msg.Topic(), payload, rootTopic); err != nil {
			log.Printf("Errore inserimento DB per topic %s: %v", msg.Topic(), err)
		}
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
