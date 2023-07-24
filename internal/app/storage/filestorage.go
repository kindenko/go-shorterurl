package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/kindenko/go-shorterurl/internal/app/utils"
)

type FileStorage struct {
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

type File struct {
	path string
}

func InitFileDB(fileStoragePath string) *File {
	return &File{
		path: fileStoragePath,
	}
}

func (f *File) Save(fullURL string) (string, error) {
	var fs FileStorage

	fs.Original = fullURL
	fs.Short = utils.RandString()

	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return "", err
	}

	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(fs)
	return fs.Short, err
}

func (f *File) Get(shortURL string) (string, error) {

	file, err := os.OpenFile(f.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	data := make(map[string]string)

	for scanner.Scan() {
		var d FileStorage
		err = json.Unmarshal(scanner.Bytes(), &d)
		if err != nil {
			log.Println(err)
		}

		data[d.Short] = d.Original
	}
	original := data[shortURL]

	return original, nil

}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func SaveToFile(fs *FileStorage, fileName string) error {

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(fs)
	return err
}

func LoadFromFile(fileName string) (map[string]string, error) {

	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	data := make(map[string]string)

	for scanner.Scan() {
		var d FileStorage
		err = json.Unmarshal(scanner.Bytes(), &d)
		if err != nil {
			log.Println(err)
		}

		data[d.Short] = d.Original
	}
	return data, nil
}
