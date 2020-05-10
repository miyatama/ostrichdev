package ostrich

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"unsafe"
)

type FileAccesserInterface interface {
	ReadAll(filepath string) ([]string, error)
	WriteAll(filepath string, contents []string) error
	RemoveFile(filepath string) error
}

type FileAccesser struct {
}

// ReadAll is return content splited '\n' string
func (f *FileAccesser) ReadAll(filepath string) ([]string, error) {
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return []string{}, err
	}
	contentStr := *(*string)(unsafe.Pointer(&contents))
	splittedContent := strings.Split(contentStr, "\n")
	return splittedContent, nil

}

// WriteAll is write file
func (f *FileAccesser) WriteAll(filepath string, contents []string) error {
	byteContent := f.strings2Bytes(contents)
	err := ioutil.WriteFile(filepath, byteContent, 0644)
	if err != nil {
		return err
	}
	return nil

}

// RemoveFile is removing file.
func (f *FileAccesser) RemoveFile(filepath string) error {
	return os.Remove(filepath)
}

func (f *FileAccesser) strings2Bytes(texts []string) []byte {
	content := bytes.NewBuffer(make([]byte, 0, 1024)) //1K bytes capacity
	recode := "\n"
	for _, line := range texts {
		content.WriteString(line)
		content.WriteString(recode)
	}
	return content.Bytes()
}
