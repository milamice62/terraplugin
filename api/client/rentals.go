package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Rental struct {
	RentalID string    `json:"_id"`
	Customer *Customer `json:"customer"`
	Movie    *Movie    `json:"movie"`
	DateOut  string    `json:"dateOut"`
}

type RentalID struct {
	MovieID    string `json:"movieId"`
	CustomerID string `json:"customerId"`
}

// GetAll Retrieves all of the Items from the server
func (c *Client) GetAllRentals() (*map[string]Rental, error) {
	body, err := c.httpRequest("", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	rental := map[string]Rental{}
	err = json.NewDecoder(body).Decode(&rental)
	if err != nil {
		return nil, err
	}
	return &rental, nil
}

// GetItem gets an item with a specific name from the server
func (c *Client) GetRental(rentalID string) (*Rental, error) {
	body, err := c.httpRequest(fmt.Sprintf("api/rentals/%s", rentalID), "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	rental := &Rental{}
	err = json.NewDecoder(body).Decode(rental)
	if err != nil {
		return nil, err
	}
	return rental, nil
}

// create new genre
func (c *Client) NewRental(rentalID *RentalID) (*io.ReadCloser, error) {
	buf := bytes.Buffer{}

	err := json.NewEncoder(&buf).Encode(&rentalID)
	if err != nil {
		return nil, err
	}
	body, err := c.httpRequest("api/rentals", "POST", buf)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

// UpdateItem updates the values of an item
func (c *Client) UpdateRental(rental *Rental) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(rental)
	if err != nil {
		return err
	}
	_, err = c.httpRequest(fmt.Sprintf("api/rentals/%s", rental.RentalID), "PUT", buf)
	if err != nil {
		return err
	}
	return nil
}

// DeleteItem removes an item from the server
func (c *Client) DeleteRental(rentalID string) error {
	_, err := c.httpRequest(fmt.Sprintf("api/rentals/%s", rentalID), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}
