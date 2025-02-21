# go-github-release
git tool to cut a releaser on github with assets

Tool create a new release at Github for a given owner and repos using a Github Personal Token and the upload assets from a target directory to this release. Finally we give the option to update the message body or publish the release.

## global install

```
go install github.com/dearing/go-github-release
```

## go tool install
```
go get -tool github.com/dearing/go-github-release@latest
```

# dev notes

We are going to avoid any SDKs and just called the native HTTP endpoint in order to keep the project within just the standard library.  The scope is limited to just uploading assets and creating releases so it shouldn't be rough.

## authentication

- `env:GITHUB_TOKEN` should be the only supported method

## create release

The idea should be to just create a OWNER/REPO release using our auth with a set of args.

- https://api.github.com/repos/OWNER/REPO/releases
- https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#create-a-release

### example from docs

```
curl -L \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <YOUR-TOKEN>" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/OWNER/REPO/releases \
  -d '{"tag_name":"v1.0.0","target_commitish":"master","name":"v1.0.0","body":"Description of the release","draft":false,"prerelease":false,"generate_release_notes":false}'

```

## upload assets

the idea should be that with an active release, we upload assets and name them.  Workflow could be a set of names passed in the cli AND/OR a directory to target and upload.


- support SNI for endpoint
- set `Release-Id` from CreateRelease
- set `Content-Length` to asset size
- set `Content-Type` to a MIME type (https://pkg.go.dev/mime)
- binaries => `application/octet-stream`
- archives => `application/zip`
- sumfiles => `text/plain`
- https://www.iana.org/assignments/media-types/media-types.xhtml
- https://docs.github.com/en/rest/releases/assets?apiVersion=2022-11-28#upload-a-release-asset


### example from docs
```
curl -L \
  -X POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <YOUR-TOKEN>" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  -H "Content-Type: application/octet-stream" \
  "https://uploads.github.com/repos/OWNER/REPO/releases/RELEASE_ID/assets?name=example.zip" \
  --data-binary "@example.zip"

```