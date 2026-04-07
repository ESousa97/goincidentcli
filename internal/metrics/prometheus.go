package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// PrometheusClient queries the Prometheus HTTP API.
type PrometheusClient struct {
	baseURL string
	client  *http.Client
}

// NewPrometheusClient creates a client targeting the given Prometheus base URL.
func NewPrometheusClient(baseURL string) *PrometheusClient {
	return &PrometheusClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

type apiResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []any             `json:"value"` // [unix_timestamp, string_value]
		} `json:"result"`
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

// QueryCurrent executes an instant PromQL query and returns the scalar string value.
// Returns "no data" when the query yields no results.
func (p *PrometheusClient) QueryCurrent(query string) (string, error) {
	endpoint := p.baseURL + "/api/v1/query"
	params := url.Values{"query": {query}}

	resp, err := p.client.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return "", fmt.Errorf("prometheus request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read prometheus response: %w", err)
	}

	var result apiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse prometheus response: %w", err)
	}

	if result.Status != "success" {
		errMsg := result.Error
		if errMsg == "" {
			errMsg = "unknown error"
		}
		return "", fmt.Errorf("prometheus error: %s", errMsg)
	}

	if len(result.Data.Result) == 0 {
		return "no data", nil
	}

	values := result.Data.Result[0].Value
	if len(values) < 2 {
		return "no value", nil
	}

	if v, ok := values[1].(string); ok {
		return v, nil
	}

	return fmt.Sprintf("%v", values[1]), nil
}
