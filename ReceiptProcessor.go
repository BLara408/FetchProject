package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string `json:"retailer`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items`
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

func main() {
	http.HandleFunc("receipts/process", processReceipts)
	http.HandleFunc("receipts", getPoints)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func processReceipts(writer http.ResponseWriter, request *http.Request) {
	var receipts Receipt
	error := json.NewDecoder(request.Body).Decode(&receipts)
	if error != nil {
		http.Error(writer, error.Error(), http.StatusBadRequest)
	}

	id := generateReceiptUUID()

	response := ReceiptIDResponse{ID: id}
	json.NewEncoder(writer).Encode(response)

}

func generateReceiptUUID() string {
	newUUID := uuid.New()
	return newUUID.String()
}

func getPoints(writer http.ResponseWriter, rquest *http.Request) {

}

func calculatePointsByID(id string) int {

	return 0
}
