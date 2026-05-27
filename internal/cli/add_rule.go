package cli

//
//import (
//	"encoding/json"
//	"fmt"
//	"log"
//	"os"
//	"strings"
//
//	"github.com/barisaydogdu/magic-byte-organizer/internal/config"
//	"github.com/spf13/cobra"
//)
//
//var addRuleCmd = &cobra.Command{
//	Use:   "add-rule [file_type] [target_folder]",
//	Short: "Adds a new file transfer rule to the system.",
//
//	Args: cobra.ExactArgs(2),
//
//	Run: func(cmd *cobra.Command, args []string) {
//		fileType := args[0]
//		targetFolder := args[1]
//
//		log.Printf("adding rule for %s to %s", fileType, targetFolder)
//
//		configs := config.LoadConfig("config.json")
//
//		val, exists := configs[fileType]
//		if exists {
//			fmt.Printf("Warning: There is already a rule for '%s'. Should it be overwritten? (yes/no):", val)
//		}
//
//		var answer string
//		_, err2 := fmt.Scanln(&answer)
//		if err2 != nil {
//			return
//		}
//
//		if strings.ToLower(answer) != "yes" && strings.ToLower(answer) != "y" {
//			fmt.Println("Aborting...")
//			return
//		}
//
//		configs[fileType] = targetFolder
//
//		newJson, err := json.MarshalIndent(configs, "", "  ")
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		err = os.WriteFile("config.json", newJson, 0644)
//		if err != nil {
//			return
//		}
//
//	},
//}
