FROM golang:1.24

WORKDIR /usr/src/app

# Install Reflex
RUN go install github.com/cespare/reflex@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o /usr/local/bin/app .

EXPOSE 8080

# Use Reflex to watch for changes and restart the app
CMD ["reflex", "-r", "\\.go$", "--", "/usr/local/bin/app"]

