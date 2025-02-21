package main

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

var argVersion = flag.Bool("version", false, "print version and exit")
var assetDir = flag.String("asset-dir", "build", "directory to search for assets to upload")

var tagName = flag.String("tag-name", "", "release tag to reference or create")
var targetCommitish = flag.String("target-commitish", "", "when creating a tag, the commit to reference")

var githubOwner = flag.String("github-owner", "", "github owner")
var githubRepo = flag.String("github-repo", "", "github repo")

var releaseName = flag.String("release-name", "", "release name, defaults to tag name")
var releaseBody = flag.String("release-body", "no results", "release body contents (markdown supported)")
var releaseDraft = flag.Bool("release-draft", true, "set draft flag on release")
var releasePrerelease = flag.Bool("release-prerelease", false, "set prerelease flag on release")
var generateReleaseNotes = flag.Bool("generate-release-notes", false, "tell GitHub to generate release notes")

func usage() {
	println(`Usage: [go tool] go-github-release [options]
Simple GitHub tool to cut a release and upload assets.

NOTE: An env:GITHUB_TOKEN is required to authenticate with GitHub and its 
recommended to use a personal access token scoped to repository.

GitHub releases are fairly simple, its a ccollection of assets with name and 
message body that reference a commit, which can be a newly created tag as a 
commitish (--target-commitish) or existing tag in which case the target 
commitish is ignored. Assets are uploaded to this release based on mime-type 
which is inferred from the file extension.

If you are unsure of what workflow to use, the author recommends creating a 
tag ahead of release and using that tag as the --tag-name. This way you can
seperate the release process from the tag creation process and ensure your 
tag is referencing the proper moment in history. The automatic tag creation 
feature is a GitHub conveinence and not related to git.

ex: to draft a new release for an existing tag v1:

[go tool] go-github-release \
	--github-owner myorg \
	--github-repo myrepo \
	--tag-name v1 \
	--asset-dir build

ex: to draft a new release and new tag v1 pointing to abc123:

[go tool] go-github-release \
	--github-owner myorg \
	--github-repo myrepo \
	--tag-name v1 \
	--target-commitish abc123 \
	--asset-dir dist

Tips:
  - drafts that create a tag, won't do so until the release is published
  - if --github-repo is not set, the current directory is used
  - if --release-name is not set, the tag name is used
  - if --target-commitish is not set, main is used
  - a release body supports markdown
  - you can safely edit the release body while the assets are uploading
  - find the latest tag by name: git describe --tags --abbrev=0
  - find a tag's reference commit: git rev-list -n 1 tags/v0.0.1

Options:
`)
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
			slog.Warn("--github-repo not set, using working dir", "repo", *githubRepo)
		} else {
			slog.Error("--github-repo is required")
			return RepoNotFound
		}
	}

	if *tagName == "" {
		slog.Error("--tag-name is required")
		return TagNameRequired
	}

	// ensure the asset directory exists
	if _, err := os.Stat(*assetDir); os.IsNotExist(err) {
		slog.Error("asset directory not found", "dir", *assetDir)
		return AssetDirNotFound
	}

	// set release name to tag name if not set
	if *releaseName == "" {
		*releaseName = *tagName
		slog.Warn("--release-name not set, using tag name", "name", *releaseName)
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
