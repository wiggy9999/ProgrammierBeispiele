FROM golang:1.19-alpine

WORKDIR /app

# Abhängigkeiten kopieren und installieren
COPY go.mod go.sum ./
RUN go mod download

# Quellcode kopieren
COPY . .

# Anwendung bauen
RUN go build -o main .

# Server starten
CMD ["./main"]
