package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yv3000/siig/rules"
)

const (
	baseURL    = "https://raw.githubusercontent.com/yv3000/siig-rules/main/db/"
	timeout    = 5 * time.Second
	userAgent  = "siig/2.0"
)

// FetchStack fetches a single stack JSON by ID from siig-rules.
// Returns nil, nil if the stack doesn't exist (404).
// Returns nil, err on network/parse errors.
func FetchStack(stackID string) (*rules.StackRule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	url := baseURL + stackID + ".json"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch failed for %s: %w", stackID, err)
	}
	defer resp.Body.Close()

	// 404 means the stack doesn't exist in siig-rules — not an error
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, stackID)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response for %s: %w", stackID, err)
	}

	var rule rules.StackRule
	if err := json.Unmarshal(body, &rule); err != nil {
		return nil, fmt.Errorf("cannot parse JSON for %s: %w", stackID, err)
	}

	return &rule, nil
}
