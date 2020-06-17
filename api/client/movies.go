package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Movie struct {
	MovieID string  `json:"_id"`
	Title   string  `json:"title"`
	Genre   *Genre  `json:"genre"`
	Stock   int     `json:"numberInStock"`
	Rate    float64 `json:"dailyRentalRate"`
}

// GetAll Retrieves all of the Items from the server
func (c *Client) GetAllMovies() (*map[string]Movie, error) {
	body, err := c.httpRequest("", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	movies := map[string]Movie{}
	err = json.NewDecoder(body).Decode(&movies)
	if err != nil {
		return nil, err
	}
	return &movies, nil
}

// GetItem gets an item with a specific name from the server
func (c *Client) GetMovie(movieID string) (*Movie, error) {
	body, err := c.httpRequest(fmt.Sprintf("api/movies/%s", movieID), "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	movie := &Movie{}
	err = json.NewDecoder(body).Decode(movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

// create new genre
func (c *Client) NewMovie(movie *Movie) (*io.ReadCloser, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(movie)
	if err != nil {
		return nil, err
	}
	body, err := c.httpRequest("api/movies", "POST", buf)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

// UpdateItem updates the values of an item
func (c *Client) UpdateMovie(movie *Movie) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(movie)
	if err != nil {
		return err
	}
	_, err = c.httpRequest(fmt.Sprintf("api/movies/%s", movie.Title), "PUT", buf)
	if err != nil {
		return err
	}
	return nil
}

// DeleteItem removes an item from the server
func (c *Client) DeleteMovie(movieID string) error {
	_, err := c.httpRequest(fmt.Sprintf("api/movies/%s", movieID), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}
