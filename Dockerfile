FROM golang:1.20

# Set destination for COPY
WORKDIR /banner_rotation

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download


COPY ./ ./
# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o rotation

EXPOSE 8081

# Run
CMD ["./rotation"]
