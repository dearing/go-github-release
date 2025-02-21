package main

import (
	"flag"
	"log/slog"
	"os"
	"runtime/debug"
	"time"
)

var argVersion = flag.Bool("version", false, "print version and exit")

func main() {
	status := run()
	os.Exit(status)
}

func run() int {
	start := time.Now()
	defer func() {
		slog.Info("operation complete", "duration", time.Since(start))
	}()

	flag.Usage = usage
	flag.Parse()

	if *argVersion {
		version()
		return NoError
	}

	return NoError
}

func usage() {
	slog.Info("Usage: [go tool] go-github-release [options] <file>...")
	flag.PrintDefaults()
}

func version() {
	// seems like a nice place to sneak in some debug information
	info, ok := debug.ReadBuildInfo()
	if ok {
		slog.Info("build info", "main", info.Main.Path, "version", info.Main.Version)
		for _, setting := range info.Settings {
			slog.Info("build info", "key", setting.Key, "value", setting.Value)
		}
	}
}
