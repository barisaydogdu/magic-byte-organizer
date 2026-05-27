package detector

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var allowedTypes = map[string][]string{
	"image/jpeg":      {".jpg", ".jpeg"},
	"image/png":       {".png"},
	"application/pdf": {".pdf"},
}

func DetectFileType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	fmt.Println("file.FD", file.Fd())

	defer file.Close()

	buf := make([]byte, 512)

	n, err := file.Read(buf)
	if err != io.EOF {
		return "", err
	}
	fmt.Printf("Read %d bytes: %x\n", n, buf[n:])

	return http.DetectContentType(buf), nil
}

func CheckExtAndType(filePath string, contentType string) (bool, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	validExtensions, exist := allowedTypes[contentType]

	if !exist {
		return false, fmt.Errorf("file type %s not allowed", contentType)
	}

	isMatched := false

	for _, v := range validExtensions {
		if ext == v {
			isMatched = true
			break
		}
	}

	return isMatched, nil
}

func tryExclusiveLock(filePath string) bool {
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0666)
	if err != nil {
		return false
	}

	defer file.Close()

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return false
	}

	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)

	return true
}

func WaitForFileLock(filePath string) bool {
	maxRetry := 60

	for i := 0; i < maxRetry; i++ {
		if tryExclusiveLock(filePath) {
			if CheckFileSize(filePath) {
				return true
			}
		}

		log.Printf("file is busy waiting %i/%d seconds\n", i, maxRetry)
		time.Sleep(time.Second * 2)
	}

	return false
}

func CheckFileSize(filePath string) bool {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}

	initialSize := initialStat.Size()

	time.Sleep(2 * time.Second)

	finalStat, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}

	finalSize := finalStat.Size()

	if initialSize == finalSize {
		log.Printf("file size is fixed: %d bytes", finalSize)
		return true
	}

	log.Printf("Dosya büyüyor: %d -> %d byte\n", initialSize, finalSize)
	return false
}
