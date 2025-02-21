package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	TagName              string `json:"tag_name,omitzero"`
	TargetCommitish      string `json:"target_commitish,omitzero"`
	Name                 string `json:"name,omitzero"`
	Body                 string `json:"body,omitzero"`
	Draft                bool   `json:"draft,omitzero"`
	Prerelease           bool   `json:"prerelease,omitzero"`
	GenerateReleaseNotes bool   `json:"generate_release_notes,omitzero"`
}

// CreateReleaseResponse is the response from creating a release on GitHub
type CreateReleaseResponse struct {
	URL             string    `json:"url,omitzero"`
	HTMLURL         string    `json:"html_url,omitzero"`
	AssetsURL       string    `json:"assets_url,omitzero"`
	UploadURL       string    `json:"upload_url,omitzero"`
	TarballURL      string    `json:"tarball_url,omitzero"`
	ZipballURL      string    `json:"zipball_url,omitzero"`
	DiscussionURL   string    `json:"discussion_url,omitzero"`
	ID              int       `json:"id,omitzero"`
	NodeID          string    `json:"node_id,omitzero"`
	TagName         string    `json:"tag_name,omitzero"`
	TargetCommitish string    `json:"target_commitish,omitzero"`
	Name            string    `json:"name,omitzero"`
	Body            string    `json:"body,omitzero"`
	Draft           bool      `json:"draft,omitzero"`
	Prerelease      bool      `json:"prerelease,omitzero"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Author          Author    `json:"author"`
	Assets          []Assets  `json:"assets,omitzero"`
}

// UploadAssetResponse is the response from uploading an asset to a release on GitHub
type UploadAssetResponse struct {
	URL                string    `json:"url,omitzero"`
	BrowserDownloadURL string    `json:"browser_download_url,omitzero"`
	ID                 int       `json:"id,omitzero"`
	NodeID             string    `json:"node_id,omitzero"`
	Name               string    `json:"name,omitzero"`
	Label              string    `json:"label,omitzero"`
	State              string    `json:"state,omitzero"`
	ContentType        string    `json:"content_type,omitzero"`
	Size               int       `json:"size,omitzero"`
	DownloadCount      int       `json:"download_count,omitzero"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Uploader           Uploader  `json:"uploader"`
}

// Author is the author of the release
type Author struct {
	Login             string `json:"login,omitzero"`
	ID                int    `json:"id,omitzero"`
	NodeID            string `json:"node_id,omitzero"`
	AvatarURL         string `json:"avatar_url,omitzero"`
	GravatarID        string `json:"gravatar_id,omitzero"`
	URL               string `json:"url,omitzero"`
	HTMLURL           string `json:"html_url,omitzero"`
	FollowersURL      string `json:"followers_url,omitzero"`
	FollowingURL      string `json:"following_url,omitzero"`
	GistsURL          string `json:"gists_url,omitzero"`
	StarredURL        string `json:"starred_url,omitzero"`
	SubscriptionsURL  string `json:"subscriptions_url,omitzero"`
	OrganizationsURL  string `json:"organizations_url,omitzero"`
	ReposURL          string `json:"repos_url,omitzero"`
	EventsURL         string `json:"events_url,omitzero"`
	ReceivedEventsURL string `json:"received_events_url,omitzero"`
	Type              string `json:"type,omitzero"`
	SiteAdmin         bool   `json:"site_admin,omitzero"`
}

// Uploader is the uploader of the asset
type Uploader struct {
	Login             string `json:"login,omitzero"`
	ID                int    `json:"id,omitzero"`
	NodeID            string `json:"node_id,omitzero"`
	AvatarURL         string `json:"avatar_url,omitzero"`
	GravatarID        string `json:"gravatar_id,omitzero"`
	URL               string `json:"url,omitzero"`
	HTMLURL           string `json:"html_url,omitzero"`
	FollowersURL      string `json:"followers_url,omitzero"`
	FollowingURL      string `json:"following_url,omitzero"`
	GistsURL          string `json:"gists_url,omitzero"`
	StarredURL        string `json:"starred_url,omitzero"`
	SubscriptionsURL  string `json:"subscriptions_url,omitzero"`
	OrganizationsURL  string `json:"organizations_url,omitzero"`
	ReposURL          string `json:"repos_url,omitzero"`
	EventsURL         string `json:"events_url,omitzero"`
	ReceivedEventsURL string `json:"received_events_url,omitzero"`
	Type              string `json:"type,omitzero"`
	SiteAdmin         bool   `json:"site_admin,omitzero"`
}

// Assets are the assets of the release
type Assets struct {
	URL                string    `json:"url,omitzero"`
	BrowserDownloadURL string    `json:"browser_download_url,omitzero"`
	ID                 int       `json:"id,omitzero"`
	NodeID             string    `json:"node_id,omitzero"`
	Name               string    `json:"name,omitzero"`
	Label              string    `json:"label,omitzero"`
	State              string    `json:"state,omitzero"`
	ContentType        string    `json:"content_type,omitzero"`
	Size               int       `json:"size,omitzero"`
	DownloadCount      int       `json:"download_count,omitzero"`
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

// uploadAsset uploads an asset to a release on GitHub.
// The endpoint you call to upload release assets is specific to your release.
// Use the upload_url returned in the response of the Create a release endpoint to upload a release asset.
//
//	ex: uploadAsset("github_pat_abc123", releaseResponse.UploadURL, "file.zip")
func uploadAsset(token, url, asset string) error {

	// https://uploads.github.com/repos/dearing/go-github-release/releases/201625270/assets{?name,label}

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

	slog.Info("upload asset stats", "name", stat.Name(), "content-type", mimeType, "content-length", stat.Size())

	// replace the query stub with the asset name key-value pair
	query := fmt.Sprintf("?name=%s", stat.Name()) // TODO: use label?
	url = strings.Replace(url, "{?name,label}", query, 1)

	slog.Info(">", "url", url)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", url, file)
	if err != nil {
		return fmt.Errorf("asset upload issue: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", mimeType)
	req.ContentLength = stat.Size()

	//req.Header.Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("asset upload issue: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != ReleaseCreated {
		data, _ := io.ReadAll(resp.Body)
		slog.Error("upload asset", "status", resp.Status, "data", string(data))

		return fmt.Errorf("asset upload issue status: %s", resp.Status)
	}

	slog.Info("upload asset", "status", resp.Status)

	info := UploadAssetResponse{}

	// decode the response into our struct
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return fmt.Errorf("asset upload decode issue: %w", err)
	}

	prettyPrint(info)

	return nil
}
