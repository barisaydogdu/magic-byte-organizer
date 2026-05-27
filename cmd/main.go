package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	configRules := loadConfig("config.json")
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Create) {
					//fmt.Println("created:", event.Name)
					//go detectFileType(event.Name)
					go func(path string) {
						if waitForFileLock(path) {
							fileType, err := detectFileType(path)
							if err != nil {
								log.Println(err)
							}

							if fileType == "" {
								return
							}

							rawTargerDir, exist := configRules[fileType]

							if exist {

								realTargetDir := expandHomePath(rawTargerDir)

								moveFile(path, realTargetDir)
							} else {
								log.Printf("Kural bulunamadı, dosya atlandı: %s (Tür: %s)\n", path, fileType)
							}

						} else {
							log.Println("file lock timeout")
						}
					}(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add("/tmp")
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}
