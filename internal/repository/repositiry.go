package repository

type URLRepository interface {
	Set(url string) string
	GetById(id string) (string, error)
}
