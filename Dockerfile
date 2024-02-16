FROM golang:1.17-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ReceiptProcessor.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fetch-project-app .

FROM alpine:latest  

WORKDIR /root/

# Copy the Go application
COPY --from=builder /app/fetch-project-app .

# Copy the static files
COPY  ./static ./static

EXPOSE 8080

CMD ["./fetch-project-app"]
