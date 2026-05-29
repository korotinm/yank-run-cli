package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client is a thin typed wrapper over net/http for the yank-run backend.
type Client struct {
	BaseURL string
	HTTP    *http.Client
	UA      string
}

// New builds a Client with the given base URL, timeout, and version-derived UA.
func New(baseURL string, timeout time.Duration, version string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{Timeout: timeout},
		UA:      "yank-cli/" + version,
	}
}

// Create posts a new snippet. Returns CreateResp with New=true on HTTP 201.
func (c *Client) Create(req CreateReq) (*CreateResp, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal create request: %w", err)
	}
	httpReq, err := http.NewRequest("POST", c.BaseURL+"/api/snippets", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", c.UA)
	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, decodeError(resp)
	}
	var out CreateResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode create response: %w", err)
	}
	out.New = resp.StatusCode == 201
	return &out, nil
}

// Get fetches a snippet by id.
func (c *Client) Get(id string) (*Snippet, error) {
	httpReq, err := http.NewRequest("GET", c.BaseURL+"/api/snippets/"+id, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("User-Agent", c.UA)
	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, decodeError(resp)
	}
	var out Snippet
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode snippet: %w", err)
	}
	return &out, nil
}

// Search runs a full-text query and returns the raw hit list.
func (c *Client) Search(q string, limit int) ([]Hit, error) {
	u, err := url.Parse(c.BaseURL + "/api/snippets")
	if err != nil {
		return nil, err
	}
	qs := u.Query()
	qs.Set("q", q)
	if limit > 0 {
		qs.Set("limit", strconv.Itoa(limit))
	}
	u.RawQuery = qs.Encode()

	httpReq, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("User-Agent", c.UA)
	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, decodeError(resp)
	}
	var out []Hit
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}
	return out, nil
}

func decodeError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	var e struct {
		Error string `json:"error"`
	}
	_ = json.Unmarshal(body, &e)
	if e.Error == "" {
		e.Error = resp.Status
	}
	return &APIError{Status: resp.StatusCode, Message: e.Error}
}
