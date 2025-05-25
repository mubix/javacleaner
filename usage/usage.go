// usage/usage.go
package usage

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"javacleaner/config"
	"javacleaner/types"

	"github.com/shirou/gopsutil/v3/process"
)

func CheckJavaUsage(installs []types.JavaInstall, cfg *config.Config) map[string]types.UsageInfo {
	usageMap := make(map[string]types.UsageInfo)

	procs, err := process.Processes()
	if err != nil {
		log.Println("Error retrieving process list:", err)
		return usageMap
	}

	svcCmd := exec.Command("powershell", "-Command", "Get-Service | Where-Object {$_.PathName -match 'java'} | Select-Object -ExpandProperty Name")
	svcOut, err := svcCmd.Output()
	var services []string
	if err == nil {
		services = strings.Split(strings.TrimSpace(string(svcOut)), "\n")
	}

	taskCmd := exec.Command("schtasks.exe", "/query", "/fo", "LIST")
	taskOut, err := taskCmd.Output()
	var taskLines []string
	if err == nil {
		taskLines = strings.Split(string(taskOut), "\n")
	}

	for _, inst := range installs {
		info := types.UsageInfo{
			LastUsed: "unknown",
			UsedBy:   []string{},
		}

		instPathLower := strings.ToLower(inst.Path)

		for _, p := range procs {
			exePath, err := p.Exe()
			if err != nil {
				continue
			}
			exePathLower := strings.ToLower(exePath)

			if strings.HasPrefix(exePathLower, instPathLower) || strings.HasSuffix(exePathLower, "\\java.exe") {
				name, _ := p.Name()
				info.UsedBy = append(info.UsedBy, fmt.Sprintf("PID %d: %s", p.Pid, name))
			}
		}

		for _, svc := range services {
			if strings.TrimSpace(svc) != "" {
				info.UsedBy = append(info.UsedBy, fmt.Sprintf("Service: %s", strings.TrimSpace(svc)))
			}
		}

		for _, line := range taskLines {
			if strings.Contains(strings.ToLower(line), "java") && !strings.Contains(strings.ToLower(line), "javaupdatesched") {
				info.UsedBy = append(info.UsedBy, fmt.Sprintf("Scheduled Task Ref: %s", strings.TrimSpace(line)))
			}
		}

		usageMap[inst.Path] = info
	}

	return usageMap
}
