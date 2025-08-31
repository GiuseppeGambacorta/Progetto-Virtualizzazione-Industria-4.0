#!/bin/bash

POD_NAME=${1:-"default"}
HOST_PORT=${2:-"1883"}
GRAFANA_PORT=${3:-"5000"}

#echo "Compilando Go binary MQTT->TimescaleDB per Linux..."
#cd go_to_timescale && GOOS=linux GOARCH=arm64 go build -o mqtt_to_timescale mqtt_to_timescale.go
#cd ..


#echo "Compilando Go binary Simulazione per Linux..."
#cd go_simulation && GOOS=linux GOARCH=arm64 go build -o go_simulation go_simulation.go
#cd ..


echo " Generando dashboard per: ${POD_NAME}"
# Genera dashboard con POD_NAME sostituito
sed "s/\${POD_NAME}/${POD_NAME}/g" \
    grafana/provisioning/dashboards/mqtt-dashboard-template.json > \
    grafana/provisioning/dashboards/mqtt-dashboard.json


echo "Lanciando pod: digitaltwin-${POD_NAME} su porta MQTT:${HOST_PORT}, Grafana:${GRAFANA_PORT}"


# Sostituisci variabili e lancia
sed -e "s/\${POD_NAME}/${POD_NAME}/g" \
    -e "s/\${HOST_PORT}/${HOST_PORT}/g" \
    -e "s/\${GRAFANA_PORT}/${GRAFANA_PORT}/g" \
    -e "s|\${PWD}|$(pwd)|g" \
    pod-template.yaml | podman play kube -

echo " Pod digitaltwin-${POD_NAME} avviato:"
echo "    MQTT Broker: 127.0.0.1:${HOST_PORT}"
echo "    Grafana: http://127.0.0.1:${GRAFANA_PORT} (admin/admin)"
echo "    Database: timescale/data/${POD_NAME}/"
echo "    Grafana data: grafana/data/${POD_NAME}/"