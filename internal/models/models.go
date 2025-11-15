package models

import "time"

// Video represents a video in the library with its indexing status and metadata
type Video struct {
	VideoID       string        `json:"video_id"`
	URL           string        `json:"url"`
	IndexingStatus string       `json:"indexing_status"` // indexed, processing, failed
	Metadata      VideoMetadata `json:"metadata"`
	IndexingType  string        `json:"indexing_type"`
}

// VideoMetadata contains detailed information about a video
type VideoMetadata struct {
	Width                    int     `json:"width"`
	Height                   int     `json:"height"`
	AvgFPS                   float64 `json:"avg_fps"`
	VideoName                string  `json:"video_name"`
	Title                    string  `json:"title"`
	VideoStartTimestampUTCMs *int64  `json:"video_start_timestamp_utc_ms"`
	Duration                 float64 `json:"duration"`
	Thumbnail                string  `json:"thumbnail"`
	Description              string  `json:"description"`
	Source                   string  `json:"source"`
}

// QueryHistory represents a saved question/answer pair from the history
type QueryHistory struct {
	ID         int       `json:"id"`
	VideoID    string    `json:"video_id"`
	VideoTitle string    `json:"video_title"`
	Question   string    `json:"question"`
	Answer     string    `json:"answer"`
	Error      *string   `json:"error"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// ChatMessage represents a message in the chat API request
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// QARequest represents the request body for the QA chat API
type QARequest struct {
	VideoID  string        `json:"video_id"`
	Messages []ChatMessage `json:"messages"`
}

// QAResponse represents the response from the QA chat API
type QAResponse struct {
	ChatResponse string  `json:"chat_response"`
	SystemMessage *string `json:"system_message"`
	Error        *string `json:"error"`
	Status       string  `json:"status"`
	DebugChunks  *string `json:"debug_chunks"`
	DebugPredictedStartTime string `json:"debug_predicted_start_time"`
	DebugPredictedEndTime   string `json:"debug_predicted_end_time"`
}

// VideosGetRequest represents the request to get video information
type VideosGetRequest struct {
	VideoIDs []string `json:"video_ids,omitempty"`
}

// VideosGetResponse represents the response from the videos get API
type VideosGetResponse struct {
	Results []Video `json:"results"`
}
