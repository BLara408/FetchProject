package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
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
	http.HandleFunc("receipts/process", processReceipts)
	http.HandleFunc("receipts/{id}/points", getPoints)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func processReceipts(writer http.ResponseWriter, request *http.Request) {
	var receipt Receipt
	error := json.NewDecoder(request.Body).Decode(&receipts)
	if error != nil {
		http.Error(writer, error.Error(), http.StatusBadRequest)
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
		http.Error(writer, "Receipt not found", http.StatusNotFound)
		return
	}

	points := calculatePointsForReceipt(receipt)

	response := PointsResponse{Points: points}

	json.NewEncoder(writer).Encode(response)

}

func calculatePointsForReceipt(receipt Receipt) int {
	points := 0
	retailName := strings.ReplaceAll(receipt.Retailer, " ", "")
	points += len(retailName)

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
			if err == nil {
				points += int(math.Ceil(price * 0.2))
			}

		}

	}

	return 0
}
