package detector

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
)

func detectFileType(filePath string) (string, error) {
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

func waitForFileLock(filePath string) bool {
	maxRetry := 60

	for i := 0; i < maxRetry; i++ {
		if tryExclusiveLock(filePath) {
			if checkFileSize(filePath) {
				return true
			}
		}

		log.Printf("file is busy waiting %i/%d seconds\n", i, maxRetry)
		time.Sleep(time.Second * 2)
	}

	return false
}

func checkFileSize(filePath string) bool {
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
