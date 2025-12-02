package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/fboucher/be-my-eyes/internal/models"
)

const (
	// Base URL for the Reka Vision API
	baseURL = "https://vision-agent.api.reka.ai"
)

// Client represents an API client for the Reka Vision API
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// DoRawRequest allows custom API calls for endpoints not covered by typed methods
func (c *Client) DoRawRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// UploadVideo uploads a video with multipart/form-data format
func (c *Client) UploadVideo(videoName, videoURL string, index bool) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add form fields
	if err := writer.WriteField("index", fmt.Sprintf("%t", index)); err != nil {
		return nil, fmt.Errorf("failed to write index field: %w", err)
	}
	if err := writer.WriteField("video_name", videoName); err != nil {
		return nil, fmt.Errorf("failed to write video_name field: %w", err)
	}
	if err := writer.WriteField("video_url", videoURL); err != nil {
		return nil, fmt.Errorf("failed to write video_url field: %w", err)
	}

	// Close the writer to finalize the multipart message
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", baseURL+"/videos/upload", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// NewClient creates a new API client with the given API key
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request with the API key header
func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetVideos retrieves information about one or more videos by their IDs
func (c *Client) GetVideos(videoIDs []string) (*models.VideosGetResponse, error) {
	req := models.VideosGetRequest{
		VideoIDs: videoIDs,
	}

	respBody, err := c.doRequest("POST", "/videos/get", req)
	if err != nil {
		return nil, err
	}

	var response models.VideosGetResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetAllVideos retrieves information about all videos (calls API without video_ids)
func (c *Client) GetAllVideos() (*models.VideosGetResponse, error) {
	// Call with empty request to get all videos
	req := models.VideosGetRequest{}

	respBody, err := c.doRequest("POST", "/videos/get", req)
	if err != nil {
		return nil, err
	}

	var response models.VideosGetResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// AskQuestion sends a question about a video to the API and returns the response
func (c *Client) AskQuestion(videoID, question string) (*models.QAResponse, error) {
	req := models.QARequest{
		VideoID: videoID,
		Messages: []models.ChatMessage{
			{
				Role:    "user",
				Content: question,
			},
		},
	}

	respBody, err := c.doRequest("POST", "/qa/chat", req)
	if err != nil {
		return nil, err
	}

	var response models.QAResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
