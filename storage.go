package main

import (
	"encoding/json"
	"os"
)

// Storage is a generic struct that handles saving and loading data of type T to/from a JSON file
type Storage[T any] struct {
	FileName string
}

// NewStorage creates a new Storage instance for the given file name. It returns a pointer to the Storage instance.
func NewStorage[T any](fileName string) *Storage[T] {
	return &Storage[T]{FileName: fileName}
}

// Save serializes the provided data to JSON and writes it to the file. It uses indentation for readability. and sets file permissions to 0644.Returns an error if the operation fails.
func (s *Storage[T]) Save(data T) error {
	fileData, err := json.MarshalIndent(data, "", "  ") // serialize data to JSON with indentation
	if err != nil {
		return err
	}
	return os.WriteFile(s.FileName, fileData, 0644) // write JSON data to file
}

func (s *Storage[T]) Load(data *T) error {
	fileData, err := os.ReadFile(s.FileName) // read file data
	if err != nil {
		if os.IsNotExist(err) { // if file does not exist, return nil to indicate no data to load
			return nil
		}
		return err
	}
	return json.Unmarshal(fileData, data) // deserialize JSON data into the provided variable
}
