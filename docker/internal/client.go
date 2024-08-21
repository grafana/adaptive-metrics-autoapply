package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	resp, err := c.makeNewRequest(http.MethodGet, "aggregations/rules/segments", nil, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var segments []Segment
	if err = json.NewDecoder(resp.Body).Decode(&segments); err != nil {
		return nil, err
	}

	return segments, nil
}

func (c *Client) FetchRecommendations(segment Segment, verbose bool) ([]Recommendation, error) {
	resp, err := c.makeNewRequest(http.MethodGet, "aggregations/recommendations", url.Values{
		"segment": []string{segment.Identifier},
		"verbose": []string{strconv.FormatBool(verbose)},
	}, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var recs []Recommendation
	if err = json.NewDecoder(resp.Body).Decode(&recs); err != nil {
		return nil, err
	}

	return recs, nil
}

func (c *Client) ValidateRules(rules []Recommendation) error {
	buf, err := json.Marshal(rules)
	if err != nil {
		return err
	}

	resp, err := c.makeNewRequest(http.MethodPost, "aggregations/check-rules", nil, nil, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBuf, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
		if err != nil {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		return fmt.Errorf("unexpected status code: %d with body %q", resp.StatusCode, bodyBuf)
	}

	return nil
}

func (c *Client) GetRules(segment Segment) ([]Recommendation, string, error) {
	resp, err := c.makeNewRequest(http.MethodGet, "aggregations/rules", url.Values{
		"segment": []string{segment.Identifier},
	}, nil, nil)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	etag := resp.Header.Get("ETag")
	var rules []Recommendation
	if err = json.NewDecoder(resp.Body).Decode(&rules); err != nil {
		return nil, "", err
	}

	return rules, etag, nil
}

func (c *Client) UpdateRules(segment Segment, etag string, rules []Recommendation) error {
	buf, err := json.Marshal(rules)
	if err != nil {
		return err
	}

	resp, err := c.makeNewRequest(http.MethodPost, "aggregations/rules", url.Values{
		"segment": []string{segment.Identifier},
	}, http.Header{"If-Match": []string{etag}}, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBuf, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
		if err != nil {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		return fmt.Errorf("unexpected status code: %d with body %q", resp.StatusCode, bodyBuf)
	}

	return nil
}

func (c *Client) makeNewRequest(method, subPath string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	p := fmt.Sprintf("%s/%s", c.apiURL, subPath)
	if queryParams != nil {
		p += "?" + queryParams.Encode()
	}

	req, err := http.NewRequest(method, p, body)
	if err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(http.Header)
	}
	req.Header = headers.Clone()
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	return c.httpClient.Do(req)
}
