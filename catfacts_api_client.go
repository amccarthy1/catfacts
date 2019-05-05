package catfacts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const defaultURLString = "https://catfact.ninja"
const listBreedsEndpoint = "breeds"
const listFactsEndpoint = "facts"
const randomFactEndpoint = "fact"

// NewClientWithStringURL returns a new facts client using the given string as
// a base URL. If the given value is not a valid URL, returns an error.
func NewClientWithStringURL(urlString string) (*Client, error) {
	baseURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	return NewClientWithURL(baseURL), nil
}

// NewClientWithURL returns a new facts client using the given URL object as a
// base URL
func NewClientWithURL(baseURL *url.URL) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		pageSize:   10,
	}
}

// NewClient returns a new catfacts client using the default url.
func NewClient() *Client {
	client, err := NewClientWithStringURL(defaultURLString)
	if err != nil {
		panic(err)
	}
	return client
}

// Client is a catfacts API client
type Client struct {
	baseURL    *url.URL
	httpClient httpGetClient
	pageSize   int
}

type httpGetClient interface {
	Get(url string) (*http.Response, error)
}

// WithPageSize sets the page size of the API client. Note that this is not
// a streaming client, and smaller page sizes (<100) usually result in slower
// API calls due to round-trip request latency.
func (c *Client) WithPageSize(pageSize int) *Client {
	c.pageSize = pageSize
	return c
}

// WithHTTPClient sets the HTTP client of the API client. This can be useful
// for mocking or using a custom HTTP client implementation.
func (c *Client) WithHTTPClient(httpClient httpGetClient) *Client {
	c.httpClient = httpClient
	return c
}

// GetRandomFact returns a random fact from the /fact endpoint
func (c *Client) GetRandomFact() (*Fact, error) {
	u := fmt.Sprintf("%s/%s", c.baseURL.String(), randomFactEndpoint)
	res, err := c.httpClient.Get(u)

	var parsed Fact
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}

// ListAllBreeds returns a list of Breeds from the /breeds endpoint
func (c *Client) ListAllBreeds() ([]Breed, error) {
	u := fmt.Sprintf("%s/%s", c.baseURL.String(), listBreedsEndpoint)
	var allBreeds []Breed

	err := c.paginatedGet(u, func(page *json.Decoder) (*pagination, error) {
		var parsed breeds
		err := page.Decode(&parsed)
		if err != nil {
			return nil, err
		}
		allBreeds = append(allBreeds, parsed.Data...)
		return &parsed.pagination, nil
	})
	if err != nil {
		return nil, err
	}
	return allBreeds, nil
}

// ListAllFacts returns a list of Facts from the /facts endpoint
func (c *Client) ListAllFacts() ([]Fact, error) {
	u := fmt.Sprintf("%s/%s", c.baseURL.String(), listFactsEndpoint)
	var allFacts []Fact

	err := c.paginatedGet(u, func(page *json.Decoder) (*pagination, error) {
		var parsed facts
		err := page.Decode(&parsed)
		if err != nil {
			return nil, err
		}
		allFacts = append(allFacts, parsed.Data...)
		return &parsed.pagination, nil
	})
	if err != nil {
		return nil, err
	}
	return allFacts, nil
}

type pageHandler func(page *json.Decoder) (*pagination, error)

func (c *Client) paginatedGet(urlString string, handle pageHandler) error {
	currentPage := 0
	lastPage := 1
	for currentPage != lastPage {
		currentPage++
		res, err := c.httpClient.Get(fmt.Sprintf("%s?limit=%d&page=%d", urlString, c.pageSize, currentPage))
		if err != nil {
			return err
		}
		decoder := json.NewDecoder(res.Body)
		pagination, err := handle(decoder)
		if err != nil {
			return err
		}
		currentPage = pagination.CurrentPage
		lastPage = pagination.LastPage
	}
	return nil
}
