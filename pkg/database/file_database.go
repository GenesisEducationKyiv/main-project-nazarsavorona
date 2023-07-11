package database

import (
	"bufio"
	"io"
	"os"
)

type FileDatabase struct {
	file *os.File
}

func NewFileDatabase(file *os.File) *FileDatabase {
	return &FileDatabase{
		file: file,
	}
}

func (f *FileDatabase) AddEmail(email string) error {
	_, err := f.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	_, err = f.file.WriteString(email + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (f *FileDatabase) Emails() ([]string, error) {
	_, err := f.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f.file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}
