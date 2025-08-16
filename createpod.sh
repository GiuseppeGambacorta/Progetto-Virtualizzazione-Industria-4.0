#!/bin/bash

POD_NAME=${1:-"default"}
HOST_PORT=${2:-"1883"}

echo "🔨 Compilando Go binary per Linux..."
cd go && GOOS=linux GOARCH=arm64 go build -o mqtt_to_timescale mqtt_to_timescale.go
cd ..

echo "📁 Creando directory database per pod: ${POD_NAME}"
mkdir -p timescale/data/${POD_NAME}

echo "Lanciando pod: digitaltwin-${POD_NAME} su porta ${HOST_PORT}"

# Sostituisci variabili e lancia]
sed -e "s/\${POD_NAME}/${POD_NAME}/g" \
    -e "s/\${HOST_PORT}/${HOST_PORT}/g" \
    pod-template.yaml | podman play kube -
    
echo "Pod digitaltwin-${POD_NAME} avviato su 127.0.0.1:${HOST_PORT}"
echo "🗄️ Database in: timescale/data/${POD_NAME}/"