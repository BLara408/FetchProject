package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ReceiptIDResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

var receipts = make(map[string]Receipt)

func main() {
	// Attaching the router to the correct functions.
	r := mux.NewRouter()

	// Serve static files from the FetchProject directory
	fs := http.FileServer(http.Dir(`./static`))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", loggingMiddleware(fs)))

	r.HandleFunc("/receipts/process", processReceipts).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")
	port := ":8080"
	log.Printf("Starting server on port %s\n", port)
	// Running on port 8080
	log.Fatal(http.ListenAndServe(":8080", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request URL:", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func processReceipts(writer http.ResponseWriter, request *http.Request) {
	// Decode the request body and check for required fields
	var receipt Receipt
	err := json.NewDecoder(request.Body).Decode(&receipt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if required fields are missing
	if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || receipt.Total == "" || len(receipt.Items) == 0 {
		http.Error(writer, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Call to generate unique UUID
	id := generateReceiptUUID()

	// Map the receipt to ID
	receipts[id] = receipt

	response := ReceiptIDResponse{ID: id}
	json.NewEncoder(writer).Encode(response)
}

func generateReceiptUUID() string {
	//Generates and returns a unique UUID
	newUUID := uuid.New()
	return newUUID.String()
}

func getPoints(writer http.ResponseWriter, request *http.Request) {
	//Splits the path into parts
	parts := strings.Split(request.URL.Path, "/")

	id := parts[len(parts)-2]

	//If a receipt is not found throws 404 Not found
	receipt, found := receipts[id]
	if !found {
		http.Error(writer, "No receipt found for that id", http.StatusNotFound)
		return
	}

	//Call to calculate the points for receipt
	points := calculatePointsForReceipt(receipt)

	response := PointsResponse{Points: points}

	json.NewEncoder(writer).Encode(response)
}

func calculatePointsForReceipt(receipt Receipt) int {
	points := countAlphaNumericChars(receipt.Retailer) +
		calculateRoundDollarPoints(receipt.Total) +
		calculateQuarterPoints(receipt.Total) +
		calculateItemPoints(receipt.Items) +
		calculateItemDescriptionPoints(receipt.Items) +
		calculateOddDayPoints(receipt.PurchaseDate) +
		calculateTimeRangePoints(receipt.PurchaseTime)

	return points
}
func countAlphaNumericChars(s string) int {
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

func calculateRoundDollarPoints(total string) int {
	points := 0
	if totalFloat, err := strconv.ParseFloat(total, 64); err == nil && math.Mod(totalFloat, 1) == 0 {
		points += 50
	}
	return points
}

func calculateQuarterPoints(total string) int {
	points := 0
	if totalFloat, err := strconv.ParseFloat(total, 64); err == nil && math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}
	return points
}

func calculateItemPoints(items []Item) int {
	return len(items) / 2 * 5
}

func calculateItemDescriptionPoints(items []Item) int {
	points := 0
	for _, item := range items {
		trimLength := len(strings.TrimSpace(item.ShortDescription))
		if trimLength%3 == 0 {
			if priceFloat, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int(math.Ceil(priceFloat * 0.2))
			}
		}
	}
	return points
}

func calculateOddDayPoints(purchaseDate string) int {
	points := 0
	if date, err := time.Parse("2006-01-02", purchaseDate); err == nil && date.Day()%2 != 0 {
		points += 6
	}
	return points
}

func calculateTimeRangePoints(purchaseTime string) int {
	points := 0
	if time, err := time.Parse("15:04", purchaseTime); err == nil {
		hour := time.Hour()
		minute := time.Minute()
		if (hour > 14 || (hour == 14 && minute >= 0)) && (hour < 16 || (hour == 16 && minute == 0)) {
			points += 10
		}
	}
	return points
}
