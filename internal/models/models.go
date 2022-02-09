package models

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponce struct {
	Result string `json:"result"`
}

type ShortenEty struct {
	UserID      string `json:"user_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type UserHistoryItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
