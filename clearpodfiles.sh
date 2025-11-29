#!/bin/bash

POD_NAME=${1:-"default"}

if [ "$POD_NAME" = "default" ] && [ $# -eq 0 ]; then
    echo " Specifica il nome del pod da pulire!"
    echo "Utilizzo: $0 <POD_NAME>"
    echo "Esempio: $0 palletizer"
    echo ""
    echo "Pod con dati esistenti:"
    (ls -d timescale/data/*/ 2>/dev/null; ) | sed 's|.*/\(.*\)/|\1|' | sort -u
    exit 1
fi

echo "🧹 Pulizia dati per pod: ${POD_NAME}"

# Controlla se il pod è in esecuzione
if podman pod exists dt-${POD_NAME} 2>/dev/null; then
    echo "  Il pod dt-${POD_NAME} è ancora in esecuzione!"
    read -p "Vuoi fermarlo e rimuoverlo? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo " Fermando e rimuovendo pod..."
        podman pod stop dt-${POD_NAME}
        podman pod rm -f dt-${POD_NAME}
    else
        echo " Operazione annullata. Ferma il pod manualmente prima di pulire i dati."
        exit 1
    fi
fi

# Percorsi da pulire
TIMESCALE_PATH="timescale/data/${POD_NAME}"



echo " Directory da rimuovere:"
echo "    ${TIMESCALE_PATH}"
echo "    ${METABASE_PATH}"

# Verifica che le directory esistano
if [ ! -d "$TIMESCALE_PATH" ]; then
    echo "ℹ  Nessuna directory di dati trovata per il pod ${POD_NAME}"
    exit 0
fi

# Conferma prima di eliminare
read -p "  Sei sicuro di voler eliminare TUTTI i dati per questo pod? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo " Operazione annullata."
    exit 0
fi

# Rimuovi le directory
if [ -d "$TIMESCALE_PATH" ]; then
    echo "  Rimuovendo dati TimescaleDB..."
    rm -rf "$TIMESCALE_PATH"
    echo " Rimosso: $TIMESCALE_PATH"
fi




echo ""
echo " Pulizia completata per pod: ${POD_NAME}"
echo " Ora puoi ricreare il pod con: ./createpod.sh ${POD_NAME}"