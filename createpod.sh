#!/bin/bash

# Imposta il valore predefinito per la compilazione a 'false'
COMPILE_GOLANG=false

# Analizza le opzioni della riga di comando
# Se viene trovato il flag -c, imposta COMPILE_GOLANG a 'true'
while getopts ":c" opt; do
  case ${opt} in
    c)
      COMPILE_GOLANG=true
      ;;
    \?)
      echo "Opzione non valida: -$OPTARG" >&2
      exit 1
      ;;
  esac
done
shift $((OPTIND-1)) # Rimuove le opzioni analizzate, lasciando gli argomenti posizionali

# Assegna gli argomenti posizionali (POD_NAME, etc.)
POD_NAME=${1:-"default"}
HOST_PORT=${2:-"1883"}


# Esegue la compilazione solo se il flag -c è stato fornito
if [ "$COMPILE_GOLANG" = true ] ; then
    echo "--- Compilazione Go binary MQTT->TimescaleDB per Linux (flag -c rilevato) ---"
    cd go_to_timescale && CGO_ENABLED=0 GOOS=linux go build -v -a -o mqtt_to_timescale mqtt_to_timescale.go
    # Controlla se la compilazione è fallita
    if [ $? -ne 0 ]; then
        echo "ERRORE: La compilazione è fallita. Uscita."
        cd ..
        exit 1
    fi
    cd ..
else
    echo "--- Compilazione saltata (usa il flag -c per forzarla) ---"
fi


echo "Lanciando pod: digitaltwin-${POD_NAME} su porta MQTT:${HOST_PORT}"


# Sostituisci variabili e lancia
sed -e "s/\${POD_NAME}/${POD_NAME}/g" \
    -e "s/\${HOST_PORT}/${HOST_PORT}/g" \
    -e "s|\${PWD}|$(pwd)|g" \
    pod-template.yaml | podman play kube -

echo " Pod digitaltwin-${POD_NAME} avviato:"
echo "    MQTT Broker: 127.0.0.1:${HOST_PORT}"

echo "    Database: timescale/data/${POD_NAME}/"
