package detector

import (
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
	"text/plain":      {".txt", ".csv"},
}

func DetectFileType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	buf := make([]byte, 512)

	_, err = file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	
	return http.DetectContentType(buf), nil
}

func CheckExtAndType(filePath string, contentType string) (bool, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	validExtensions, exist := allowedTypes[contentType]

	if !exist {
		//return false, fmt.Errorf("file type %s not allowed", contentType)
		log.Println("Invalid content type:", contentType)
		return true, nil
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
		log.Println("open file error", err)
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
	retryCount := 0

	ticker := time.NewTicker(time.Second * 5)

	defer ticker.Stop()

	for i := 0; i < maxRetry; i++ {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Println("File does not exist")
			return false
		}

		if tryExclusiveLock(filePath) {
			if CheckFileSize(filePath) {
				return true
			}
		}

		retryCount++

		if retryCount >= maxRetry {
			return false
		}

		log.Printf("File is busy waiting %d/%d seconds\n", i, maxRetry)

		select {
		case <-ticker.C:
			continue
		}
	}

	return false
}

func CheckFileSize(filePath string) bool {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		log.Printf("File size check error: %s\n", err)
		return false
	}

	initialSize := initialStat.Size()

	time.Sleep(6 * time.Second)

	finalStat, err := os.Stat(filePath)
	if err != nil {
		log.Printf("File size check error: %s\n", err)
		return false
	}

	finalSize := finalStat.Size()

	if initialSize == finalSize {
		log.Printf("File size is fixed: %d bytes", finalSize)
		return true
	}

	log.Printf("File is growing: %d -> %d byte\n", initialSize, finalSize)
	return false
}
