package models

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponce struct {
	Result string `json:"result"`
}

type ShortenEty struct {
	UserID      string `json:"user_id"`
	ShortID     string `json:"short_id"`
	OriginalURL string `json:"original_url"`
	Deleted     bool   `json:"deleted"`
}

type ShortenItem struct {
	ShortID     string `json:"short_id"`
	OriginalURL string `json:"original_url"`
}

type UserHistoryItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ShortenBatchRequest []BatchItemRequest

type BatchItemRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	shortID       string
}

type BatchItemResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
