FROM golang:1.23-bullseye

RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y build-essential wget curl zip lsb-release git ssh

RUN curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | BINDIR=/usr/local/bin sh

RUN arduino-cli core update-index
RUN arduino-cli core install arduino:avr

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY ./main.go ./
#COPY ./auth/*.go ./auth/
#COPY ./build/*.go ./build/
COPY ./database/*.go ./database/
#COPY ./parameter/*.go ./parameter/
COPY ./web/*.go ./web/
COPY ./common/*.go ./common/
# RUN go test -v ./...
RUN go build -mod=readonly -v -o server

# For local development environment only
COPY sketch-bridge-c8804059e16c.json ./

# EXPOSE 8088

ENTRYPOINT /app/server
