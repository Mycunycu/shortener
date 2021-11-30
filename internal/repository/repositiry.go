package repository

type URLRepository interface {
	Set(url string) string
	GetByID(id string) (string, error)
}
