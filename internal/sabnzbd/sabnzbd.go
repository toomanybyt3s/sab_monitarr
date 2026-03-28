package sabnzbd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Status represents the top-level response from the SABnzbd API.
type Status struct {
	Status string `json:"status"`
	Queue  Queue  `json:"queue"`
}

// Queue contains download queue information.
type Queue struct {
	Status   string `json:"status"`
	Speed    string `json:"speed"`
	SizeLeft string `json:"sizeleft"`
	TimeLeft string `json:"timeleft"`
	Slots    []Slot `json:"slots"`
}

// Slot represents a single item in the download queue.
type Slot struct {
	Filename   string `json:"filename"`
	Status     string `json:"status"`
	SizeLeft   string `json:"sizeleft"`
	Percentage string `json:"percentage"`
	TimeLeft   string `json:"timeleft"`
}

// Client is a reusable HTTP client with a sensible timeout.
var Client = &http.Client{
	Timeout: 5 * time.Second,
}

// FetchStatus calls the SABnzbd queue API and returns the parsed status.
// The API key is passed as a query parameter (SABnzbd requirement).
func FetchStatus(sabnzbdURL, apiKey string, debug bool) (*Status, error) {
	apiURL := fmt.Sprintf("%s/api?output=json&mode=queue&apikey=%s", sabnzbdURL, apiKey)

	if debug {
		safeURL := strings.Replace(apiURL, apiKey, "[REDACTED]", 1)
		fmt.Printf("[DEBUG] Requesting SABnzbd API: %s\n", safeURL)
	}

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build API request: %w", err)
	}

	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response body: %w", err)
	}

	if debug {
		fmt.Printf("[DEBUG] API response received, status: %s, length: %d bytes\n", resp.Status, len(bodyBytes))
	}

	var status Status
	if err = json.Unmarshal(bodyBytes, &status); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return &status, nil
}
