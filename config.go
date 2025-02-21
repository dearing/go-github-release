package main

type Config struct {
	Owner       string
	Repos       string
	Token       string
	ReleaseName string
	ReleaseTag  string
	ReleaseBody string

	Draft                bool
	Prerelease           bool
	GenerateReleaseNotes bool
	MakeLatest           bool

	//DiscussionCategoryName string

	// Assets is a list of files to be uploaded to the release.
	Assets []string
}

// Status codes
const (
	ReleaseCreated  = 201
	ReleaseNotFound = 404
	ReleaseNotValid = 422
)
