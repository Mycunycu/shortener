package repository

type Repositorier interface {
	Set(url string) string
	GetByID(id string) (string, error)
	WriteData(data string) error
	ReadAllData() map[string]string
}
