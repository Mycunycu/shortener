package services

import "context"

type ShortURLService interface {
	PingDB(context.Context) error
}
