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

func usage() {
	slog.Info("Usage: [go tool] go-github-release [options] <file>...")
	flag.PrintDefaults()
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

	// clock the whole enchalada
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

	// create a release request populated with some defaults
	r := CreateReleaseRequest{
		TagName:              "v0.0.1",
		TargetCommitish:      "main",
		Name:                 "v0.0.1",
		Body:                 "initial release",
		Draft:                true,
		Prerelease:           false,
		GenerateReleaseNotes: false,
	}

	// create a release on GitHub and get the response struct back
	resp, err := r.CreateRelease("dearing", "go-github-release", token)
	if err != nil {
		slog.Error("release issue", "err", err)
		return ErrorCreateRequest
	}
	slog.Info("release created", "id", resp.ID, "url", resp.HTMLURL)

	// get a list of assets in the directory
	assets, err := filepath.Glob(filepath.Join(*argDir, "*"))
	if err != nil {
		slog.Error("glob issue", "err", err)
		return ErrorBadPattern
	}
	slog.Info("assets", "count", len(assets))

	// upload each file to the release
	for _, file := range assets {
		uploadStart := time.Now()
		err := uploadAsset(token, resp.UploadURL, file)
		if err != nil {
			slog.Error("upload issue", "err", err)
			return AssetUploadError
		}
		slog.Info("asset uploaded", "file", file, "duration", time.Since(uploadStart))
	}
	slog.Info("assets uploaded", "count", len(assets))

	return NoError
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
