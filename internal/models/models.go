package models

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponce struct {
	Result string `json:"result"`
}

type UserUrlItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
