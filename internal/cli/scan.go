package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/barisaydogdu/magic-byte-organizer/internal/config"
	"github.com/barisaydogdu/magic-byte-organizer/internal/detector"
	"github.com/barisaydogdu/magic-byte-organizer/internal/mover"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a file to move to a target directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Scanning is starting: %s\n", watchDir)

		configRules := config.LoadConfig("config.json")

		entries, err := os.ReadDir(watchDir)
		if err != nil {
			return
		}

		for i := 0; i < len(entries); i++ {

			if entries[i].IsDir() {
				continue
			}

			if strings.HasSuffix(entries[i].Name(), ".md") || strings.HasSuffix(entries[i].Name(), ".crdownload") ||
				strings.HasSuffix(entries[i].Name(), ".part") || strings.HasSuffix(entries[i].Name(), ".part") || strings.HasPrefix(entries[i].Name(), ".") {
				continue
			}
			path := filepath.Join(watchDir, entries[i].Name())

			fileType, err := detector.DetectFileType(path)
			if err != nil {
				fmt.Println(err)
				return
			}

			isMatch, err := detector.CheckExtAndType(path, fileType)
			if err != nil {
				fmt.Println(err)
				return
			}

			if !isMatch {
				//todo: burada raporlayabiliriz
				//	fmt.Println("File type not match")
				//	return
			}

			rawTargetDir, exist := configRules[fileType]

			if exist {
				realTargetDir := config.ExpandHomePath(rawTargetDir)

				if isDryRun {
					log.Printf("\033[33m[DRY RUN]\033[0m Simülasyon: '%s' dosyası '%s' klasörüne taşınacaktı.\n", path, realTargetDir)
				} else {
					err := mover.MoveFile(path, realTargetDir)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			} else {
				log.Printf("Kural bulunamadı, dosya atlandı: %s (Tür: %s)\n", path, fileType)
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&watchDir, "dir", "d", "/Downloads", "The full path to the source directory to scan and organize.")

	scanCmd.Flags().BoolVarP(&isDryRun, "dry-run", "", false, "Fiziksel taşıma yapmadan simülasyon loglarını gösterir.")
}
