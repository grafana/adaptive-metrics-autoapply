package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	httpClient *http.Client
	userAgent  string

	apiURL string
	apiKey string
}

func NewClient(httpClient *http.Client, userAgent, apiURL, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		userAgent:  userAgent,
		apiURL:     apiURL,
		apiKey:     apiKey,
	}
}

func (c *Client) FetchSegments() ([]Segment, error) {
	resp, err := c.makeNewRequest("aggregations/rules/segments", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var segments []Segment
	if err = json.NewDecoder(resp.Body).Decode(&segments); err != nil {
		return nil, err
	}

	return segments, nil
}

func (c *Client) FetchRecommendations(segment Segment, verbose bool) ([]Recommendation, error) {
	resp, err := c.makeNewRequest("aggregations/recommendations", url.Values{
		"segment": []string{segment.Identifier},
		"verbose": []string{strconv.FormatBool(verbose)},
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var recs []Recommendation
	if err = json.NewDecoder(resp.Body).Decode(&recs); err != nil {
		return nil, err
	}

	return recs, nil
}

func (c *Client) makeNewRequest(subPath string, queryParams url.Values) (*http.Response, error) {
	p := fmt.Sprintf("%s/%s", c.apiURL, subPath)
	if queryParams != nil {
		p += "?" + queryParams.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, p, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	return c.httpClient.Do(req)
}
