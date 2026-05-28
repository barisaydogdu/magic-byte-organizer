package cli

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/barisaydogdu/magic-byte-organizer/internal/config"
	"github.com/barisaydogdu/magic-byte-organizer/internal/detector"
	"github.com/barisaydogdu/magic-byte-organizer/internal/mover"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var watchDir string
var isDryRun bool
var delayStr string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "It starts the file monitoring service.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

		defer cancel()
		parsedDelay, err := time.ParseDuration(delayStr)
		if err != nil {
			log.Println("cannot parsed delaystr", err)
			return
		}

		log.Println("MagicSort tracking service is launching...")

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Println("cannot create watcher", err)
			return
		}
		defer watcher.Close()

		absConfigPath := config.GetAbsoluteConfigPath()
		configRules, err := config.LoadConfig(absConfigPath)
		if err != nil {
			log.Println("cannot load config", err)
			return
		}

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Has(fsnotify.Create) || event.Has(fsnotify.Rename) {
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
										if parsedDelay > 0 {
											select {
											case <-time.After(parsedDelay):
												log.Println("parsed delay done")
											case <-ctx.Done():
												return
											}

											_, err := os.Stat(path)
											if err != nil {
												if os.IsNotExist(err) {
													return
												}
											}
										}

										err = mover.MoveFile(path, realTargetDir)
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
			log.Println("cannot add watcher", err)
			return
		}

		<-ctx.Done()

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&watchDir, "dir", "d", "../home/baris/downloads", "The full path to the folder (source) to be monitored.")
	startCmd.Flags().BoolVarP(&isDryRun, "dry-run", "", false, "Runs a simulation of a file transfer scenario without making any changes to the disk.")
	startCmd.Flags().StringVarP(&delayStr, "delay", "t", "0s", "Wait time before moving the file (e.g., '10m', '30s', '1h'). Default is 0s.")
}
