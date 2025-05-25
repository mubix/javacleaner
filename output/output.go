// output/output.go
package output

import (
	"encoding/json"
	"fmt"
	"javacleaner/config"
	"javacleaner/types"
	"log"
)

func PrintReport(installs []types.JavaInstall, usageMap map[string]types.UsageInfo, cfg *config.Config) {
	if cfg.OutputJSON {
		report := make(map[string]interface{})
		for _, inst := range installs {
			report[inst.Path] = map[string]interface{}{
				"Version":  inst.Version,
				"Source":   inst.Source,
				"UsedBy":   usageMap[inst.Path].UsedBy,
				"LastUsed": usageMap[inst.Path].LastUsed,
			}
		}
		jsonBytes, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			log.Println("Error generating JSON output:", err)
			return
		}
		fmt.Println(string(jsonBytes))
		log.Println("Generated JSON report")
		return
	}

	log.Println("Printing human-readable Java report")
	fmt.Println("\nJava Installation Report")
	fmt.Println("========================")
	for _, inst := range installs {
		info := usageMap[inst.Path]
		fmt.Printf("Path: %s\n", inst.Path)
		fmt.Printf("  Version : %s\n", inst.Version)
		fmt.Printf("  Source  : %s\n", inst.Source)
		fmt.Printf("  Used By :\n")
		if len(info.UsedBy) == 0 {
			fmt.Println("    (No known usage)")
		} else {
			for _, used := range info.UsedBy {
				fmt.Printf("    %s\n", used)
			}
		}
		fmt.Printf("  Last Used: %s\n\n", info.LastUsed)
	}
}

func PrintRemovable(installs []types.JavaInstall, usageMap map[string]types.UsageInfo, cfg *config.Config) {
	log.Println("Listing safe-to-remove Java versions")
	fmt.Println("\nSafe-to-Remove Java Versions")
	fmt.Println("============================")
	for _, inst := range installs {
		info := usageMap[inst.Path]
		if len(info.UsedBy) == 0 {
			fmt.Printf("%s (%s)\n", inst.Path, inst.Version)
			log.Printf("Marked for removal: %s (%s)\n", inst.Path, inst.Version)
		}
	}
}
