package mover

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func moveFile(sourcePath, destPath string) {
	fileName := filepath.Base(sourcePath)

	targetPath := filepath.Join(destPath, fileName)

	if _, err := os.Stat(targetPath); err == nil {
		ext := filepath.Ext(fileName)

		newFileName := fmt.Sprintf("%s_%s%s", strings.TrimSuffix(fileName, ext), time.Now().Format("20060102_150405"), ext)

		targetPath = filepath.Join(destPath, newFileName)
		log.Printf("İsim çakışması önlendi. Yeni isim: %s\n", newFileName)
	}

	err := os.Rename(sourcePath, targetPath)
	if err != nil {
		log.Printf("error cannot move (%s): %v\n", sourcePath, err)
		return
	}
	log.Printf("Succeess: %s -> %s moved folder.\n", fileName, destPath)
}
