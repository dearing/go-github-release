package main

const (
	NoError = iota
	UnknownError
	OwnerNotFound
	RepoNotFound
	TokenNotFound
	ErrorCreateRequest
	ErrorBadPattern
	AssetReadError
	AssetUploadError
)
