package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"dddstructure/proto"
)

// Service defines the user service.
type Service struct{}

// New creates a new service.
func New() *Service {
	return &Service{}
}

// Create creates a new user.
func (s *Service) Create(u *proto.User) (*proto.User, error) {
	// Marshal the user to JSON.
	userJSON, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	// Call user microservice.
	fmt.Println("Calling user.Create microservice...")
	req, err := http.NewRequest("POST", "http://localhost:8081/rest/user", bytes.NewBuffer(userJSON))
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

	// Unmarshal the user.
	var resu proto.User
	if err := json.NewDecoder(resp.Body).Decode(&resu); err != nil {
		return nil, err
	}

	return &resu, nil
}

// GetByID gets a user by the given ID.
func (s *Service) GetByID(id uint) (*proto.User, error) {
	// Call user microservice.
	fmt.Println("Calling user.GetByID microservice...")
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8081/rest/user/%d", id), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshal the user.
	var resu proto.User
	if err := json.NewDecoder(resp.Body).Decode(&resu); err != nil {
		return nil, err
	}

	return &resu, nil
}
