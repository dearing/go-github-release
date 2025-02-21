package main

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

var argVersion = flag.Bool("version", false, "print version and exit")
var assetDir = flag.String("asset-dir", "build", "directory to search for assets to upload")

var tagName = flag.String("tag-name", "v0.0.1", "release tag name")

var githubOwner = flag.String("github-owner", "", "github owner")
var githubRepo = flag.String("github-repo", "", "github repo")

var targetCommitish = flag.String("target-commitish", "main", "target commitish tie the release to")
var releaseName = flag.String("release-name", "v0.0.1", "release name")
var releaseBody = flag.String("release-body", "initial release", "release body contents (markdown supported)")
var releaseDraft = flag.Bool("release-draft", true, "draft")
var releasePrerelease = flag.Bool("release-prerelease", false, "prerelease")
var generateReleaseNotes = flag.Bool("generate-release-notes", false, "generate release notes")

func usage() {
	slog.Info("Usage: [go tool] go-github-release [options]")
	flag.PrintDefaults()
}

func main() {

	flag.Usage = usage
	flag.Parse()

	// print some debug information and exit
	if *argVersion {
		version()
		return
	}

	// we wrap the run inside the main function so we can capture a status
	// code and cleanly let defer functions within run, uh run
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

	if *githubOwner == "" {
		slog.Error("--github-owner is required")
		return OwnerNotFound
	}

	if *githubRepo == "" {

		pwd, err := os.Getwd()
		if err == nil {
			*githubRepo = filepath.Base(pwd)
			slog.Info("repo inferred", "repo", *githubRepo)
		}
	}

	// create a release request populated with some defaults
	r := CreateReleaseRequest{
		TagName:              *tagName,
		TargetCommitish:      *targetCommitish,
		Name:                 *releaseName,
		Body:                 *releaseBody,
		Draft:                *releaseDraft,
		Prerelease:           *releasePrerelease,
		GenerateReleaseNotes: *generateReleaseNotes,
	}

	// create a release on GitHub and get the response struct back
	resp, err := r.CreateRelease(token, *githubOwner, *githubRepo)
	if err != nil {
		slog.Error("release issue", "err", err)
		return ErrorCreateRequest
	}
	slog.Info("release created", "id", resp.ID, "url", resp.HTMLURL)

	// get a list of assets in the directory
	assets, err := filepath.Glob(filepath.Join(*assetDir, "*"))
	if err != nil {
		slog.Error("glob issue", "err", err)
		return ErrorBadPattern
	}
	slog.Info("assets", "count", len(assets))

	// upload each file to the release
	for _, asset := range assets {
		uploadStart := time.Now()
		err := uploadAsset(token, resp.UploadURL, asset)
		if err != nil {
			slog.Error("upload issue", "err", err)
			return AssetUploadError
		}
		slog.Info("asset uploaded", "asset", asset, "duration", time.Since(uploadStart))
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
