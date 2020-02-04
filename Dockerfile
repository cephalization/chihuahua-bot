# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/cephalization/chihuahua-bot

# Build the chibot command inside the container.
RUN go install github.com/cephalization/chihuahua-bot

# Run the chihuahua-bot command by default when the container starts.
ENTRYPOINT /go/bin/chihuahua-bot

# Document that the service listens on port 8080.
EXPOSE 8080