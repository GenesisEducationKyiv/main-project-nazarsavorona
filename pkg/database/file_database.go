package database

import (
	"bufio"
	"log"
	"os"
)

type FileDatabase struct {
	file string
}

func createFolderIfNotExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewFileDatabase(folder, file string) *FileDatabase {
	err := createFolderIfNotExists(folder)
	if err != nil {
		return nil
	}

	_, err = os.Create(file)
	if err != nil {
		return nil
	}

	return &FileDatabase{
		file: file,
	}
}

func (f *FileDatabase) AddNewEmail(email string) error {
	file, err := os.OpenFile(f.file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(file)

	if _, err = file.WriteString(email + "\n"); err != nil {
		panic(err)
	}

	return nil
}

func (f *FileDatabase) GetEmails() ([]string, error) {
	file, err := os.Open(f.file)

	if err != nil {
		return []string{}, err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(file)

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
