package repository

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type IStorage interface {
	WriteData(data string) error
	ReadAll() map[string]string
}

type Storage struct {
	*os.File
}

func NewStorage(path string) (*Storage, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &Storage{file}, nil
}

func (s *Storage) WriteData(data string) error {
	_, err := s.Write([]byte(data))
	return err
}

func (s *Storage) ReadAll() map[string]string {
	result := make(map[string]string)

	scanner := bufio.NewScanner(s)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ln := strings.Split(scanner.Text(), "\n")
		splited := strings.Split(ln[0], "-")

		fmt.Println(splited)

		result[splited[0]] = splited[1]
	}

	return result
}
