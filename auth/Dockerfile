FROM golang:1.23.4-bullseye

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app

EXPOSE 8080

CMD ["/app"]
