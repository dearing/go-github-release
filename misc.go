package main

import (
	"log/slog"
	"mime"
	"runtime/debug"
)

const (
	NoError = iota
	UnknownError
	OwnerNotFound
	RepoNotFound
	TagNameRequired
	TokenNotFound
	AssetDirNotFound
	ErrorCreateRequest
	ErrorBadPattern
	AssetReadError
	AssetUploadError
)

func init() {
	// feels silly to have to do this
	mime.AddExtensionType(".exe", "application/octet-stream")
	mime.AddExtensionType(".zip", "application/zip")
	mime.AddExtensionType(".tar.gz", "application/gzip")
	mime.AddExtensionType(".txt", "text/plain")
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
