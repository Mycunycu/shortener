package services

import "context"

type ShortURLService interface {
	ShortenURL(context.Context, string, string) (string, error)
	ExpandURL(context.Context, string) (string, error)
	PingDB(context.Context) error
}
