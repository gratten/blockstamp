# Use an official Go image as the base
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download


# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main ./app
# RUN chmod +x ./app

# Set environment variables (override at runtime)
# ENV BITCOIN_RPC_HOST=http://bitcoind:18443
# ENV BITCOIN_RPC_USER=your_rpc_user
# ENV BITCOIN_RPC_PASSWORD=your_rpc_password

# Run the app
CMD ["./main"]