package restaurant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type CreateOrderRequest struct {
	OrderID string      `json:"order_id"`
	Items   []OrderItem `json:"items"`
}

type OrderItem struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
	Price      int64  `json:"price"`
	Quantity   int    `json:"quantity"`
}

type CreateOrderResponse struct {
	ExternalOrderID string `json:"external_order_id"`
	Status          string `json:"status"`
}

func (c *Client) CreateOrder(
	ctx context.Context,
	request CreateOrderRequest,
) (CreateOrderResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return CreateOrderResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/orders",
		bytes.NewReader(body),
	)
	if err != nil {
		return CreateOrderResponse{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return CreateOrderResponse{}, fmt.Errorf("send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return CreateOrderResponse{}, fmt.Errorf(
			"restaurant returned status %d",
			response.StatusCode,
		)
	}

	var result CreateOrderResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return CreateOrderResponse{}, fmt.Errorf(
			"decode response: %w",
			err,
		)
	}

	return result, nil
}
