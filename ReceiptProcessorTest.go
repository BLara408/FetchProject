package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculatePointsForReceipt(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Sample Retailer",
		PurchaseDate: "2024-03-03",
		PurchaseTime: "15:30",
		Total:        "100.00",
		Items: []Item{
			{ShortDescription: "Item 1", Price: "20.00"},
			{ShortDescription: "Item 2", Price: "30.00"},
		},
	}

	points := calculatePointsForReceipt(receipt)

	expectedPoints := 50 + 25 + 5 + 6 + 10

	if points != expectedPoints {
		t.Errorf("calculatePointsForReceipt() returned %d points; expected %d points", points, expectedPoints)
	}
}

func TestProcessReceipts(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Sample Retailer",
		PurchaseDate: "2024-03-03",
		PurchaseTime: "15:30",
		Total:        "100.00",
		Items: []Item{
			{ShortDescription: "Item 1", Price: "20.00"},
			{ShortDescription: "Item 2", Price: "30.00"},
		},
	}

	requestBody, err := json.Marshal(receipt)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/receipts/process", strings.NewReader(string(requestBody)))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	processReceipts(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("processReceipts() returned status code %d; expected %d", rr.Code, http.StatusOK)
	}

	var response ReceiptIDResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if response.ID == "" {
		t.Errorf("processReceipts() returned empty ID")
	}
}
