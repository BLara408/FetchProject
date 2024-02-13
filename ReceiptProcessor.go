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
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", processReceipts).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))

}

func processReceipts(writer http.ResponseWriter, request *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(request.Body).Decode(&receipt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	id := generateReceiptUUID()

	receipts[id] = receipt

	response := ReceiptIDResponse{ID: id}
	json.NewEncoder(writer).Encode(response)

}

func generateReceiptUUID() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func getPoints(writer http.ResponseWriter, request *http.Request) {
	parts := strings.Split(request.URL.Path, "/")

	id := parts[len(parts)-2]

	receipt, found := receipts[id]
	if !found {
		http.Error(writer, "No receipt found for that id", http.StatusNotFound)
		return
	}
	points := calculatePointsForReceipt(receipt)

	response := PointsResponse{Points: points}

	json.NewEncoder(writer).Encode(response)

}

func calculatePointsForReceipt(receipt Receipt) int {
	points := 0
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points++
		}
	}

	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(total, 1) == 0 {
		points += 50
	}

	if err == nil && math.Mod(total, 0.25) == 0 {
		points += 25
	}

	points += len(receipt.Items) / 2 * 5

	for _, item := range receipt.Items {
		trimLength := len(strings.TrimSpace(item.ShortDescription))

		if trimLength%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)

			if err != nil {

				continue
			}

			points += int(math.Ceil(price * 0.2))

		}
	}

	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDate.Day()%2 != 0 {
		points += 6

	}

	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil && purchaseTime.After(time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)) && purchaseTime.Before(time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)) {
		points += 10
	}

	return points
}
