FROM golang:1.22-bookworm as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY main.go ./
COPY . ./

# Build the binary.
RUN CGO_ENABLED=0 go build -v -o go-svc-restart

# Use the official distroless image for a lean production container.
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static-debian12

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/go-svc-restart /

# Run the web service on container startup.
CMD ["/go-svc-restart"]
