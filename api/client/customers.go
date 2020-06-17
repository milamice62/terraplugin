package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Customer struct {
	CustomerID string `json:"_id"`
	Name       string `json:"name"`
	IsGold     bool   `json:"isGold"`
	Phone      string `json:"phone"`
}

// GetAll Retrieves all of the Items from the server
func (c *Client) GetAllCustomers() (*map[string]Customer, error) {
	body, err := c.httpRequest("", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	customers := map[string]Customer{}
	err = json.NewDecoder(body).Decode(&customers)
	if err != nil {
		return nil, err
	}
	return &customers, nil
}

// GetItem gets an item with a specific name from the server
func (c *Client) GetCustomer(customerID string) (*Customer, error) {
	body, err := c.httpRequest(fmt.Sprintf("api/customers/%s", customerID), "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	customer := &Customer{}
	err = json.NewDecoder(body).Decode(customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

// create new genre
func (c *Client) NewCustomer(customer *Customer) (*io.ReadCloser, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(customer)
	if err != nil {
		return nil, err
	}
	body, err := c.httpRequest("api/customers", "POST", buf)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

// UpdateItem updates the values of an item
func (c *Client) UpdateCustomer(customer *Customer) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(customer)
	if err != nil {
		return err
	}
	_, err = c.httpRequest(fmt.Sprintf("api/customers/%s", customer.CustomerID), "PUT", buf)
	if err != nil {
		return err
	}
	return nil
}

// DeleteItem removes an item from the server
func (c *Client) DeleteCustomer(customerID string) error {
	_, err := c.httpRequest(fmt.Sprintf("api/customers/%s", customerID), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}
