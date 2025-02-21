package main

import (
	"flag"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

var argVersion = flag.Bool("version", false, "print version and exit")
var argDir = flag.String("dir", "build", "directory to search for files")

func init() {
	// feels silly to have to do this
	mime.AddExtensionType(".exe", "application/octet-stream")
	mime.AddExtensionType(".zip", "application/zip")
	mime.AddExtensionType(".tar.gz", "application/gzip")
	mime.AddExtensionType(".txt", "text/plain")
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if *argVersion {
		version()
		return
	}

	status := run()
	os.Exit(status)
}

func run() int {
	// clock the operation
	start := time.Now()
	defer func() {
		slog.Info("operation complete", "duration", time.Since(start))
	}()

	// get GITHUB_TOKEN from environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		slog.Error("missing GITHUB_TOKEN")
		return TokenNotFound
	}

	//errChan := make(chan error)

	r := CreateReleaseRequest{
		TagName:              "v0.0.1",
		TargetCommitish:      "main",
		Name:                 "v0.0.1",
		Body:                 "initial release",
		Draft:                true,
		Prerelease:           false,
		GenerateReleaseNotes: false,
	}

	resp, err := r.CreateRelease("dearing", "go-github-release", token)
	if err != nil {
		slog.Error("release issue", "err", err)
		return UnknownError
	}
	slog.Info("release created", "id", resp.ID, "url", resp.HTMLURL)

	if *argDir != "" {
		files, err := filepath.Glob(filepath.Join(*argDir, "*"))
		if err != nil {
			slog.Error("file glob issue", "err", err)
			return UnknownError
		}
		slog.Info("files", "count", len(files))

		for _, file := range files {
			UploadAsset(1, file)
		}

		slog.Info("files uploaded", "count", len(files))
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

// func prettyPrint(v any) {
// 	b, _ := json.MarshalIndent(v, "", "  ")
// 	println(string(b))
// }
