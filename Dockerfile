FROM golang:1.20-bullseye

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY eddilithium2jwt/eddilithium2jwt.go ./eddilithium2jwt/eddilithium2jwt.go
COPY cmd/server/server.go ./cmd/server/server.go

RUN go build -o /dilithiumwebdemo ./cmd/server

EXPOSE 1323

CMD [ "/dilithiumwebdemo"]
