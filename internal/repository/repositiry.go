package repository

type IRepository interface {
	Set(url string) string
	GetByID(id string) (string, error)
}
