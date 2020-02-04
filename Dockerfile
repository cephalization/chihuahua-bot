# build stage
FROM golang as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# final stage
FROM scratch
COPY --from=builder /app/chihuahua-bot /app/
EXPOSE 8080
ENTRYPOINT ["/app/chihuahua-bot"]
