FROM golang:latest

WORKDIR /app

COPY main.go go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

CMD ["/docker-gs-ping"]

LABEL authors="ranv"
