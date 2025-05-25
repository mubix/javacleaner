// javacleaner.go
package main

import (
	"fmt"
	"os"
	"strings"

	"javacleaner/cleanup"
	"javacleaner/config"
	"javacleaner/detector"
	"javacleaner/output"
	"javacleaner/usage"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.Interactive {
		runInteractive(cfg)
	} else {
		runHeadless(cfg)
	}
}

func runHeadless(cfg *config.Config) {
	fmt.Println("[+] Scanning for installed Java versions...")
	installs := detector.DetectJavaVersions(cfg)
	cfg.Installs = installs
	fmt.Printf("[+] Found %d Java installations\n", len(installs))

	fmt.Println("[+] Checking Java usage...")
	usageMap := usage.CheckJavaUsage(installs, cfg)
	cfg.UsageMap = usageMap

	output.PrintReport(installs, usageMap, cfg)

	if cfg.AutoClean {
		fmt.Println("[!] Running auto-clean...")
		cleanup.SafeRemove(installs, usageMap, cfg)
	}
}

func runInteractive(cfg *config.Config) {
	fmt.Println("=== Java Cleaner (Interactive Mode) ===")
	for {
		fmt.Println("\n1. Scan for Java Versions")
		fmt.Println("2. Check Usage")
		fmt.Println("3. Show Safe-to-Remove")
		fmt.Println("4. Generate Report")
		fmt.Println("5. Remove Versions")
		fmt.Println("6. Exit")
		fmt.Print("Select an option: ")

		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			cfg.Installs = detector.DetectJavaVersions(cfg)
			fmt.Printf("Found %d installations.\n", len(cfg.Installs))
		case "2":
			cfg.UsageMap = usage.CheckJavaUsage(cfg.Installs, cfg)
			fmt.Println("Usage checked.")
		case "3":
			output.PrintRemovable(cfg.Installs, cfg.UsageMap, cfg)
		case "4":
			output.PrintReport(cfg.Installs, cfg.UsageMap, cfg)
		case "5":
			cleanup.SafeRemove(cfg.Installs, cfg.UsageMap, cfg)
		case "6":
			os.Exit(0)
		default:
			fmt.Println("Invalid option")
		}
	}
}
