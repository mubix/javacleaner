// types/types.go
package types

type JavaInstall struct {
	Version string
	Path    string
	Source  string // Registry, Path, JAVA_HOME, etc.
}

type UsageInfo struct {
	LastUsed string   // e.g., timestamp or 'unknown'
	UsedBy   []string // process names, services, etc.
}
