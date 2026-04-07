package pagerduty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.pagerduty.com"
)

type Client struct {
	apiToken   string
	httpClient *http.Client
}

// OnCallUser represents a user currently on-call
type OnCallUser struct {
	User struct {
		ID    string `json:"id"`
		Name  string `json:"summary"`
		Email string `json:"email"`
	} `json:"user"`
	EscalationPolicy struct {
		ID   string `json:"id"`
		Name string `json:"summary"`
	} `json:"escalation_policy"`
	Schedule struct {
		ID   string `json:"id"`
		Name string `json:"summary"`
	} `json:"schedule"`
}

// Schedule represents a PagerDuty schedule
type Schedule struct {
	ID   string `json:"id"`
	Name string `json:"summary"`
}

// EscalationPolicy represents a PagerDuty escalation policy
type EscalationPolicy struct {
	ID   string `json:"id"`
	Name string `json:"summary"`
}

// OnCallResponse is the API response for on-call users
type OnCallResponse struct {
	OnCalls []OnCallUser `json:"oncalls"`
}

// SchedulesResponse is the API response for schedules
type SchedulesResponse struct {
	Schedules []Schedule `json:"schedules"`
}

// EscalationPoliciesResponse is the API response for escalation policies
type EscalationPoliciesResponse struct {
	EscalationPolicies []EscalationPolicy `json:"escalation_policies"`
}

func NewClient(apiToken string) *Client {
	return &Client{
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetOnCallUsers retrieves current on-call users for given escalation policy IDs
func (c *Client) GetOnCallUsers(ctx context.Context, escalationPolicyIDs []string) ([]OnCallUser, error) {
	url := BaseURL + "/oncalls?time_zone=UTC"
	for _, id := range escalationPolicyIDs {
		url += fmt.Sprintf("&escalation_policy_ids[]=%s", id)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.apiToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pagerduty api returned status %d: %s", resp.StatusCode, string(body))
	}

	var result OnCallResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.OnCalls, nil
}

// GetSchedules retrieves all schedules from PagerDuty
func (c *Client) GetSchedules(ctx context.Context) ([]Schedule, error) {
	url := BaseURL + "/schedules"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.apiToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pagerduty api returned status %d: %s", resp.StatusCode, string(body))
	}

	var result SchedulesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Schedules, nil
}

// GetEscalationPolicies retrieves all escalation policies from PagerDuty
func (c *Client) GetEscalationPolicies(ctx context.Context) ([]EscalationPolicy, error) {
	url := BaseURL + "/escalation_policies"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.apiToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pagerduty api returned status %d: %s", resp.StatusCode, string(body))
	}

	var result EscalationPoliciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.EscalationPolicies, nil
}

// CreateIncident creates a new incident in PagerDuty
func (c *Client) CreateIncident(ctx context.Context, serviceID, title, description string, urgency string) error {
	url := BaseURL + "/incidents"

	payload := map[string]interface{}{
		"incident": map[string]interface{}{
			"type":    "incident",
			"title":   title,
			"service": map[string]string{"id": serviceID, "type": "service_reference"},
			"urgency": urgency,
			"body":    map[string]string{"type": "incident_body", "details": description},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.apiToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("From", "sentinel-incident@system.local")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pagerduty api returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
