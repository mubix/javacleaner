// cleanup/remover.go
package cleanup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"javacleaner/config"
	"javacleaner/types"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func findUninstallCommand(javaPath string) (string, bool) {
	subkeys := []string{
		`SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall`,
		`SOFTWARE\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall`,
	}

	roots := []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER}
	for _, root := range roots {
		for _, base := range subkeys {
			key, err := registry.OpenKey(root, base, registry.READ)
			if err != nil {
				continue
			}
			defer key.Close()

			names, err := key.ReadSubKeyNames(-1)
			if err != nil {
				continue
			}

			for _, name := range names {
				subkey, err := registry.OpenKey(key, name, registry.READ)
				if err != nil {
					continue
				}
				path, _, _ := subkey.GetStringValue("InstallLocation")
				if filepath.Clean(strings.ToLower(javaPath)) == filepath.Clean(strings.ToLower(path)) {
					cmd, _, err := subkey.GetStringValue("UninstallString")
					subkey.Close()
					if err == nil && cmd != "" {
						return cmd, true
					}
				}
				subkey.Close()
			}
		}
	}
	return "", false
}

func isAdmin() bool {
	const TokenElevationTypeDefault = 1
	const TokenElevationTypeFull = 2

	hToken := windows.Token(0)
	type elevationType uint32
	var et elevationType
	var returnedLen uint32

	err := windows.GetTokenInformation(
		hToken,
		windows.TokenElevationType,
		(*byte)(unsafe.Pointer(&et)),
		uint32(unsafe.Sizeof(et)),
		&returnedLen,
	)
	if err != nil {
		return false
	}
	return et == TokenElevationTypeFull
}

func SafeRemove(installs []types.JavaInstall, usageMap map[string]types.UsageInfo, cfg *config.Config) {
	log.Println("Beginning safe removal process")
	admin := isAdmin()

	for _, inst := range installs {
		info := usageMap[inst.Path]
		if len(info.UsedBy) > 0 {
			continue
		}

		fmt.Printf("[?] Remove %s (version: %s)? [y/N]: ", inst.Path, inst.Version)
		var input string
		fmt.Scanln(&input)
		if strings.ToLower(strings.TrimSpace(input)) == "y" {
			if uninstallCmd, found := findUninstallCommand(inst.Path); found {
				fmt.Printf("[~] Running uninstaller for %s...\n", inst.Path)
				fmt.Printf("[~] Command: %s\n", uninstallCmd)
				log.Printf("Executing uninstall command: %s\n", uninstallCmd)
				cmd := exec.Command("cmd", "/C", uninstallCmd)
				err := cmd.Run()
				if err != nil {
					fmt.Printf("[-] Uninstaller failed: %v\n", err)
					log.Printf("Uninstaller failed for %s: %v\n", inst.Path, err)
					if !admin {
						fmt.Println("[!] Try re-running Java Cleaner as Administrator for better results.")
						log.Println("Recommendation: rerun as Administrator")
					}
				} else {
					fmt.Printf("[+] Uninstaller completed for %s\n", inst.Path)
					log.Printf("Uninstalled via uninstaller: %s\n", inst.Path)
				}
			} else {
				err := os.RemoveAll(filepath.Clean(inst.Path))
				if err != nil {
					fmt.Printf("[-] Failed to remove %s: %v\n", inst.Path, err)
					log.Printf("Failed to remove %s: %v\n", inst.Path, err)
					if !admin {
						fmt.Println("[!] Try re-running Java Cleaner as Administrator for better results.")
						log.Println("Recommendation: rerun as Administrator")
					}
				} else {
					fmt.Printf("[+] Removed %s\n", inst.Path)
					log.Printf("Removed %s\n", inst.Path)
				}
			}
		} else {
			fmt.Printf("[-] Skipped %s\n", inst.Path)
			log.Printf("Skipped removal of %s\n", inst.Path)
		}
	}
}
