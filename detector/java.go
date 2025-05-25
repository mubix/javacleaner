// detector/java.go
package detector

import (
	"os"
	"path/filepath"
	"strings"

	"javacleaner/types"

	"golang.org/x/sys/windows/registry"
)

func DetectJavaVersions(cfg interface{}) []types.JavaInstall {
	var installs []types.JavaInstall

	// Search registry (64-bit)
	regPaths := []registry.Key{
		registry.LOCAL_MACHINE,
		registry.CURRENT_USER,
	}

	for _, baseKey := range regPaths {
		k, err := registry.OpenKey(baseKey, `SOFTWARE\JavaSoft\Java Runtime Environment`, registry.READ)
		if err == nil {
			versions, _ := k.ReadSubKeyNames(-1)
			for _, version := range versions {
				path, _, _ := k.GetStringValue("JavaHome")
				if path != "" {
					installs = append(installs, types.JavaInstall{
						Version: version,
						Path:    path,
						Source:  "Registry",
					})
				}
			}
			k.Close()
		}
	}

	// Check JAVA_HOME env var
	if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		installs = append(installs, types.JavaInstall{
			Version: "unknown",
			Path:    javaHome,
			Source:  "JAVA_HOME",
		})
	}

	// Check Program Files
	roots := []string{
		os.Getenv("ProgramFiles"),
		os.Getenv("ProgramFiles(x86)"),
	}

	for _, root := range roots {
		if root == "" {
			continue
		}
		javaRoot := filepath.Join(root, "Java")
		entries, err := os.ReadDir(javaRoot)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() && strings.HasPrefix(strings.ToLower(entry.Name()), "jre") || strings.HasPrefix(strings.ToLower(entry.Name()), "jdk") {
					installs = append(installs, types.JavaInstall{
						Version: entry.Name(),
						Path:    filepath.Join(javaRoot, entry.Name()),
						Source:  "ProgramFiles",
					})
				}
			}
		}
	}

	return installs
}
