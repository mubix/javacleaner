package config

import (
	"flag"
	"javacleaner/types"
)

type Config struct {
	Interactive bool
	Verbose     bool
	AutoClean   bool
	OutputJSON  bool

	Installs []types.JavaInstall
	UsageMap map[string]types.UsageInfo
}

func LoadConfig() *Config {
	interactive := flag.Bool("interactive", false, "Run in interactive mode")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	autoClean := flag.Bool("auto-clean", false, "Automatically clean unused Java versions")
	outputJSON := flag.Bool("output", false, "Output result in JSON")

	flag.Parse()

	return &Config{
		Interactive: *interactive,
		Verbose:     *verbose,
		AutoClean:   *autoClean,
		OutputJSON:  *outputJSON,
		Installs:    []types.JavaInstall{},
		UsageMap:    map[string]types.UsageInfo{},
	}
}
