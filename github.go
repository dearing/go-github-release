package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Status codes
const (
	ReleaseCreated  = 201
	ReleaseNotFound = 404
	ReleaseNotValid = 422
)

// CreateReleaseRequest is the request to create a release on GitHub
type CreateReleaseRequest struct {
	TagName              string `json:"tag_name,omitempty"`
	TargetCommitish      string `json:"target_commitish,omitempty"`
	Name                 string `json:"name,omitempty"`
	Body                 string `json:"body,omitempty"`
	Draft                bool   `json:"draft,omitempty"`
	Prerelease           bool   `json:"prerelease,omitempty"`
	GenerateReleaseNotes bool   `json:"generate_release_notes,omitempty"`
}

// CreateReleaseResponse is the response from creating a release on GitHub
type CreateReleaseResponse struct {
	URL             string    `json:"url,omitempty"`
	HTMLURL         string    `json:"html_url,omitempty"`
	AssetsURL       string    `json:"assets_url,omitempty"`
	UploadURL       string    `json:"upload_url,omitempty"`
	TarballURL      string    `json:"tarball_url,omitempty"`
	ZipballURL      string    `json:"zipball_url,omitempty"`
	DiscussionURL   string    `json:"discussion_url,omitempty"`
	ID              int       `json:"id,omitempty"`
	NodeID          string    `json:"node_id,omitempty"`
	TagName         string    `json:"tag_name,omitempty"`
	TargetCommitish string    `json:"target_commitish,omitempty"`
	Name            string    `json:"name,omitempty"`
	Body            string    `json:"body,omitempty"`
	Draft           bool      `json:"draft,omitempty"`
	Prerelease      bool      `json:"prerelease,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Author          Author    `json:"author"`
	Assets          []Assets  `json:"assets,omitempty"`
}

// UploadAssetResponse is the response from uploading an asset to a release on GitHub
type UploadAssetResponse struct {
	URL                string    `json:"url,omitempty"`
	BrowserDownloadURL string    `json:"browser_download_url,omitempty"`
	ID                 int       `json:"id,omitempty"`
	NodeID             string    `json:"node_id,omitempty"`
	Name               string    `json:"name,omitempty"`
	Label              string    `json:"label,omitempty"`
	State              string    `json:"state,omitempty"`
	ContentType        string    `json:"content_type,omitempty"`
	Size               int       `json:"size,omitempty"`
	DownloadCount      int       `json:"download_count,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Uploader           Uploader  `json:"uploader"`
}

// Author is the author of the release
type Author struct {
	Login             string `json:"login,omitempty"`
	ID                int    `json:"id,omitempty"`
	NodeID            string `json:"node_id,omitempty"`
	AvatarURL         string `json:"avatar_url,omitempty"`
	GravatarID        string `json:"gravatar_id,omitempty"`
	URL               string `json:"url,omitempty"`
	HTMLURL           string `json:"html_url,omitempty"`
	FollowersURL      string `json:"followers_url,omitempty"`
	FollowingURL      string `json:"following_url,omitempty"`
	GistsURL          string `json:"gists_url,omitempty"`
	StarredURL        string `json:"starred_url,omitempty"`
	SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
	OrganizationsURL  string `json:"organizations_url,omitempty"`
	ReposURL          string `json:"repos_url,omitempty"`
	EventsURL         string `json:"events_url,omitempty"`
	ReceivedEventsURL string `json:"received_events_url,omitempty"`
	Type              string `json:"type,omitempty"`
	SiteAdmin         bool   `json:"site_admin,omitempty"`
}

// Uploader is the uploader of the asset
type Uploader struct {
	Login             string `json:"login,omitempty"`
	ID                int    `json:"id,omitempty"`
	NodeID            string `json:"node_id,omitempty"`
	AvatarURL         string `json:"avatar_url,omitempty"`
	GravatarID        string `json:"gravatar_id,omitempty"`
	URL               string `json:"url,omitempty"`
	HTMLURL           string `json:"html_url,omitempty"`
	FollowersURL      string `json:"followers_url,omitempty"`
	FollowingURL      string `json:"following_url,omitempty"`
	GistsURL          string `json:"gists_url,omitempty"`
	StarredURL        string `json:"starred_url,omitempty"`
	SubscriptionsURL  string `json:"subscriptions_url,omitempty"`
	OrganizationsURL  string `json:"organizations_url,omitempty"`
	ReposURL          string `json:"repos_url,omitempty"`
	EventsURL         string `json:"events_url,omitempty"`
	ReceivedEventsURL string `json:"received_events_url,omitempty"`
	Type              string `json:"type,omitempty"`
	SiteAdmin         bool   `json:"site_admin,omitempty"`
}

// Assets are the assets of the release
type Assets struct {
	URL                string    `json:"url,omitempty"`
	BrowserDownloadURL string    `json:"browser_download_url,omitempty"`
	ID                 int       `json:"id,omitempty"`
	NodeID             string    `json:"node_id,omitempty"`
	Name               string    `json:"name,omitempty"`
	Label              string    `json:"label,omitempty"`
	State              string    `json:"state,omitempty"`
	ContentType        string    `json:"content_type,omitempty"`
	Size               int       `json:"size,omitempty"`
	DownloadCount      int       `json:"download_count,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Uploader           Uploader  `json:"uploader"`
}

// CreateRelease creates a release on GitHub
//
//	ex: request.CreateRelease("dearing", "go-github-release", token)
func (r *CreateReleaseRequest) CreateRelease(owner, repo, token string) (CreateReleaseResponse, error) {

	// create the URL for the request
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)

	// marshal this request into JSON for the body
	data, err := json.Marshal(r)
	if err != nil {
		return CreateReleaseResponse{}, fmt.Errorf("create release issue: %w", err)
	}

	ctx := context.Background()

	// create the request passing a bytes buffer of our payload
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return CreateReleaseResponse{}, fmt.Errorf("create request issue: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Accept", "application/vnd.github+json") // github vendor specific media type; recommended
	req.Header.Set("Content-Type", "application/json")      // we'll marshal our data into JSON

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CreateReleaseResponse{}, fmt.Errorf("create release issue: %w", err)
	}
	defer resp.Body.Close()

	// 422 is unprocessable entity and would likely mean the user set a bad commit-ish
	if resp.StatusCode != ReleaseCreated {
		return CreateReleaseResponse{}, fmt.Errorf("create release unexpected status: %s", resp.Status)
	}

	slog.Debug("create release", "status", resp.Status)
	slog.Debug("create release", "headers", resp.Header)

	info := CreateReleaseResponse{}

	// decode the response into our struct
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return CreateReleaseResponse{}, fmt.Errorf("create release decode issue: %w", err)
	}

	return info, nil
}

// UploadAsset uploads an asset to a release on GitHub
func UploadAsset(id int, asset string) error {

	file, err := os.Open(asset)
	if err != nil {
		return fmt.Errorf("asset upload issue: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("asset upload issue: %w", err)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(asset))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	slog.Info("upload asset", "name", stat.Name(), "content-type", mimeType, "content-length", stat.Size())
	return nil
}
