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
### Point Calculation:
   - One point for every alphanumeric character in the retailer name.
   - 50 points if the total is a round dollar amount with no cents.
   - 25 points if the total is a multiple of 0.25.
   - 5 points for every two items on the receipt.
   - If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
   - 6 points if the day in the purchase date is odd.
   - 10 points if the time of purchase is after 2:00pm and before 4:00pm.
