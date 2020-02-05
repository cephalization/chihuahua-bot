# build stage
FROM golang as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

ENTRYPOINT ["/app/chihuahua-bot"]
