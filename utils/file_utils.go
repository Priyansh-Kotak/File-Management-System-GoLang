package utils

import (
	"io"
	"os"
)

// SaveFileLocally saves a file to the local storage
func SaveFileLocally(file io.Reader, fileName string) error {
	out, err := os.Create("./uploads/" + fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	return nil
}