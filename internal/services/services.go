package services

import "context"

type ShortURLService interface {
	ShortenURL(context.Context, string, string) (string, error)
	PingDB(context.Context) error
}
