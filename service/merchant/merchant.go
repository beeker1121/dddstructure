package merchant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"dddstructure/proto"
	"dddstructure/storage"
)

// Service defines the merchant service.
type Service struct {
	s *storage.Storage
}

// New creates a new service.
func New() *Service {
	return &Service{}
}

// Create creates a new merchant.
func (s *Service) Create(m *proto.Merchant) (*proto.Merchant, error) {
	// Marshal the merchant to JSON.
	merchantJSON, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	// Call merchant microservice.
	fmt.Println("Calling merchant.Create microservice...")
	req, err := http.NewRequest("POST", "http://localhost:8082/rest/merchant", bytes.NewBuffer(merchantJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshal the merchant.
	var resm proto.Merchant
	if err := json.NewDecoder(resp.Body).Decode(&resm); err != nil {
		return nil, err
	}

	return &resm, nil
}

// GetByID gets a merchant by the given ID.
func (s *Service) GetByID(id uint) (*proto.Merchant, error) {
	// Call merchant microservice.
	fmt.Println("Calling merchant.GetByID microservice...")
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8082/rest/merchant/%d", id), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshal the merchant.
	var resm proto.Merchant
	if err := json.NewDecoder(resp.Body).Decode(&resm); err != nil {
		return nil, err
	}

	return &resm, nil
}
