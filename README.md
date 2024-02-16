#Receipt Processor Webservice

This is a Go application for processing receipts and calculating points based on certain criteria.

## Running with Docker
### Prerequisites
- Docker installed on your machine

1. **Build Docker Image:**
docker build -t fetch-project-app .

2. **Running the Docker Container:**
docker run -p 8080:8080 fetch-project-app

## Running with Go exectutable 

1. **Building Go executable:**
   go build -o fetch-project-app
   
2. **Running Go executable:**
   ./fetch-project-app

### Testing
1: Curl Commands

2: Postman

A unique UUID will be generated when POST http://localhost:8080/receipts/process is called.
Take this UUID and insert into http://localhost:8080/receipts/{id}/points

Example:
POST http://localhost:8080/receipts/process
Return
{"id":"94891dfa-a67f-4733-9405-8177beb96384"}

GET http://localhost:8080/receipts/94891dfa-a67f-4733-9405-8177beb96384/points
