package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"errors"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	result := make(map[string][]string)
	lines := strings.Split(fileContent, "\n")

	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return nil, errors.New("empty file")
	}

	headers := strings.Split(lines[0], ",")
	for _, header := range headers {
		result[header] = []string{}
	}

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		values := strings.Split(line, ",")

		if len(values) != len(headers) {
			return nil, errors.New("invalid CSV data: number of columns doesn't match header")
		}

		for i, value := range values {
			result[headers[i]] = append(result[headers[i]], value)
		}
	}
	return result, nil
}
