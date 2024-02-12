package main

import (
	"net/http"
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

}

func processReceptis(writer http.ResponseWriter, request *http.Request) {

}

func getPoints(writer http.ResponseWriter, rquest *http.Request) {

}

func calculatePointsByID(id string) int {

	return 0
}
