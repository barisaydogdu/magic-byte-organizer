package cli

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/barisaydogdu/magic-byte-organizer/internal/config"
	"github.com/barisaydogdu/magic-byte-organizer/internal/detector"
	"github.com/barisaydogdu/magic-byte-organizer/internal/mover"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var watchDir string
var isDryRun bool

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "It starts the file monitoring service.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MagicSort tracking service is launching...")
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		configRules := config.LoadConfig("config.json")
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Has(fsnotify.Create) {
						go func(path string) {
							fileName := filepath.Base(path)

							if strings.HasSuffix(fileName, ".md") || strings.HasSuffix(fileName, ".crdownload") ||
								strings.HasSuffix(fileName, ".part") || strings.HasSuffix(fileName, ".part") || strings.HasPrefix(fileName, ".") {

								return
							}

							if detector.WaitForFileLock(path) {
								fileType, err := detector.DetectFileType(path)
								if err != nil {
									log.Println(err)
								}

								if fileType == "" {
									return
								}

								isMatch, err := detector.CheckExtAndType(path, fileType)
								if err != nil {
									log.Println(err)
								}

								if !isMatch {
									//todo: burada raporlayabiliriz
									//	fmt.Println("File type not match")
									//	return
								}

								rawTargerDir, exist := configRules[fileType]

								if exist {
									realTargetDir := config.ExpandHomePath(rawTargerDir)

									if isDryRun {
										log.Printf("\033[33m[DRY RUN]\033[0m Simülasyon: '%s' dosyası '%s' klasörüne taşınacaktı.\n", path, realTargetDir)
									} else {
										err := mover.MoveFile(path, realTargetDir)
										if err != nil {
											return
										}
									}

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
		err = watcher.Add(watchDir)
		if err != nil {
			log.Fatal(err)
		}

		<-make(chan struct{})
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	//./magicsort start -d /tmp
	startCmd.Flags().StringVarP(&watchDir, "dir", "d", "../home/baris/downloads", "The full path to the folder (source) to be monitored.")
	startCmd.Flags().BoolVarP(&isDryRun, "dry-run", "", false, "")

}
