package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spaceapegames/terraform-provider-example/api/server"
)

// Client holds all of the information required to connect to a server
type Client struct {
	hostname   string
	port       int
	authToken  string
	httpClient *http.Client
}

type Genre struct {
	Name string `json:"name"`
	ID   string `json:"_id"`
}

// NewClient returns a new client configured to communicate on a server with the
// given hostname and port and to send an Authorization Header with the value of
// token
func NewClient(hostname string, port int, token string) *Client {
	return &Client{
		hostname:   hostname,
		port:       port,
		authToken:  token,
		httpClient: &http.Client{},
	}
}

// GetAll Retrieves all of the Items from the server
func (c *Client) GetAllGenres() (*map[string]server.Item, error) {
	body, err := c.httpRequest("item", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	items := map[string]server.Item{}
	err = json.NewDecoder(body).Decode(&items)
	if err != nil {
		return nil, err
	}
	return &items, nil
}

// GetItem gets an item with a specific name from the server
func (c *Client) GetGenre(genreID string) (*Genre, error) {
	body, err := c.httpRequest(fmt.Sprintf("api/genres/%s", genreID), "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	genre := &Genre{}
	err = json.NewDecoder(body).Decode(genre)
	if err != nil {
		return nil, err
	}
	return genre, nil
}

// create new genre
func (c *Client) NewGenre(genre *Genre) (io.ReadCloser, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(genre)
	if err != nil {
		return nil, err
	}
	resBody, err := c.httpRequest("api/genres", "POST", buf)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}

// UpdateItem updates the values of an item
func (c *Client) UpdateGenre(genre *Genre) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(genre)
	if err != nil {
		return err
	}
	_, err = c.httpRequest(fmt.Sprintf("api/genres/%s", genre.Name), "PUT", buf)
	if err != nil {
		return err
	}
	return nil
}

// DeleteItem removes an item from the server
func (c *Client) DeleteGenre(genreID string) error {
	_, err := c.httpRequest(fmt.Sprintf("api/genres/%s", genreID), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) httpRequest(path, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-auth-token", c.authToken)
	switch method {
	case "GET":
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", resp.StatusCode, respBody.String())
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("%s:%v/%s", c.hostname, c.port, path)
}
