package models

type ShortenRequest struct {
	Url string `json:"url"`
}

type ShortenResponce struct {
	Result string `json:"result"`
}
