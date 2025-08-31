# Digital Twin - Prova Lettura MQTT

## Requisiti

- [Podman](https://podman.io/) installato e funzionante
- [Go](https://golang.org/) installato (per compilare i binari)
- Sistema operativo host: **macOS**, **Linux** o **Windows**

## 1. Compilazione dei binari Go

Compila i due eseguibili necessari (simulatore e client):

**Su macOS/Linux/WSL:**
```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o go_simulation/go_simulation go_simulation/mqtt_simulation.go
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o go_to_timescale/mqtt_to_timescale go_to_timescale/mqtt_to_timescale.go
```

**Su Windows (PowerShell):**
```powershell
$env:GOOS="linux"; $env:GOARCH="arm64"; $env:CGO_ENABLED="0"; go build -o go_simulation/go_simulation go_simulation/mqtt_simulation.go
$env:GOOS="linux"; $env:GOARCH="arm64"; $env:CGO_ENABLED="0"; go build -o go_to_timescale/mqtt_to_timescale go_to_timescale/mqtt_to_timescale.go
```

## 2. Avvio della VM Podman

**macOS:**
```bash
podman machine start
```

**Windows:**
```powershell
podman machine start
```

## 3. Creazione e avvio del Pod

Assicurati di esportare la variabile `PWD` per i path assoluti:

**macOS/Linux/WSL:**
```bash
export PWD=$(pwd)
envsubst < pod-preconfigured.yaml | podman play kube -
```

**Windows (PowerShell):**
```powershell
$env:PWD = (Get-Location).Path
envsubst < pod-preconfigured.yaml | podman play kube -
```

Oppure, per creare un pod parametrico con lo script:

**macOS/Linux/WSL:**
```bash
./createpod.sh palletizer1 1883 4000
```

**Windows (PowerShell):**
```powershell
bash createpod.sh palletizer1 1883 4000
```

## 4. Gestione dei Pod

Visualizza i pod attivi:

```bash
podman pod ps
```

Visualizza i container di un pod:

```bash
podman ps
```

Ferma, avvia o elimina un pod:

```bash
podman pod stop <nome>
podman pod start <nome>
podman pod rm -f <nome>
```

## 5. Accesso a Grafana

Apri il browser su:

```
http://localhost:4000
```

Credenziali di default:
- **Username:** admin
- **Password:** admin

---

Per pulire i dati di un pod:

**macOS/Linux/WSL:**
```bash
./clearpodfiles.sh palletizer1
```

**Windows (PowerShell):**
```powershell
bash clearpodfiles.sh palletizer1
```

---

**Nota:**  
Puoi creare più pod contemporaneamente cambiando nome e porte, ad esempio:

```bash
./createpod.sh palletizer2 1884 4001
```